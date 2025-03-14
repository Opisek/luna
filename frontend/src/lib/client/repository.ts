import { browser } from "$app/environment";

import { writable } from "svelte/store";

import { AllChangesCalendar, AllChangesEvent, AllChangesSource, NoOp } from "./placeholders";
import { fetchJson, fetchResponse } from "./net";
import { getMetadata } from "./metadata";
import { queueNotification } from "./notifications";

import { atLeastOnePromise, deepCopy } from "$lib/common/misc";

class Repository {
  //
  // Constants
  //

  private readonly spoolerDelay = 50; // 50ms
  private readonly maxCacheAge = 1000 * 60 * 10; // 10 minutes

  //
  // Subscribeable Stores
  //

  readonly sources = writable([] as SourceModel[]);
  readonly calendars = writable([] as CalendarModel[]);
  readonly events = writable([] as EventModel[]);

  //
  // Constructor
  //
  constructor() {
    getMetadata().hiddenCalendars.subscribe(() => {
      this.compileEvents(this.eventsRangeStart, this.eventsRangeEnd);
    });

    if (browser) {
      window.addEventListener("storage", () => this.loadCache());
      this.loadCache();
    }
  }

  //
  // Misc
  //

  private async mapBaseEventToRecurrenceInstance(idAndDate: [string, number]): Promise<EventModel | null> {
    const baseEvent = this.eventsMap.get(idAndDate[0]);
    if (baseEvent == undefined) return null;

    const baseDuration = baseEvent.date.end.getTime() - baseEvent.date.start.getTime();

    const copiedEvent = await deepCopy(baseEvent);
    copiedEvent.date.start = new Date(idAndDate[1]);
    copiedEvent.date.end = new Date(idAndDate[1]);
    copiedEvent.date.end.setTime(copiedEvent.date.end.getTime() + baseDuration);

    return copiedEvent;
  }

  private async mapAllRecurrenceInstances(idAndDates: [string, number][]): Promise<EventModel[]> {
    return (await Promise.all(idAndDates.map(async (idAndDate) => this.mapBaseEventToRecurrenceInstance(idAndDate)))).filter(x => x != null);
  }

  //
  // Caching
  //

  private eventsRangeStart: Date = new Date();
  private eventsRangeEnd: Date = new Date();

  private readonly emptyCache = { date: 0, value: null };

  private lastCacheSave = Date.now();

  private sourcesCache: CacheEntry<SourceModel[]> = this.emptyCache; // sources
  private sourceDetailsCache: Map<string, CacheEntry<SourceModel>> = new Map(); // source -> details
  private calendarsCache: Map<string, CacheEntry<string[]>> = new Map(); // source -> calendars
  private eventsCache: Map<string, Map<number, CacheEntry<[string, number][]>>> = new Map(); // calendar -> month -> event id, start date (because of recurring events having the same id)
  private eventsMap: Map<string, EventModel> = new Map(); // event id -> event
  private calendarsMap: Map<string, CalendarModel> = new Map(); // calendar id -> calendar

  private cacheOk<T>(cache: CacheEntry<T> | undefined): (T | null) {
    return (cache && Date.now() - cache.date < this.maxCacheAge) ? cache.value : null;
  }

  invalidateCache() {
    this.sourcesCache.date = 0;
    this.sourceDetailsCache.forEach((cache) => cache.date = 0);
    this.calendarsCache.forEach((cache) => cache.date = 0);
    this.eventsCache.forEach((cache) => cache.forEach((entry) => entry.date = 0));
    this.saveCache();
  }

  //
  // Form Data
  //

  private getSourceFormData(source: SourceModel, changes: SourceModelChanges = AllChangesSource): FormData {
    const formData = new FormData();
    if (changes.name) formData.set("name", source.name);
    if (changes.type) formData.set("type", source.type);
    if (changes.type || changes.settings) {
      switch (source.type) {
        case "caldav":
          formData.set("url", source.settings.url);
          break;
        case "ical":
          switch (source.settings.location) {
            case "remote":
              formData.set("location", "remote");
              formData.set("url", source.settings.url);
              break
            case "database":
              formData.set("location", "database");
              formData.set("file", source.settings.file.item(0));
              break;
            case "local":
              formData.set("location", "local");
              formData.set("path", source.settings.path);
              break;
            default:
              throw new Error("Unsupported iCal file location");
          }
          break;
        default:
          throw new Error("Unsupported source type");
      }
    }
    if (changes.auth) {
      formData.set("auth_type", source.auth_type);
      switch (source.auth_type) {
        case "none":
          break;
        case "basic":
          formData.set("auth_username", source.auth.username);
          formData.set("auth_password", source.auth.password);
          break;
        case "bearer":
          formData.set("auth_token", source.auth.token);
          break;
        default:
          throw new Error("Unsupported auth type");
      }
    }
    return formData;
  }

  private getCalendarFormData(calendar: CalendarModel, changes: CalendarModelChanges = AllChangesCalendar): FormData {
    const formData = new FormData();
    if (changes.name) formData.set("name", calendar.name);
    if (changes.desc) formData.set("desc", calendar.desc);
    if (changes.color) formData.set("color", calendar.color);
    return formData;
  }

  private getEventFormData(event: EventModel, changes: EventModelChanges = AllChangesEvent): FormData {
    const formData = new FormData();
    if (changes.name) formData.set("name", event.name);
    if (changes.desc) formData.set("desc", event.desc);
    if (changes.date) {
      if (event.date.allDay) {
        const start = new Date(event.date.start.getTime() - (event.date.start.getTimezoneOffset() * 60000));
        const end = new Date(event.date.end.getTime() - (event.date.end.getTimezoneOffset() * 60000));
        formData.set("date_start", start.toISOString());
        formData.set("date_end", end.toISOString());
      } else {
        formData.set("date_start", event.date.start.toISOString());
        formData.set("date_end", event.date.end.toISOString());
      }
      formData.set("date_all_day", event.date.allDay ? "true" : "false");
    }
    if (changes.color) {
      if (event.color && event.color !== "") {
        formData.set("color", event.color);
      } else {
        formData.set("color", "null");
      }
    }
    return formData;
  }

  //
  // Visibility
  //

  private compileSources() {
    this.sources.set(this.sourcesCache.value || []);
  }

  private compileCalendarsTimeout: (ReturnType<typeof setTimeout> | undefined) = undefined;
  private compileCalendars() {
    clearTimeout(this.compileCalendarsTimeout);
    this.compileCalendarsTimeout = setTimeout(() => {
      const allCalendars = Array.from(this.calendarsCache.values().map(x => x.value).filter(x => x != null)).flat();
      this.calendars.set(allCalendars.map(x => this.calendarsMap.get(x)).filter(x => x != null));
    }, this.spoolerDelay)
  }

  private compileEventsTimeout: (ReturnType<typeof setTimeout> | undefined) = undefined;
  private compileEvents(start: Date, end: Date) {
    clearTimeout(this.compileEventsTimeout);

    this.compileEventsTimeout = setTimeout(async () => {
      const allEvents = 
        Array.from(this.eventsCache.entries())
        .filter(x => !getMetadata().hiddenCalendars.has(x[0]) && x[1] != null) // Event must be visible
        .map(x => Array.from(x[1].entries()))
        .flat()
        .filter(x => x[1] != null && x[0] >= start.getTime() && x[0] <= end.getTime()) // Event must be in the time frame
        .map(x => x[1].value)
        .filter(x => x != null) // Event must exist
        .flat();
      
      const uniqueEvents = [ ...new Map(allEvents.map(x => [`${x[0]}${x[1]}`, x])).values() ];

      const eventsWithData = await this.mapAllRecurrenceInstances(uniqueEvents);

      this.events.set(eventsWithData);
    }, this.spoolerDelay)
  }

  //
  // Local Storage
  //

  private loadCache() {
    const cacheTimestamp = localStorage.getItem("cache.timestamp");
    if (cacheTimestamp != null && Number.parseInt(cacheTimestamp) == this.lastCacheSave) return;

    const newSourcesCache = localStorage.getItem("cache.sources");
    if (newSourcesCache) this.sourcesCache = JSON.parse(newSourcesCache);

    const newSourceDetailsCache = localStorage.getItem("cache.sourceDetails");
    if (newSourceDetailsCache) this.sourceDetailsCache = new Map(JSON.parse(newSourceDetailsCache));

    const newCalendarsCache = localStorage.getItem("cache.calendars");
    if (newCalendarsCache) this.calendarsCache = new Map(JSON.parse(newCalendarsCache));

    const newEventsCache = localStorage.getItem("cache.events");
    if (newEventsCache) this.eventsCache = new Map(JSON.parse(newEventsCache).map((x: [string, [number, CacheEntry<string>][]]) => [x[0], new Map(x[1])]));

    const newCalendarsMap = localStorage.getItem("cache.calendarsMap");
    if (newCalendarsMap) this.calendarsMap = new Map(JSON.parse(newCalendarsMap));

    const newEventsMap = localStorage.getItem("cache.eventsMap");
    if (newEventsMap) {
      this.eventsMap = new Map(JSON.parse(newEventsMap));
      this.eventsMap.forEach((event) => {
        event.date.start = new Date(event.date.start);
        event.date.end = new Date(event.date.end);
        if (event.date.allDay) {
          event.date.start.setHours(0, 0, 0, 0);
          event.date.end.setHours(0, 0, 0, 0);
        }
      });
    }

    this.compileSources();
    this.compileCalendars();
    this.compileEvents(this.eventsRangeStart, this.eventsRangeEnd);
  }

  private saveCacheTimeout: (ReturnType<typeof setTimeout> | undefined) = undefined;
  private saveCache() {
    if (browser) {
      // @ts-ignore if only typescript were consistent...
      clearTimeout(this.saveCacheTimeout);
      setTimeout(() => {
        this.lastCacheSave = Date.now();
        localStorage.setItem("cache.timestamp", this.lastCacheSave.toString());
        localStorage.setItem("cache.sources", JSON.stringify(this.sourcesCache));
        localStorage.setItem("cache.sourceDetails", JSON.stringify(Array.from(this.sourceDetailsCache.entries())));
        localStorage.setItem("cache.calendars", JSON.stringify(Array.from(this.calendarsCache.entries())));
        localStorage.setItem("cache.events", JSON.stringify(Array.from(this.eventsCache.entries().map(x => [x[0], Array.from(x[1].entries())]))));
        localStorage.setItem("cache.calendarsMap", JSON.stringify(Array.from(this.calendarsMap.entries())));
        localStorage.setItem("cache.eventsMap", JSON.stringify(Array.from(this.eventsMap.entries())));
      }, this.spoolerDelay)
    }
  }

  // 
  // Sources
  // 

  async getSources(forceRefresh = false): Promise<SourceModel[]> {
    if (!browser) return [];

    if (!forceRefresh) {
      const cached = this.cacheOk(this.sourcesCache);
      if (cached) return Promise.resolve(cached);
    }

    const stopLoading = getMetadata().startLoading();

    const fetchedSources: SourceModel[] = await fetchJson("/api/sources").catch((err) => {
      throw err;
    }).finally(() => {
      stopLoading();
    });

    this.sourcesCache.date = Date.now(),
    this.sourcesCache.value = fetchedSources;
    this.compileSources();
    this.saveCache();
    return fetchedSources;
  }

  async getSourceDetails(id: string, forceRefresh = false): Promise<SourceModel> {
    if (!browser) return {} as SourceModel;

    if (!forceRefresh) {
      const cached = this.cacheOk(this.sourceDetailsCache.get(id));
      if (cached) return Promise.resolve(cached);
    }

    const fetched: SourceModel = await fetchJson(`/api/sources/${id}`).catch((err) => { throw err; });

    this.sourceDetailsCache.set(id, {
      date: Date.now(),
      value: fetched
    });
    this.saveCache();
    return fetched;
  }

  async createSource(newSource: SourceModel): Promise<void> {
    if (!browser) return;

    const formData = this.getSourceFormData(newSource);

    const json = await fetchJson(`/api/sources`, { method: "PUT", body: formData }).catch((err) => { throw err; });

    newSource.id = json.id;
    this.sourcesCache.value = this.sourcesCache.value?.concat(newSource) || [ newSource ];
    this.sources.update((sources) => sources.concat(newSource));

    this.getCalendars(newSource.id).then(async (cals) => {
      this.compileCalendars();
      const [_, errors] = await atLeastOnePromise(cals.map((cal) => this.getEventsFromCalendar(cal.id, this.eventsRangeStart, this.eventsRangeEnd))).catch(() => {
        throw new Error("Failed to fetch events");
      });
      errors.forEach((err) => {
        queueNotification(
          "failure",
          `Failed to fetch events from ${cals[err[0]].name}: ${err[1].message}`
        );
      });
      this.compileEvents(this.eventsRangeStart, this.eventsRangeEnd);
    }).catch((err) => {
      queueNotification(
        "failure",
        `Failed to fetch calendars from ${newSource.name}: ${err.message}`
      );
    });

    this.saveCache();
  }

  async editSource(modifiedSource: SourceModel, changes: SourceModelChanges): Promise<void> {
    if (!browser) return;

    let formData = this.getSourceFormData(modifiedSource, changes);

    await fetchResponse(`/api/sources/${modifiedSource.id}`, { method: "PATCH", body: formData }).catch((err) => { throw err; });
    
    this.sourcesCache.value = this.sourcesCache.value?.map((source => source.id === modifiedSource.id ? modifiedSource : source)) || [ modifiedSource ];
    this.compileSources();

    this.getCalendars(modifiedSource.id).then(async (cals) => {
      this.compileCalendars();
      const [_, errors] = await atLeastOnePromise(cals.map((cal) => this.getEventsFromCalendar(cal.id, this.eventsRangeStart, this.eventsRangeEnd))).catch(() => {
        throw new Error("Failed to fetch events");
      });
      errors.forEach((err) => {
        queueNotification(
          "failure",
          `Failed to fetch events from ${cals[err[0]].name}: ${err[1].message}`
        );
      });
      this.compileEvents(this.eventsRangeStart, this.eventsRangeEnd);
    }).catch((err) => {
      queueNotification(
        "failure",
        `Failed to fetch calendars from ${modifiedSource.name}: ${err.message}`
      );
    });

    this.saveCache();
  }

  async deleteSource(id: string): Promise<void> {
    if (!browser) return;

    await fetchResponse(`/api/sources/${id}`, { method: "DELETE" }).catch((err) => { throw err; });

    this.sourcesCache.value = this.sourcesCache.value?.filter((source) => source.id !== id) || [];
    this.compileSources();

    for (const calendar of this.calendarsCache.get(id)?.value || []) {
      this.eventsCache.delete(calendar);
      this.calendarsMap.delete(calendar);
    }

    this.compileCalendars();
    this.compileEvents(this.eventsRangeStart, this.eventsRangeEnd);

    this.saveCache();
  }

  //
  // Calendars
  //

  async getAllCalendars(forceRefresh = false): Promise<CalendarModel[]> {
    if (!browser) return [];

    const allSources = await this.getSources(forceRefresh);

    const [calendars, errors] = await atLeastOnePromise(allSources.map((source) => this.getCalendars(source.id, forceRefresh))).catch(() => {
      throw new Error("Failed to fetch calendars");
    });

    errors.forEach((err) => {
      queueNotification(
        "failure",
        `Failed to fetch calendars from ${allSources[err[0]].name}: ${err[1].message}`
      );
    });

    return calendars.flat();
  }

  private async getCalendars(id: string, forceRefresh = false): Promise<CalendarModel[]> {
    if (!browser) return [];

    if (!forceRefresh) {
      const cached = this.cacheOk(this.calendarsCache.get(id));
      if (cached) return Promise.resolve(cached.map(x => this.calendarsMap.get(x)).filter(x => x != null));
    }

    const stopLoading = getMetadata().startLoadingSource(id);

    const response = await fetchJson(`/api/sources/${id}/calendars`).catch((err) => {
      getMetadata().addFaultySource(id, err.message);
      throw err;
    }).finally(() => {
      stopLoading();
    });

    getMetadata().removeFaultySource(id);

    const fetched: CalendarModel[] = response.calendars;

    // Delete orphans
    let anyOrphaned = false;
    const fetchedSet = new Set(fetched.map(x => x.id));
    for (const calendar of this.calendarsCache.get(id)?.value || []) {
      if (!fetchedSet.has(calendar)) {
        this.eventsCache.delete(calendar);
        anyOrphaned = true;
      }
    }

    for (const calendar of fetched) {
      this.calendarsMap.set(calendar.id, calendar)
    }
    
    this.calendarsCache.set(id, {
      date: Date.now(),
      value: fetched.map(x => x.id)
    });

    this.compileCalendars();
    if (anyOrphaned) this.compileEvents(this.eventsRangeStart, this.eventsRangeEnd);
    this.saveCache();
    return fetched;
  }

  async getCalendar(id: string, forceRefresh = false): Promise<CalendarModel | null> {
    if (!browser) return {} as CalendarModel;

    return this.calendarsMap.get(id) || null; // TODO: needs a bit of refactoring plus actual fetch of the relevant endpoint depending on cache age
  }

  async createCalendar(newCalendar: CalendarModel): Promise<void> {
    if (!browser) return;

    // add to database
    const formData = this.getCalendarFormData(newCalendar);

    const json = await fetchJson(`/api/sources/${newCalendar.source}/calendars`, { method: "PUT", body: formData }).catch((err) => { throw err; });

    newCalendar.id = json.id;

    // add to cache
    this.calendarsMap.set(newCalendar.id, newCalendar);
    this.calendarsCache.set(newCalendar.source, {
      date: this.calendarsCache.get(newCalendar.source)?.date || Date.now(),
      value: this.calendarsCache.get(newCalendar.source)?.value?.concat(newCalendar.id) || [ newCalendar.id ]
    });

    // add to display
    this.calendars.update((calendars) => calendars.concat(newCalendar));

    this.saveCache();
  };

  async editCalendar(modifiedCalendar: CalendarModel, changes: CalendarModelChanges): Promise<void> {
    if (!browser) return;

    // update in database
    const formData = this.getCalendarFormData(modifiedCalendar, changes);

    await fetchResponse(`/api/calendars/${modifiedCalendar.id}`, { method: "PATCH", body: formData }).catch((err) => { throw err; });

    // update in cache
    this.calendarsMap.set(modifiedCalendar.id, modifiedCalendar);

    // update on display
    this.calendars.update((calendars) => calendars.map((cal) => cal.id === modifiedCalendar.id ? modifiedCalendar : cal));

    this.saveCache();
  }

  async deleteCalendar(id: string): Promise<void> {
    if (!browser) return;

    // remove from database
    await fetchResponse(`/api/calendars/${id}`, { method: "DELETE" }).catch((err) => { throw err; });

    const calendar = this.calendarsMap.get(id);
    if (!calendar) return;
    this.calendarsMap.delete(id);

    // remove from cache
    this.eventsCache.delete(id);
    this.calendarsCache.set(calendar.source, {
      date: this.calendarsCache.get(calendar.source)?.date || Date.now(),
      value: this.calendarsCache.get(calendar.source)?.value?.filter((cal) => cal !== id) || []
    });

    // remove from display
    this.calendars.update((calendars) => calendars.filter((cal) => cal.id !== id));
    this.compileEvents(this.eventsRangeStart, this.eventsRangeEnd);

    this.saveCache();
  }

  async moveCalendar(calendar: CalendarModel): Promise<void> {
    throw new Error("Not implemented");

    if (!browser) return;

    const oldId = calendar.id;

    // add to the new calendar
    await this.createCalendar(calendar).catch((err) => { throw err; });

    // TODO: MOVE ALL EVENTS!!!

    // remove from the old calendar
    await this.deleteCalendar(calendar.id).catch((err) => {
      // undo changes
      this.deleteCalendar(calendar.id).catch(NoOp);
      calendar.id = oldId;
      throw err;
    });
  }

  //
  // Events
  //

  private determineEventMonths(event: EventModel): Date[] {
    const start = new Date(event.date.start);
    start.setUTCDate(1);
    start.setUTCHours(0, 0, 0, 0);
    const end = new Date(event.date.end);

    const months = [];
    while (start < end) {
      start.setUTCMonth(start.getUTCMonth() + 1);
      months.push(new Date(start));
    }

    return months;
  }

  private addEventToCache(event: EventModel, date: Date) {
    let calendarEventsCache = this.eventsCache.get(event.calendar);
    if (!calendarEventsCache) {
      calendarEventsCache = new Map();
      this.eventsCache.set(event.calendar, calendarEventsCache);
    }

    let cacheEntry = calendarEventsCache.get(date.getTime());
    if (!cacheEntry) {
      cacheEntry = { date: Date.now(), value: [] };
      calendarEventsCache.set(date.getTime(), cacheEntry);
    }

    calendarEventsCache.set(date.getTime(), {
      date: cacheEntry.date,
      value: [ ...cacheEntry.value || [], [event.id, event.date.start.getTime()] ]
    });
  }

  private removeEventFromCache(event: EventModel, date: Date) {
    if (!this.eventsCache.has(event.calendar)) return;
    const cacheEntry = this.eventsCache.get(event.calendar)?.get(date.getTime());
    if (!cacheEntry || !cacheEntry.value) return;
    cacheEntry.value = cacheEntry.value.filter((idAndDate) => idAndDate[0] !== event.id);
  }

  async getAllEvents(start: Date, end: Date, forceRefresh = false): Promise<EventModel[]> {
    if (!browser) return [];
    this.eventsRangeStart = start;
    this.eventsRangeEnd = end;

    start.setUTCHours(0, 0, 0, 0);
    end.setUTCHours(23, 59, 59, 999);

    // Set start and end to the start and end of each month
    start.setDate(1);
    end.setMonth(end.getMonth() + 1);
    end.setDate(0);

    // Add one month of padding in both directions
    start.setMonth(start.getMonth() - 1);
    end.setMonth(end.getMonth() + 1);

    this.compileEvents(start, end);
    const allSources = await this.getSources(forceRefresh).catch((err) => { throw err; });
    const [events, errors] = await atLeastOnePromise(allSources.map((source) => this.getEventsFromSource(source.id, start, end, forceRefresh))).catch(() => {
      throw new Error("Failed to fetch events");
    });
    errors.forEach((err) => {
      queueNotification(
        "failure",
        `Failed to fetch events from ${allSources[err[0]].name}: ${err[1].message}`
      );
    });
    return events.flat();
  }

  private async getEventsFromSource(source: string, start: Date, end: Date, forceRefresh = false): Promise<EventModel[]> {
    let cals = await this.getCalendars(source, forceRefresh).catch((err) => { throw err; });
    cals = cals.filter(x => !getMetadata().hiddenCalendars.has(x.id)); // only fetch events from visible calendars
    const [events, errors] = await atLeastOnePromise(cals.map((calendar) => this.getEventsFromCalendar(calendar.id, start, end, forceRefresh))).catch(() => {
      throw new Error("Failed to fetch events");
    })
    errors.forEach((err) => {
      queueNotification(
        "failure",
        `Failed to fetch events from ${cals[err[0]].name}: ${err[1].message}`
      );
    });
    return events.flat();
  }

  private async getEventsFromCalendar(calendar: string, start: Date, end: Date, forceRefresh = false): Promise<EventModel[]> {
    let result: EventModel[] = [];

    const cache = (forceRefresh ? null : this.eventsCache.get(calendar)) || new Map<number, CacheEntry<[string, number][]>>();

    // Limit the range we need to ask for but only use one request per calendar
    const fetchStart = new Date(start);
    const fetchEnd = new Date(end);

    while (fetchStart.getTime() <= fetchEnd.getTime()) {
      const cached = this.cacheOk(cache.get(fetchStart.getTime()));
      if (!cached) break;
      result = result.concat(await this.mapAllRecurrenceInstances(cached));
      fetchStart.setMonth(fetchStart.getMonth() + 1);
    }

    while (fetchEnd.getTime() >= fetchStart.getTime()) {
      const cached = this.cacheOk(cache.get(fetchEnd.getTime()));
      if (!cached) break;
      result = result.concat(await this.mapAllRecurrenceInstances(cached));
      fetchEnd.setMonth(fetchEnd.getMonth() - 1);
    }

    if (fetchStart.getTime() > fetchEnd.getTime()) {
      return result;
    }

    const stopLoading = getMetadata().startLoadingCalendar(calendar);

    const fetchedEvents = await this.fetchEvents(calendar, start, end).catch((err) => {
      getMetadata().addFaultyCalendar(calendar, err.message);
      throw err;
    }).finally(() => {
      stopLoading();
    });

    getMetadata().removeFaultyCalendar(calendar);

    this.compileEvents(start, end);
    this.saveCache();
    return result.concat(fetchedEvents);
  }

  async fetchEvents(calendar: string, start: Date, end: Date): Promise<EventModel[]> {
    const localStart = new Date(start);
    localStart.setHours(0, 0, 0, 0);
    const localEnd = new Date(end);
    localEnd.setHours(23, 59, 59, 999);

    const fetched: EventModel[] = (await fetchJson(`/api/calendars/${calendar}/events?start=${encodeURIComponent(localStart.toISOString())}&end=${encodeURIComponent(localEnd.toISOString())}`)).events;

    let calendarEventsCache = this.eventsCache.get(calendar);
    if (!calendarEventsCache) {
      calendarEventsCache = new Map();
      this.eventsCache.set(calendar, calendarEventsCache);
    }
    for (let i = new Date(start); i.getTime() < end.getTime(); i.setMonth(i.getMonth() + 1)) {
      calendarEventsCache.set(i.getTime(), {
        date: Date.now(),
        value: []
      });
    }

    for (const event of fetched) {
      event.date.start = new Date(event.date.start);
      event.date.end = new Date(event.date.end);

      if (event.date.allDay) {
        event.date.start.setHours(0, 0, 0, 0);
        event.date.end.setHours(0, 0, 0, 0);
      }

      this.eventsMap.set(event.id, event);
      for (const month of this.determineEventMonths(event)) {
        this.addEventToCache(event, month);
      }
    }

    return fetched;
  }

  async getEventsFromPreviouslyHiddenCalendar(calendar: string) {
    if (!browser) return;

    this.getEventsFromCalendar(calendar, this.eventsRangeStart, this.eventsRangeEnd).catch((err) => {
      const calendarName = this.calendarsMap.get(this.calendarsCache.get(calendar)?.value?.find((cal) => cal === calendar) || "")?.name;
      queueNotification(
        "failure",
        `Failed to fetch events from calendar${calendarName ? " " + calendarName : ""}: ${err.message}`
      );
    });
  }

  async createEvent(newEvent: EventModel): Promise<void> {
    if (!browser) return;

    // add to database
    if (newEvent.date.allDay) {
      newEvent.date.start.setHours(0, 0, 0, 0);
      newEvent.date.end.setHours(0, 0, 0, 0);
    }

    const formData = this.getEventFormData(newEvent);

    const json = await fetchJson(`/api/calendars/${newEvent.calendar}/events`, { method: "PUT", body: formData }).catch((err) => { throw err; });

    newEvent.id = json.id;

    // add to cache
    this.eventsMap.set(newEvent.id, newEvent);
    for (const month of this.determineEventMonths(newEvent)) this.addEventToCache(newEvent, month);

    // add to display
    const isHidden = getMetadata().hiddenCalendars.has(newEvent.calendar)
    if (!isHidden && newEvent.date.start <= this.eventsRangeEnd && newEvent.date.end >= this.eventsRangeStart) this.events.update((events) => events.concat(newEvent));

    // info if the event is hidden
    if (isHidden) {
      queueNotification(
        "info",
        `Event ${newEvent.name} was added to a hidden calendar.`
      );
    }

    this.saveCache();
  };

  async editEvent(modifiedEvent: EventModel, changes: EventModelChanges): Promise<void> {
    if (!browser) return;

    // update in database
    if (modifiedEvent.date.allDay) {
      modifiedEvent.date.start.setHours(0, 0, 0, 0);
      modifiedEvent.date.end.setHours(0, 0, 0, 0);
    }

    const formData = this.getEventFormData(modifiedEvent, changes);

    await fetchResponse(`/api/events/${modifiedEvent.id}`, { method: "PATCH", body: formData }).catch((err) => { throw err; });

    // update in cache
    const previousMonths = this.determineEventMonths(this.eventsMap.get(modifiedEvent.id)!);
    const currentMonths = this.determineEventMonths(modifiedEvent);
    this.eventsMap.set(modifiedEvent.id, modifiedEvent);

    for (const month of previousMonths) {
      if (!currentMonths.includes(month)) {
        this.removeEventFromCache(modifiedEvent, month);
      }
    }

    for (const month of currentMonths) {
      if (!previousMonths.includes(month)) {
        this.addEventToCache(modifiedEvent, month);
      }
    }

    // update on display
    if (modifiedEvent.date.start <= this.eventsRangeEnd && modifiedEvent.date.end >= this.eventsRangeStart) {
      this.events.update((events) => events.map((event) => event.id === modifiedEvent.id ? modifiedEvent : event));
    } else {
      this.events.update((events) => events.filter((event) => event.id !== modifiedEvent.id));
    }

    this.saveCache();
  }

  async deleteEvent(id: string): Promise<void> {
    if (!browser) return;

    // remove from database
    await fetchResponse(`/api/events/${id}`, { method: "DELETE" }).catch((err) => { throw err; });

    const event = this.eventsMap.get(id);
    if (!event) return;
    this.eventsMap.delete(id);

    // remove from cache
    const months = this.determineEventMonths(event);
    for (const month of months) this.removeEventFromCache(event, month);

    // remove from display
    this.events.update((events) => events.filter((event) => event.id !== id));

    this.saveCache();
  }

  async moveEvent(event: EventModel): Promise<void> {
    if (!browser) return;

    const oldId = event.id;

    // add to the new calendar
    await this.createEvent(event).catch((err) => { throw err; });

    // remove from the old calendar
    await this.deleteEvent(oldId).catch((err) => {
      // undo changes
      this.deleteEvent(event.id).catch(NoOp);
      event.id = oldId;
      throw err;
    });
  }
}

let repository: Repository | null = null;
export function getRepository() {
  if (!repository) repository = new Repository();
  return repository;
}

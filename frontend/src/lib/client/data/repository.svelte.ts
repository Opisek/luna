import { browser } from "$app/environment";
import { page } from "$app/state";

import { AllChangesCalendar, AllChangesEvent, AllChangesSource, NoChangesCalendar, NoChangesEvent, NoOp } from "../placeholders";
import { ColorKeys } from "../../../types/colors";
import { fetchJson, fetchResponse } from "../net";
import { queueNotification } from "../notifications";

import { parallel } from "$lib/common/misc";
import type { Metadata } from "./metadata.svelte";


export class Repository {
  private metadata: Metadata;

  //
  // Constants
  //

  private readonly spoolerDelay = 50; // 50ms
  private readonly maxCacheAge = 1000 * 60 * 10; // 10 minutes

  //
  // Subscribeable Stores
  //

  sources = $state<SourceModel[]>([]);
  calendars = $state<CalendarModel[]>([]);
  events = $state<EventModel[]>([]);

  //
  // Constructor
  //
  constructor(metadata: Metadata) {
    this.metadata = metadata;

    this.sources = [];
    this.calendars = [];
    this.events = [];

    if (browser) {
      window.addEventListener("storage", () => this.loadCache());
      this.loadCache();
    }
  }

  //
  // Caching
  //

  private eventsRangeStart: number = this.getMonthFromDate(new Date());
  private eventsRangeEnd: number = this.eventsRangeStart;

  private readonly emptyCache = { date: 0, value: null };

  private lastCacheSave = Date.now();

  private sourcesCache: CacheEntry<SourceModel[]> = this.emptyCache; // sources
  private sourceDetailsCache: Map<string, CacheEntry<SourceModel>> = new Map(); // source -> details
  private calendarsCache: Map<string, CacheEntry<string[]>> = new Map(); // source -> calendars
  private eventsCache: Map<string, Map<number, CacheEntry<string[]>>> = new Map(); // calendar -> month -> event ids
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

  private pendingMonths: Map<string, Set<number>> = new Map();

  private getMonthFromDate(time: Date): number {
    return time.getFullYear() * 100 + time.getMonth() + 1;
  }
  private getDateFromMonth(month: number): Date {
    const onlyMonth = month % 100;
    const year = (month - onlyMonth) / 100;
    return new Date(year, onlyMonth - 1, 1);
  }
  private nextMonth(month: number): number {
    return month + ((month % 100 == 12) ? 89 : 1);
  }
  private previousMonth(month: number): number {
    return month - ((month % 100 == 1) ? 89 : 1);
  }
  private isMonthPending(calendar: string, month: number): boolean {
    let monthsPendingForCalendar = this.pendingMonths.get(calendar);
    if (!monthsPendingForCalendar) return false;
    return monthsPendingForCalendar.has(month);
  }
  private setMonthPending(calendar: string, month: number) {
    let monthsPendingForCalendar = this.pendingMonths.get(calendar);
    if (!monthsPendingForCalendar) monthsPendingForCalendar = new Set();
    monthsPendingForCalendar.add(month);
    this.pendingMonths.set(calendar, monthsPendingForCalendar);
  }
  private setMonthNotPending(calendar: string, month: number) {
    let monthsPendingForCalendar = this.pendingMonths.get(calendar);
    if (!monthsPendingForCalendar) return;
    monthsPendingForCalendar.delete(month);
    this.pendingMonths.set(calendar, monthsPendingForCalendar);
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
        case "google":
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
        case "oauth":
          formData.set("auth_client", source.auth.client_id);
          formData.set("auth_tokens", source.auth.tokens_id);
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
    this.sources = this.sourcesCache.value || [];
  }

  private compileCalendarsTimeout: (ReturnType<typeof setTimeout> | undefined) = undefined;
  private compileCalendars() {
    clearTimeout(this.compileCalendarsTimeout);
    this.compileCalendarsTimeout = setTimeout(() => {
      const allCalendars = Array.from(this.calendarsCache.values().map(x => x.value).filter(x => x != null)).flat();
      this.calendars = allCalendars.map(x => this.calendarsMap.get(x)).filter(x => x != null);
    }, this.spoolerDelay)
  }

  private compileEventsTimeout: (ReturnType<typeof setTimeout> | undefined) = undefined;
  private compileEvents(startMonth: number, endMonth: number) {
    clearTimeout(this.compileEventsTimeout);
    this.compileEventsTimeout = setTimeout(async () => {
      this.events = [ ...new Map(
        Array.from(this.eventsCache.entries())
        .filter(x => !this.metadata.hiddenCalendars.has(x[0]) && x[1] != null) // Event must be visible
        .map(x => Array.from(x[1].entries()))
        .flat()
        .filter(x => x[1] != null && x[0] >= startMonth && x[0] <= endMonth) // Event must be in the time frame
        .map(x => x[1].value)
        .filter(x => x != null) // Event must exist
        .flat()
        .map(x => this.eventsMap.get(x))
        .filter(x => x != null) // Event must exist
        .map(x => [x.id, x])
      ).values()];
    }, this.spoolerDelay)
  }
  public async recalculateEvents(calendarThatBecameVisible: (string | null) = null) {
    if (calendarThatBecameVisible != null) {
      await this.getEventsFromPreviouslyHiddenCalendar(calendarThatBecameVisible);
    }
    this.compileEvents(this.eventsRangeStart, this.eventsRangeEnd);
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

    const stopLoading = this.metadata.startLoading();

    const fetchedSources: SourceModel[] = (await fetchJson("/api/sources").catch((err) => {
      throw err;
    }).finally(() => {
      stopLoading();
    })).sources;

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

    const fetched: SourceModel = (await fetchJson(`/api/sources/${id}`).catch((err) => { throw err; })).source;

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
    this.sources.push(newSource);

    this.getCalendars(newSource.id).then(async (cals) => {
      this.compileCalendars();
      const [_, errors] = await parallel(cals.map((cal) => this.getEventsFromCalendar(cal.id, this.eventsRangeStart, this.eventsRangeEnd))).catch((err) => {
        throw new Error(`Failed to fetch events from ${newSource.name}: ${(err.cause || err).message}`, { cause: err.cause || err });
      });
      errors.forEach((err) => {
        queueNotification(
          ColorKeys.Danger,
          `Failed to fetch events from ${cals[err[0]].name}: ${err[1].message}`
        );
      });
      this.compileEvents(this.eventsRangeStart, this.eventsRangeEnd);
    }).catch((err) => {
      queueNotification(
        ColorKeys.Danger,
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
    const detailsCacheEntry = this.sourceDetailsCache.get(modifiedSource.id);
    if (detailsCacheEntry) detailsCacheEntry.value = modifiedSource;
    this.compileSources();

    if (changes.settings) {
      for (const calendar of this.calendarsCache.get(modifiedSource.id)?.value || []) {
        this.eventsCache.delete(calendar);
        this.calendarsMap.delete(calendar);
      }

      this.getCalendars(modifiedSource.id, true).then(async (cals) => {
        this.compileCalendars();
        const [_, errors] = await parallel(cals.map((cal) => this.getEventsFromCalendar(cal.id, this.eventsRangeStart, this.eventsRangeEnd))).catch((err) => {
          throw new Error(`Failed to fetch events from ${modifiedSource.name}: ${(err.cause || err).message}`, { cause: err.cause || err });
        });
        errors.forEach((err) => {
          queueNotification(
            ColorKeys.Danger,
            `Failed to fetch events from ${cals[err[0]].name}: ${err[1].message}`
          );
        });
        this.compileEvents(this.eventsRangeStart, this.eventsRangeEnd);
      }).catch((err) => {
        queueNotification(
          ColorKeys.Danger,
          `Failed to fetch calendars from ${modifiedSource.name}: ${err.message}`
        );
      });
    }

    this.saveCache();
  }

  async changeSourceDisplayOrder(movedSource: SourceModel, newIndex: number): Promise<void> {
    if (!browser) return;

    const sources = this.sourcesCache.value || [];
    const previousIndex = sources.findIndex(x => x.id == movedSource.id);
    if (previousIndex == -1) throw new Error("Could not move cached source");
    if (previousIndex == newIndex) return;

    let formData = new FormData();
    formData.append("index", newIndex.toString())

    await fetchResponse(`/api/sources/${movedSource.id}/order`, { method: "POST", body: formData }).catch((err) => { throw err; });

    const original = sources[previousIndex];

    const direction = previousIndex < newIndex ? 1 : -1;
    for (let i = previousIndex; i != newIndex; i += direction) {
      sources[i] = sources[i + direction];
    }
    sources[newIndex] = original;

    this.sourcesCache.value = sources;
    this.compileSources();
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

    const [calendars, errors] = await parallel(allSources.map((source) => this.getCalendars(source.id, forceRefresh))).catch((err) => {
      throw new Error(`Failed to fetch calendars: ${(err.cause || err).message}`, { cause: err.cause || err });
    });

    errors.forEach((err) => {
      queueNotification(
        ColorKeys.Danger,
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

    const stopLoading = this.metadata.startLoadingSource(id);

    const response = await fetchJson(`/api/sources/${id}/calendars`).catch((err) => {
      this.metadata.addFaultySource(id, err.message);
      throw err;
    }).finally(() => {
      stopLoading();
    });

    this.metadata.removeFaultySource(id);

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

  async getCalendar(id: string, forceRefresh = false): Promise<CalendarModel> {
    // TODO: needs refactoring like integrating cache age check
    const fetched: CalendarModel = (await fetchJson(`/api/calendars/${id}`).catch((err) => { throw err; })).calendar;
    this.calendarsMap.set(fetched.id, fetched);
    this.compileCalendars();
    return fetched;
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
    this.calendars.push(newCalendar);

    this.saveCache();
  };

  async editCalendar(modifiedCalendar: CalendarModel, changes: CalendarModelChanges, override: boolean): Promise<void> {
    if (!browser) return;

    // update in database
    if (override && changes !== NoChangesCalendar) changes.color = true; // we have no way to destinguish between "don't change color" and "default color"
    const formData = this.getCalendarFormData(modifiedCalendar, changes);
    formData.set("overridden", override ? "true" : "false");

    await fetchResponse(`/api/calendars/${modifiedCalendar.id}`, { method: "PATCH", body: formData }).catch((err) => { throw err; });

    // update in cache
    modifiedCalendar.overridden = override;
    this.calendarsMap.set(modifiedCalendar.id, modifiedCalendar);

    // update on display
    //this.calendars.update((calendars) => calendars.map((cal) => cal.id === modifiedCalendar.id ? modifiedCalendar : cal));
    if (changes.color) this.compileEvents(this.eventsRangeStart, this.eventsRangeEnd);

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
    this.calendars.splice(this.calendars.findIndex((cal) => cal.id === id), 1);
    this.compileEvents(this.eventsRangeStart, this.eventsRangeEnd);

    this.saveCache();
  }

  async moveCalendar(calendar: CalendarModel): Promise<void> {
    throw new Error("Not implemented");

    //if (!browser) return;

    //const oldId = calendar.id;

    //// add to the new calendar
    //await this.createCalendar(calendar).catch((err) => { throw err; });

    //// TODO: MOVE ALL EVENTS!!!

    //// remove from the old calendar
    //await this.deleteCalendar(calendar.id).catch((err) => {
    //  // undo changes
    //  this.deleteCalendar(calendar.id).catch(NoOp);
    //  calendar.id = oldId;
    //  throw err;
    //});
  }

  //
  // Events
  //

  private determineEventMonths(event: EventModel): number[] {
    let start = this.getMonthFromDate(event.date.start);
    let end = this.getMonthFromDate(event.date.end);

    const months = [];
    while (start <= end) {
      months.push(start);
      start = this.nextMonth(start);
    }

    return months;
  }

  private addEventToCache(event: EventModel, month: number) {
    let calendarEventsCache = this.eventsCache.get(event.calendar);
    if (!calendarEventsCache) {
      calendarEventsCache = new Map();
      this.eventsCache.set(event.calendar, calendarEventsCache);
    }

    let cacheEntry = calendarEventsCache.get(month);
    if (!cacheEntry) {
      cacheEntry = { date: Date.now(), value: [] };
      calendarEventsCache.set(month, cacheEntry);
    }

    calendarEventsCache.set(month, {
      date: cacheEntry.date,
      value: [ ...cacheEntry.value || [], event.id ]
    });
  }

  private removeEventFromCache(event: EventModel, month: number) {
    if (!this.eventsCache.has(event.calendar)) return;
    const cacheEntry = this.eventsCache.get(event.calendar)?.get(month);
    if (!cacheEntry || !cacheEntry.value) return;
    cacheEntry.value = cacheEntry.value.filter((idAndDate) => idAndDate[0] !== event.id);
  }

  async getAllEvents(start: Date, end: Date, forceRefresh = false): Promise<EventModel[]> {
    if (!browser) return [];

    this.eventsRangeStart = this.getMonthFromDate(start);
    this.eventsRangeEnd = this.getMonthFromDate(end);

    this.compileEvents(this.eventsRangeStart, this.eventsRangeEnd);
    const allSources = await this.getSources(forceRefresh).catch((err) => { throw err; });
    const [events, errors] = await parallel(allSources.map((source) => this.getEventsFromSource(source.id, this.previousMonth(this.eventsRangeStart), this.nextMonth(this.eventsRangeEnd), forceRefresh))).catch((err) => {
      throw new Error(`Failed to fetch events: ${(err.cause || err).message}`, { cause: err.cause || err });
    });
    errors.forEach((err) => {
      queueNotification(
        ColorKeys.Danger,
        `Failed to fetch events from ${allSources[err[0]].name}: ${err[1].message}`
      );
    });
    return events.flat();
  }

  async getEvent(id: string, forceRefresh = false): Promise<EventModel> {
    // TODO: needs refactoring like integrating cache age check
    const fetched: EventModel = (await fetchJson(`/api/events/${id}`).catch((err) => { throw err; })).event;
    fetched.date.start = new Date(fetched.date.start);
    fetched.date.end = new Date(fetched.date.end);
    if (fetched.date.allDay) {
      fetched.date.start.setHours(0, 0, 0, 0);
      fetched.date.end.setHours(0, 0, 0, 0);
    }
    this.eventsMap.set(fetched.id, fetched);
    this.compileEvents(this.eventsRangeStart, this.eventsRangeEnd);
    return fetched;
  }

  private async getEventsFromSource(source: string, startMonth: number, endMonth: number, forceRefresh = false): Promise<EventModel[]> {
    let cals = await this.getCalendars(source, forceRefresh).catch((err) => { throw err; });
    cals = cals.filter(x => !this.metadata.hiddenCalendars.has(x.id)); // only fetch events from visible calendars
    const [events, errors] = await parallel(cals.map((calendar) => this.getEventsFromCalendar(calendar.id, startMonth, endMonth, forceRefresh))).catch((err) => {
      throw new Error(`Failed to fetch events: ${(err.cause || err).message}`, { cause: err.cause || err });
    })
    errors.forEach((err) => {
      queueNotification(
        ColorKeys.Danger,
        `Failed to fetch events from ${cals[err[0]].name}: ${err[1].message}`
      );
    });
    return events.flat();
  }

  private async getEventsFromCalendar(calendar: string, startMonth: number, endMonth: number, forceRefresh = false): Promise<EventModel[]> {
    let result: EventModel[] = [];

    const cache = (forceRefresh ? null : this.eventsCache.get(calendar)) || new Map<number, CacheEntry<string[]>>();

    // Determine which months must be fetched and which can be taken from cache
    while (startMonth <= endMonth) {
      if (!this.isMonthPending(calendar, startMonth)) {
        const cached = this.cacheOk(cache.get(startMonth));
        if (!cached) break;
        result = result.concat(cached.map(x => this.eventsMap.get(x)).filter(x => x != null));
      } else {
      }
      startMonth = this.nextMonth(startMonth);
    }
    while (endMonth >= startMonth) {
      if (!this.isMonthPending(calendar, endMonth)) {
        const cached = this.cacheOk(cache.get(endMonth));
        if (!cached) break;
        result = result.concat(cached.map(x => this.eventsMap.get(x)).filter(x => x != null));
      } else {
      }
      endMonth = this.previousMonth(endMonth);
    }

    if (startMonth > endMonth) {
      return result;
    }

    const stopLoading = this.metadata.startLoadingCalendar(calendar);

    // Fetch events
    for (let month = startMonth; month <= endMonth; month = this.nextMonth(month)) this.setMonthPending(calendar, month);
    const fetchedEvents = await this.fetchEvents(calendar, startMonth, endMonth).catch((err) => {
      this.metadata.addFaultyCalendar(calendar, err.message);
      for (let month = startMonth; month <= endMonth; month = this.nextMonth(month)) this.setMonthNotPending(calendar, month);
      throw err;
    }).finally(() => {
      stopLoading();
    });

    // Clear event cache for the requested months
    let calendarEventsCache = this.eventsCache.get(calendar);
    if (!calendarEventsCache) {
      calendarEventsCache = new Map();
      this.eventsCache.set(calendar, calendarEventsCache);
    }
    for (let month = startMonth; month <= endMonth; month = this.nextMonth(month)) {
      calendarEventsCache.set(month, {
        date: Date.now(),
        value: []
      });
    }
    this.eventsCache.set(calendar, calendarEventsCache);

    // Fill the cache with the fetched events
    for (const event of fetchedEvents) {
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

    // Done
    for (let month = startMonth; month <= endMonth; month = this.nextMonth(month)) this.setMonthNotPending(calendar, month);
    this.metadata.removeFaultyCalendar(calendar);
    this.compileEvents(this.eventsRangeStart, this.eventsRangeEnd);
    this.saveCache();
    return result.concat(fetchedEvents);
  }

  async fetchEvents(calendar: string, startMonth: number, endMonth: number): Promise<EventModel[]> {
    const localStart = this.getDateFromMonth(startMonth);
    const localEnd = this.getDateFromMonth(endMonth);
    localEnd.setMonth(localEnd.getMonth() + 1);
    localEnd.setDate(0);
    localEnd.setHours(23, 59, 59, 999);

    return (
      await fetchJson(`/api/calendars/${calendar}/events?start=${encodeURIComponent(localStart.toISOString())}&end=${encodeURIComponent(localEnd.toISOString())}`)
        .catch((err) => { throw err; })
    ).events;
  }

  async getEventsFromPreviouslyHiddenCalendar(calendar: string) {
    if (!browser) return;

    this.getEventsFromCalendar(calendar, this.previousMonth(this.eventsRangeStart), this.nextMonth(this.eventsRangeEnd)).catch((err) => {
      const calendarName = this.calendarsMap.get(this.calendarsCache.get(calendar)?.value?.find((cal) => cal === calendar) || "")?.name;
      queueNotification(
        ColorKeys.Danger,
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
    const isHidden = this.metadata.hiddenCalendars.has(newEvent.calendar)
    if (!isHidden && this.getMonthFromDate(newEvent.date.start) <= this.eventsRangeEnd && this.getMonthFromDate(newEvent.date.end) >= this.eventsRangeStart) this.events.push(newEvent);

    // info if the event is hidden
    if (isHidden) {
      queueNotification(
        ColorKeys.Accent,
        `Event ${newEvent.name} was added to a hidden calendar.`
      );
    }

    this.saveCache();
  };

  async editEvent(modifiedEvent: EventModel, changes: EventModelChanges, override: boolean): Promise<void> {
    if (!browser) return;

    // update in database
    if (modifiedEvent.date.allDay) {
      modifiedEvent.date.start.setHours(0, 0, 0, 0);
      modifiedEvent.date.end.setHours(0, 0, 0, 0);
    }

    if (override && changes !== NoChangesEvent) changes.color = true; // we have no way to destinguish between "don't change color" and "default color"
    const formData = this.getEventFormData(modifiedEvent, changes);
    formData.set("overridden", override ? "true" : "false");

    await fetchResponse(`/api/events/${modifiedEvent.id}`, { method: "PATCH", body: formData }).catch((err) => { throw err; });

    // update in cache
    modifiedEvent.overridden = override;
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
    if (this.getMonthFromDate(modifiedEvent.date.start) <= this.eventsRangeEnd && this.getMonthFromDate(modifiedEvent.date.end) >= this.eventsRangeStart) {
      //this.events.update((events) => events.map((event) => event.id === modifiedEvent.id ? modifiedEvent : event));
    } else {
      this.events.splice(this.events.findIndex((event) => event.id === modifiedEvent.id), 1);
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
    this.events.splice(this.events.findIndex((event) => event.id === id), 1);

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

export function getRepository(): Repository {
  return page.data.singletons.repository;
}
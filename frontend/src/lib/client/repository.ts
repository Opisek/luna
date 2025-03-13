import { browser } from "$app/environment";

import { writable } from "svelte/store";

import { hiddenCalendars, isCalendarVisible } from "./localStorage";
import { queueNotification } from "./notifications";
import { AllChangesEvent, AllChangesSource, NoOp } from "./placeholders";

import { atLeastOnePromise, deepCopy } from "$lib/common/misc";

//
// Constants
//

const spoolerDelay = 50; // 50ms
const maxCacheAge = 1000 * 60 * 10; // 10 minutes

//
// Subscribeable Stores
//

export const sources = writable([] as SourceModel[]);
export const calendars = writable([] as CalendarModel[]);
export const events = writable([] as EventModel[]);

export const faultySources = writable(new Set<string>());
export const faultyCalendars = writable(new Set<string>());

export const loadingSources = writable(new Set<string>());
export const loadingCalendars = writable(new Set<string>());
export const loadingData = writable(false);

//
// Misc
//

let loadingCounter = 0;
function indicateStartLoading() {
  loadingCounter++;
  loadingData.set(true);
}
function indicateStopLoading() {
  if (--loadingCounter == 0) loadingData.set(false);
}

async function mapBaseEventToRecurrenceInstance(idAndDate: [string, number]): Promise<EventModel | null> {
  const baseEvent = eventsMap.get(idAndDate[0]);
  if (baseEvent == undefined) return null;

  const baseDuration = baseEvent.date.end.getTime() - baseEvent.date.start.getTime();

  const copiedEvent = await deepCopy(baseEvent);
  copiedEvent.date.start = new Date(idAndDate[1]);
  copiedEvent.date.end = new Date(idAndDate[1]);
  copiedEvent.date.end.setTime(copiedEvent.date.end.getTime() + baseDuration);

  return copiedEvent;
}

async function mapAllRecurrenceInstances(idAndDates: [string, number][]): Promise<EventModel[]> {
  return (await Promise.all(idAndDates.map(async (idAndDate) => mapBaseEventToRecurrenceInstance(idAndDate)))).filter(x => x != null);
}

//
// Caching
//

let eventsRangeStart: Date = new Date();
let eventsRangeEnd: Date = new Date();

const emptyCache = { date: 0, value: null };

let lastCacheSave = Date.now();

let sourcesCache: CacheEntry<SourceModel[]> = emptyCache; // sources
let sourceDetailsCache: Map<string, CacheEntry<SourceModel>> = new Map(); // source -> details
let calendarsCache: Map<string, CacheEntry<string[]>> = new Map(); // source -> calendars
let eventsCache: Map<string, Map<number, CacheEntry<[string, number][]>>> = new Map(); // calendar -> month -> event id, start date (because of recurring events having the same id)
let eventsMap: Map<string, EventModel> = new Map(); // event id -> event
let calendarsMap: Map<string, CalendarModel> = new Map(); // calendar id -> calendar

function cacheOk<T>(cache: CacheEntry<T> | undefined): (T | null) {
  return (cache && Date.now() - cache.date < maxCacheAge) ? cache.value : null;
}

export function invalidateCache() {
  sourcesCache.date = 0;
  sourceDetailsCache.forEach((cache) => cache.date = 0);
  calendarsCache.forEach((cache) => cache.date = 0);
  eventsCache.forEach((cache) => cache.forEach((entry) => entry.date = 0));
  saveCache();
}

// 
// Web
// 

async function fetchResponse(url: string, options: RequestInit = {}): Promise<Response> {
  const response = await fetch(url, options).catch((err) => {
    if (!err) err = new Error("Could not contact server");
    throw err;
  });
  if (response.ok) {
    return response;
  } else {
    const json = await response.json().catch(() => null);
    let err = null;
    if (!err) err = json.error;
    if (!err) err = json.message;
    if (!err) err = `${response.statusText ? response.statusText : "Could not contact server"} (${response.status})`;
    throw new Error(err);
  }
}

async function fetchJson(url: string, options: RequestInit = {}) {
  return (await fetchResponse(url, options).catch(err => { throw err; })).json();
}

//
// Form Data
//

function getSourceFormData(source: SourceModel, changes: SourceModelChanges = AllChangesSource): FormData {
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

function getEventFormData(event: EventModel, changes: EventModelChanges = AllChangesEvent): FormData {
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

function compileSources() {
  sources.set(sourcesCache.value || []);
}

let compileCalendarsTimeout: ReturnType<typeof setTimeout>;
function compileCalendars() {
  clearTimeout(compileCalendarsTimeout);
  compileCalendarsTimeout = setTimeout(() => {
    const allCalendars = Array.from(calendarsCache.values().map(x => x.value).filter(x => x != null)).flat();
    calendars.set(allCalendars.map(x => calendarsMap.get(x)).filter(x => x != null));
  }, spoolerDelay)
}

let compileEventsTimeout: ReturnType<typeof setTimeout>;
function compileEvents(start: Date, end: Date) {
  clearTimeout(compileEventsTimeout);

  compileEventsTimeout = setTimeout(async () => {
    const allEvents = 
      Array.from(eventsCache.entries())
      .filter(x => isCalendarVisible(x[0]) && x[1] != null) // Event must be visible
      .map(x => Array.from(x[1].entries()))
      .flat()
      .filter(x => x[1] != null && x[0] >= start.getTime() && x[0] <= end.getTime()) // Event must be in the time frame
      .map(x => x[1].value)
      .filter(x => x != null) // Event must exist
      .flat();
    
    const uniqueEvents = [ ...new Map(allEvents.map(x => [`${x[0]}${x[1]}`, x])).values() ];

    const eventsWithData = await mapAllRecurrenceInstances(uniqueEvents);

    events.set(eventsWithData);
  }, spoolerDelay)
}

//
// Local Storage
//

function loadCache() {
  const cacheTimestamp = localStorage.getItem("cache.timestamp");
  if (cacheTimestamp != null && Number.parseInt(cacheTimestamp) == lastCacheSave) return;

  const newSourcesCache = localStorage.getItem("cache.sources");
  if (newSourcesCache) sourcesCache = JSON.parse(newSourcesCache);

  const newSourceDetailsCache = localStorage.getItem("cache.sourceDetails");
  if (newSourceDetailsCache) sourceDetailsCache = new Map(JSON.parse(newSourceDetailsCache));

  const newCalendarsCache = localStorage.getItem("cache.calendars");
  if (newCalendarsCache) calendarsCache = new Map(JSON.parse(newCalendarsCache));

  const newEventsCache = localStorage.getItem("cache.events");
  if (newEventsCache) eventsCache = new Map(JSON.parse(newEventsCache).map((x: [string, [number, CacheEntry<string>][]]) => [x[0], new Map(x[1])]));

  const newCalendarsMap = localStorage.getItem("cache.calendarsMap");
  if (newCalendarsMap) calendarsMap = new Map(JSON.parse(newCalendarsMap));

  const newEventsMap = localStorage.getItem("cache.eventsMap");
  if (newEventsMap) {
    eventsMap = new Map(JSON.parse(newEventsMap));
    eventsMap.forEach((event) => {
      event.date.start = new Date(event.date.start);
      event.date.end = new Date(event.date.end);
      if (event.date.allDay) {
        event.date.start.setHours(0, 0, 0, 0);
        event.date.end.setHours(0, 0, 0, 0);
      }
    });
  }

  compileSources();
  compileCalendars();
  compileEvents(eventsRangeStart, eventsRangeEnd);
}

if (browser) {
  window.addEventListener("storage", () => loadCache());
  loadCache();
}

let saveCacheTimeout: ReturnType<typeof setTimeout>;
function saveCache() {
  if (browser) {
    // @ts-ignore if only typescript were consistent...
    clearTimeout(saveCacheTimeout);
    setTimeout(() => {
      lastCacheSave = Date.now();
      localStorage.setItem("cache.timestamp", lastCacheSave.toString());
      localStorage.setItem("cache.sources", JSON.stringify(sourcesCache));
      localStorage.setItem("cache.sourceDetails", JSON.stringify(Array.from(sourceDetailsCache.entries())));
      localStorage.setItem("cache.calendars", JSON.stringify(Array.from(calendarsCache.entries())));
      localStorage.setItem("cache.events", JSON.stringify(Array.from(eventsCache.entries().map(x => [x[0], Array.from(x[1].entries())]))));
      localStorage.setItem("cache.calendarsMap", JSON.stringify(Array.from(calendarsMap.entries())));
      localStorage.setItem("cache.eventsMap", JSON.stringify(Array.from(eventsMap.entries())));
    }, spoolerDelay)
  }
}

// 
// Sources
// 

export async function getSources(forceRefresh = false): Promise<SourceModel[]> {
  if (!browser) return [];

  if (!forceRefresh) {
    const cached = cacheOk(sourcesCache);
    if (cached) return Promise.resolve(cached);
  }

  indicateStartLoading();

  const fetchedSources = await fetchJson("/api/sources").catch((err) => {
    throw err;
  }).finally(() => {
    indicateStopLoading();
  });

  sourcesCache.date = Date.now(),
  sourcesCache.value = fetchedSources;
  compileSources();
  saveCache();
  return fetchedSources;
}

export async function getSourceDetails(id: string, forceRefresh = false): Promise<SourceModel> {
  if (!browser) return {} as SourceModel;

  if (!forceRefresh) {
    const cached = cacheOk(sourceDetailsCache.get(id));
    if (cached) return Promise.resolve(cached);
  }

  const fetched = await fetchJson(`/api/sources/${id}`).catch((err) => { throw err; });

  sourceDetailsCache.set(id, {
    date: Date.now(),
    value: fetched
  });
  saveCache();
  return fetched;
}

export async function createSource(newSource: SourceModel): Promise<void> {
  if (!browser) return;

  const formData = getSourceFormData(newSource);

  const json = await fetchJson(`/api/sources`, { method: "PUT", body: formData }).catch((err) => { throw err; });

  newSource.id = json.id;
  sourcesCache.value = sourcesCache.value?.concat(newSource) || [ newSource ];
  sources.update((sources) => sources.concat(newSource));

  getCalendars(newSource.id).then(async (cals) => {
    compileCalendars();
    const [_, errors] = await atLeastOnePromise(cals.map((cal) => getEventsFromCalendar(cal.id, eventsRangeStart, eventsRangeEnd))).catch(() => {
      throw new Error("Failed to fetch events");
    });
    errors.forEach((err) => {
      queueNotification(
        "failure",
        `Failed to fetch events from ${cals[err[0]].name}: ${err[1].message}`
      );
    });
    compileEvents(eventsRangeStart, eventsRangeEnd);
  }).catch((err) => {
    queueNotification(
      "failure",
      `Failed to fetch calendars from ${newSource.name}: ${err.message}`
    );
  });

  saveCache();
}

export async function editSource(modifiedSource: SourceModel, changes: SourceModelChanges): Promise<void> {
  if (!browser) return;

  let formData = getSourceFormData(modifiedSource, changes);

  await fetchResponse(`/api/sources/${modifiedSource.id}`, { method: "PATCH", body: formData }).catch((err) => { throw err; });
  
  sourcesCache.value = sourcesCache.value?.map((source => source.id === modifiedSource.id ? modifiedSource : source)) || [ modifiedSource ];
  compileSources();

  getCalendars(modifiedSource.id).then(async (cals) => {
    compileCalendars();
    const [_, errors] = await atLeastOnePromise(cals.map((cal) => getEventsFromCalendar(cal.id, eventsRangeStart, eventsRangeEnd))).catch(() => {
      throw new Error("Failed to fetch events");
    });
    errors.forEach((err) => {
      queueNotification(
        "failure",
        `Failed to fetch events from ${cals[err[0]].name}: ${err[1].message}`
      );
    });
    compileEvents(eventsRangeStart, eventsRangeEnd);
  }).catch((err) => {
    queueNotification(
      "failure",
      `Failed to fetch calendars from ${modifiedSource.name}: ${err.message}`
    );
  });

  saveCache();
}

export async function deleteSource(id: string): Promise<void> {
  if (!browser) return;

  await fetchResponse(`/api/sources/${id}`, { method: "DELETE" }).catch((err) => { throw err; });

  sourcesCache.value = sourcesCache.value?.filter((source) => source.id !== id) || [];
  compileSources();

  for (const calendar of calendarsCache.get(id)?.value || []) {
    eventsCache.delete(calendar);
    calendarsMap.delete(calendar);
  }

  compileCalendars();
  compileEvents(eventsRangeStart, eventsRangeEnd);

  saveCache();
}

//
// Calendars
//

export async function getAllCalendars(forceRefresh = false): Promise<CalendarModel[]> {
  if (!browser) return [];

  const allSources = await getSources(forceRefresh);

  const [calendars, errors] = await atLeastOnePromise(allSources.map((source) => getCalendars(source.id, forceRefresh))).catch(() => {
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

async function getCalendars(id: string, forceRefresh = false): Promise<CalendarModel[]> {
  if (!browser) return [];

  if (!forceRefresh) {
    const cached = cacheOk(calendarsCache.get(id));
    if (cached) return Promise.resolve(cached.map(x => calendarsMap.get(x)).filter(x => x != null));
  }

  indicateStartLoading();
  loadingSources.update((loading) => {
    loading.add(id);
    return loading;
  });

  const response = await fetchJson(`/api/sources/${id}/calendars`).catch((err) => {
    faultySources.update((faulty) => {
      faulty.add(id);
      return faulty;
    });
    throw err;
  }).finally(() => {
    indicateStopLoading();
    loadingSources.update((loading) => {
      loading.delete(id);
      return loading;
    });
  });

  faultySources.update((faulty) => {
    faulty.delete(id);
    return faulty;
  });

  const fetched: CalendarModel[] = response.calendars;

  // Delete orphans
  let anyOrphaned = false;
  const fetchedSet = new Set(fetched.map(x => x.id));
  for (const calendar of calendarsCache.get(id)?.value || []) {
    if (!fetchedSet.has(calendar)) {
      eventsCache.delete(calendar);
      anyOrphaned = true;
    }
  }

  for (const calendar of fetched) {
    calendarsMap.set(calendar.id, calendar)
  }
  
  calendarsCache.set(id, {
    date: Date.now(),
    value: fetched.map(x => x.id)
  });

  compileCalendars();
  if (anyOrphaned) compileEvents(eventsRangeStart, eventsRangeEnd);
  saveCache();
  return fetched;
}

export async function getCalendar(id: string, forceRefresh = false): Promise<CalendarModel | null> {
  if (!browser) return {} as CalendarModel;

  return calendarsMap.get(id) || null; // TODO: needs a bit of refactoring plus actual fetch of the relevant endpoint depending on cache age
}

hiddenCalendars.subscribe(() => {
  compileEvents(eventsRangeStart, eventsRangeEnd);
});

//
// Events
//

function determineEventMonths(event: EventModel): Date[] {
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

function addEventToCache(event: EventModel, date: Date) {
  let calendarEventsCache = eventsCache.get(event.calendar);
  if (!calendarEventsCache) {
    calendarEventsCache = new Map();
    eventsCache.set(event.calendar, calendarEventsCache);
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

function removeEventFromCache(event: EventModel, date: Date) {
  if (!eventsCache.has(event.calendar)) return;
  const cacheEntry = eventsCache.get(event.calendar)?.get(date.getTime());
  if (!cacheEntry || !cacheEntry.value) return;
  cacheEntry.value = cacheEntry.value.filter((idAndDate) => idAndDate[0] !== event.id);
}

export async function getAllEvents(start: Date, end: Date, forceRefresh = false): Promise<EventModel[]> {
  if (!browser) return [];
  eventsRangeStart = start;
  eventsRangeEnd = end;

  start.setUTCHours(0, 0, 0, 0);
  end.setUTCHours(23, 59, 59, 999);

  // Set start and end to the start and end of each month
  start.setDate(1);
  end.setMonth(end.getMonth() + 1);
  end.setDate(0);

  // Add one month of padding in both directions
  start.setMonth(start.getMonth() - 1);
  end.setMonth(end.getMonth() + 1);

  compileEvents(start, end);
  const allSources = await getSources(forceRefresh).catch((err) => { throw err; });
  const [events, errors] = await atLeastOnePromise(allSources.map((source) => getEventsFromSource(source.id, start, end, forceRefresh))).catch(() => {
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

async function getEventsFromSource(source: string, start: Date, end: Date, forceRefresh = false): Promise<EventModel[]> {
  let cals = await getCalendars(source, forceRefresh).catch((err) => { throw err; });
  cals = cals.filter(x => isCalendarVisible(x.id)); // only fetch events from visible calendars
  const [events, errors] = await atLeastOnePromise(cals.map((calendar) => getEventsFromCalendar(calendar.id, start, end, forceRefresh))).catch(() => {
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

async function getEventsFromCalendar(calendar: string, start: Date, end: Date, forceRefresh = false): Promise<EventModel[]> {
  let result: EventModel[] = [];

  const cache = (forceRefresh ? null : eventsCache.get(calendar)) || new Map<number, CacheEntry<[string, number][]>>();

  // Limit the range we need to ask for but only use one request per calendar
  const fetchStart = new Date(start);
  const fetchEnd = new Date(end);

  while (fetchStart.getTime() <= fetchEnd.getTime()) {
    const cached = cacheOk(cache.get(fetchStart.getTime()));
    if (!cached) break;
    result = result.concat(await mapAllRecurrenceInstances(cached));
    fetchStart.setMonth(fetchStart.getMonth() + 1);
  }

  while (fetchEnd.getTime() >= fetchStart.getTime()) {
    const cached = cacheOk(cache.get(fetchEnd.getTime()));
    if (!cached) break;
    result = result.concat(await mapAllRecurrenceInstances(cached));
    fetchEnd.setMonth(fetchEnd.getMonth() - 1);
  }

  if (fetchStart.getTime() > fetchEnd.getTime()) {
    return result;
  }

  indicateStartLoading();
  loadingCalendars.update((loading) => {
    loading.add(calendar);
    return loading;
  });

  const fetchedEvents = await fetchEvents(calendar, start, end).catch((err) => {
    faultyCalendars.update((faulty) => {
      faulty.add(calendar);
      return faulty;
    });

    throw err;
  }).finally(() => {
    indicateStopLoading();
    loadingCalendars.update((loading) => {
      loading.delete(calendar);
      return loading;
    });
  });

  faultyCalendars.update((faulty) => {
    faulty.delete(calendar);
    return faulty;
  });

  compileEvents(start, end);
  saveCache();
  return result.concat(fetchedEvents);
}

async function fetchEvents(calendar: string, start: Date, end: Date): Promise<EventModel[]> {
  const localStart = new Date(start);
  localStart.setHours(0, 0, 0, 0);
  const localEnd = new Date(end);
  localEnd.setHours(23, 59, 59, 999);

  const fetched: EventModel[] = (await fetchJson(`/api/calendars/${calendar}/events?start=${encodeURIComponent(localStart.toISOString())}&end=${encodeURIComponent(localEnd.toISOString())}`)).events;

  let calendarEventsCache = eventsCache.get(calendar);
  if (!calendarEventsCache) {
    calendarEventsCache = new Map();
    eventsCache.set(calendar, calendarEventsCache);
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

    eventsMap.set(event.id, event);
    for (const month of determineEventMonths(event)) {
      addEventToCache(event, month);
    }
  }

  return fetched;
}

export async function getEventsFromPreviouslyHiddenCalendar(calendar: string) {
  if (!browser) return;

  getEventsFromCalendar(calendar, eventsRangeStart, eventsRangeEnd).catch((err) => {
    const calendarName = calendarsMap.get(calendarsCache.get(calendar)?.value?.find((cal) => cal === calendar) || "")?.name;
    queueNotification(
      "failure",
      `Failed to fetch events from calendar${calendarName ? " " + calendarName : ""}: ${err.message}`
    );
  });
}

export async function createEvent(newEvent: EventModel): Promise<void> {
  if (!browser) return;

  // add to database
  if (newEvent.date.allDay) {
    newEvent.date.start.setHours(0, 0, 0, 0);
    newEvent.date.end.setHours(0, 0, 0, 0);
  }

  const formData = getEventFormData(newEvent);

  const json = await fetchJson(`/api/calendars/${newEvent.calendar}/events`, { method: "PUT", body: formData }).catch((err) => { throw err; });

  newEvent.id = json.id;

  // add to cache
  eventsMap.set(newEvent.id, newEvent);
  for (const month of determineEventMonths(newEvent)) addEventToCache(newEvent, month);

  // add to display
  if (isCalendarVisible(newEvent.calendar) && newEvent.date.start <= eventsRangeEnd && newEvent.date.end >= eventsRangeStart) events.update((events) => events.concat(newEvent));
};

export async function editEvent(modifiedEvent: EventModel, changes: EventModelChanges): Promise<void> {
  if (!browser) return;

  // update in database
  if (modifiedEvent.date.allDay) {
    modifiedEvent.date.start.setHours(0, 0, 0, 0);
    modifiedEvent.date.end.setHours(0, 0, 0, 0);
  }

  const formData = getEventFormData(modifiedEvent, changes);

  await fetchResponse(`/api/events/${modifiedEvent.id}`, { method: "PATCH", body: formData }).catch((err) => { throw err; });

  // update in cache
  const previousMonths = determineEventMonths(eventsMap.get(modifiedEvent.id)!);
  const currentMonths = determineEventMonths(modifiedEvent);
  eventsMap.set(modifiedEvent.id, modifiedEvent);

  for (const month of previousMonths) {
    if (!currentMonths.includes(month)) {
      removeEventFromCache(modifiedEvent, month);
    }
  }

  for (const month of currentMonths) {
    if (!previousMonths.includes(month)) {
      addEventToCache(modifiedEvent, month);
    }
  }

  // update on display
  if (modifiedEvent.date.start <= eventsRangeEnd && modifiedEvent.date.end >= eventsRangeStart) {
    events.update((events) => events.map((event) => event.id === modifiedEvent.id ? modifiedEvent : event));
  } else {
    events.update((events) => events.filter((event) => event.id !== modifiedEvent.id));
  }

  saveCache();
}

export async function deleteEvent(id: string): Promise<void> {
  if (!browser) return;

  // remove from database
  await fetchResponse(`/api/events/${id}`, { method: "DELETE" }).catch((err) => { throw err; });

  const event = eventsMap.get(id);
  if (!event) return;
  eventsMap.delete(id);

  // remove from cache
  const months = determineEventMonths(event);
  for (const month of months) removeEventFromCache(event, month);

  // remove from display
  events.update((events) => events.filter((event) => event.id !== id));

  saveCache();
}

export async function moveEvent(event: EventModel): Promise<void> {
  if (!browser) return;

  const oldId = event.id;

  // add to the new calendar
  await createEvent(event).catch((err) => { throw err; });

  // remove from the old calendar
  await deleteEvent(oldId).catch((err) => {
    // undo changes
    deleteEvent(event.id).catch(NoOp);
    event.id = oldId;
    throw err;
  });
}
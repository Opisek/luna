import { writable } from "svelte/store";

import { hiddenCalendars } from "./localStorage";
import { isCalendarVisible, isSourceCollapsed } from "./localStorage";
import { queueNotification } from "./notifications";

// TODO: local storage integration for PWA offline support (longterm goal)

export const sources = writable([] as SourceModel[]);
export const calendars = writable([] as CalendarModel[]);
export const events = writable([] as EventModel[]);

let sourceCalendars = new Map<string, Set<CalendarModel>>();
let calendarEvents = new Map<string, Set<EventModel>>();

let calendarMap = new Map<string, CalendarModel>();
let eventsMap = new Map<string, EventModel>();

export const faultySources = writable(new Set<string>());
export const faultyCalendars = writable(new Set<string>());

// TODO: this will depend on the month that the user is currently viewing (stored in session storaage)
let lastStart = new Date();
lastStart.setMonth(lastStart.getMonth() - 1);
lastStart.setDate(0);
let lastEnd = new Date();
lastEnd.setMonth(lastEnd.getMonth() + 2);
lastEnd.setDate(0);

function allEvents(): EventModel[] {
  const allEvents = Array.from(
    calendarEvents
      .entries()
      .filter(x => isCalendarVisible(x[0]))
      .map(x => Array.from(x[1])
    )
  ).flat();
  eventsMap = new Map(allEvents.map(event => [event.id, event]));
  return allEvents;
}

function allCalendars(): CalendarModel[] {
  const allCalendars = Array.from(sourceCalendars.values().map(x => Array.from(x))).flat();
  calendarMap = new Map(allCalendars.map(calendar => [calendar.id, calendar]));
  return allCalendars;
}

function compileEvents() {
  events.set(allEvents());
}

function compileCalendars() {
  calendars.set(allCalendars());
}

export const fetchSources = async (): Promise<void> => {
  const response = await fetch("/api/sources");
  if (response.ok) {
    const fetchedSources = await response.json() as SourceModel[];

    for (const source of fetchedSources) {
      source.collapsed = isSourceCollapsed(source.id);

      fetchCalendars(source.id).catch(err => {
        queueNotification(
          "failure",
          `Failed to fetch calendars from source "${source.name}":\n${err.message}`
        );
      });
    }

    sources.set(fetchedSources);
  } else {
    const json = await response.json();
    throw new Error(json ? json.error : "Could not contact the server");
  }
};

function getSourceFormData(source: SourceModel): FormData {
  const formData = new FormData();
  formData.set("name", source.name);
  formData.set("type", source.type);
  switch (source.type) {
    case "caldav":
      formData.set("url", source.settings.url);
      break;
    default:
      throw new Error("Unsupported source type");
  }
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
  return formData;
}

export const createSource = async (newSource: SourceModel): Promise<void> => {
  let formData: FormData;
  try {
    formData = getSourceFormData(newSource);
  } catch (e: any) {
    throw e;
  }

  const response = await fetch(`/api/sources`, { method: "PUT", body: formData });
  if (response.ok) {
    const json = await response.json();
    newSource.id = json.id;
    sources.update((sources) => sources.concat(newSource));

    fetchCalendars(newSource.id).catch(err => {
      queueNotification(
        "failure",
        `Failed to fetch calendar: ${err.message}`
      );
    });
  } else {
    const json = await response.json();
    throw new Error(json ? json.error : "Could not contact the server");
  }
}

export const editSource = async (modifiedSource: SourceModel): Promise<void> => {
  let formData = getSourceFormData(modifiedSource);

  const response = await fetch(`/api/sources/${modifiedSource.id}`, { method: "PATCH", body: formData });
  if (response.ok) {
    sources.update((sources) => sources.map((source => source.id === modifiedSource.id ? modifiedSource : source)))

    fetchCalendars(modifiedSource.id).catch(err => {
      queueNotification(
        "failure",
        `Failed to fetch calendar: ${err.message}`
      );
    });
  } else {
    const json = await response.json();
    throw new Error(json ? json.error : "Could not contact the server");
  }
}

export const deleteSource = async (id: string): Promise<void> => {
  const response = await fetch(`/api/sources/${id}`, { method: "DELETE" });
  if (response.ok) {
    sources.update((sources) => sources.filter((source) => source.id !== id));
    sourceCalendars.get(id)?.forEach((calendar) => {
      calendarEvents.delete(calendar.id);
    });
    sourceCalendars.delete(id);
    compileEvents();
  } else {
    const json = await response.json();
    throw new Error(json ? json.error : "Could not contact the server");
  }
}

export const fetchCalendars = async (id: string): Promise<void> => {
  const response = await fetch(`/api/sources/${id}/calendars`);
  if (response.ok) {
    faultySources.update((faultySources) => new Set([...faultySources].filter((faultySource) => faultySource !== id)));
    const json = await response.json() as {calendars: CalendarModel[]};
    const fetchCalendars = json.calendars;

    for (const calendar of fetchCalendars) {
      //calendar.visible = isCalendarVisible(calendar.id);

      fetchEvents(calendar.id, lastStart, lastEnd).catch(err => {
        queueNotification(
          "failure",
          `Failed to fetch events: ${err.message}`
        );
      });
    }

    sourceCalendars.set(id, new Set(fetchCalendars));
    compileCalendars();
  } else {
    faultySources.update((faultySources) => new Set(faultySources.add(id)));
    sourceCalendars.delete(id);
    compileCalendars();

    const json = await response.json();
    throw new Error(json ? json.error : "Could not contact the server");
  }
};

export const getCalendars = () => allCalendars();

export const fetchAllEvents = async (start: Date, end: Date) => {
  lastStart = start;
  lastEnd = end;
  for (const calendar of allCalendars()) {
    fetchEvents(calendar.id, start, end).catch(err => {
      queueNotification(
        "failure",
        `Failed to fetch events: ${err.message}`
      );
    });
  }
}

export const fetchSourceCalendars = async (id: string) => {
  return Array.from(sourceCalendars.get(id) || new Set<CalendarModel>());
}

export const fetchEvents = async (id: string, start: Date, end: Date): Promise<void> => {
  const url = `/api/calendars/${id}/events?start=${encodeURIComponent(start.toISOString())}&end=${encodeURIComponent(end.toISOString())}`
  const response = await fetch(url);
  if (response.ok) {
    faultyCalendars.update((faultyCalendars) => new Set([...faultyCalendars].filter((faultyCalendar) => faultyCalendar !== id)));
    const json = await response.json() as {events: EventModel[]};
    for (const event of json.events) {
      event.date.start = new Date(event.date.start);
      event.date.end = new Date(event.date.end);

      if (event.date.allDay) {
        event.date.start.setHours(0, 0, 0, 0);
        event.date.end.setHours(0, 0, 0, 0);
      }
    }

    // Do not remove events outside the requested range
    const oldEvents = Array.from(calendarEvents.get(id) || new Set<EventModel>()).filter((event) => event.date.end < start || event.date.start > end);
    calendarEvents.set(id, new Set(json.events.concat(oldEvents)));
    compileEvents();
  } else {
    faultyCalendars.update((faultyCalendars) => new Set(faultyCalendars.add(id)));
    calendarEvents.delete(id);
    compileEvents();

    const json = await response.json();
    throw new Error(json ? json.error : "Could not contact the server");
  }
};

function getEventFormData(event: EventModel): FormData {
  const formData = new FormData();
  formData.set("name", event.name);
  formData.set("desc", event.desc);
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
  if (event.color && event.color !== "") {
    formData.set("color", event.color);
  } else {
    formData.set("color", "null");
  }
  return formData;
}

export const createEvent = async (newEvent: EventModel): Promise<void> => {
  if (newEvent.date.allDay) {
    newEvent.date.start.setHours(0, 0, 0, 0);
    newEvent.date.end.setHours(0, 0, 0, 0);
  }

  const formData = getEventFormData(newEvent);

  const response = await fetch(`/api/calendars/${newEvent.calendar}/events`, { method: "PUT", body: formData });
  if (response.ok) {
    const json = await response.json();
    newEvent.id = json.id;

    calendarEvents.set(newEvent.calendar, new Set([...calendarEvents.get(newEvent.calendar) || [], newEvent]));
    compileEvents();
  } else {
    const json = await response.json();
    throw new Error(json ? json.error : "Could not contact the server");
  }
};

export const editEvent = async (modifiedEvent: EventModel): Promise<void> => {
  if (modifiedEvent.date.allDay) {
    modifiedEvent.date.start.setHours(0, 0, 0, 0);
    modifiedEvent.date.end.setHours(0, 0, 0, 0);
  }

  const formData = getEventFormData(modifiedEvent);

  const response = await fetch(`/api/events/${modifiedEvent.id}`, { method: "PATCH", body: formData });
  if (response.ok) {
    calendarEvents.set(modifiedEvent.calendar, new Set([...calendarEvents.get(modifiedEvent.calendar) || []].map((event) => event.id === modifiedEvent.id ? modifiedEvent : event)));
    compileEvents();
  } else {
    const json = await response.json();
    throw new Error(json ? json.error : "Could not contact the server");
  }
}

export const deleteEvent = async (id: string): Promise<void> => {
  const response = await fetch(`/api/events/${id}`, { method: "DELETE" });
  if (response.ok) {
    const event = eventsMap.get(id);
    const calendarId = event?.calendar;
    if (calendarId) {
      calendarEvents.set(calendarId, new Set([...calendarEvents.get(calendarId) || []].filter((event) => event.id !== id)));
      compileEvents();
    }
  } else {
    const json = await response.json();
    throw new Error(json ? json.error : "Could not contact the server");
  }
}

export const recalculateEventVisibility = () => {
  compileEvents();
}

hiddenCalendars.subscribe(() => {
  if (calendarEvents.size === 0) return;
  recalculateEventVisibility();
});
import { writable } from "svelte/store";
import { queueNotification } from "./notifications";
import { isCalendarVisible, isSourceCollapsed } from "./localStorage";

// TODO: local storage integration for PWA offline support (longterm goal)

let localCalendars: CalendarModel[] = [];

export const sources = writable([] as SourceModel[]);
export const calendars = writable(localCalendars);
export const events = writable([] as EventModel[]);

export const sourceCalendars = writable(new Map<string, CalendarModel[]>());
export const calendarEvents = writable(new Map<string, EventModel[]>());

export const faultySources = writable(new Set<string>());
export const faultyCalendars = writable(new Set<string>());

sourceCalendars.subscribe((sourceCalendars) => {
  localCalendars = Array.from(sourceCalendars.values()).flat();
  calendars.set(localCalendars);
});

calendarEvents.subscribe((calendarEvents) => {
  events.set(Array.from(calendarEvents.values()).flat());
});

export const fetchSources = async (): Promise<string> => {
  try {
    const response = await fetch("/api/sources");
    if (response.ok) {
      const fetchedSources = await response.json() as SourceModel[];

      for (const source of fetchedSources) {
        source.collapsed = isSourceCollapsed(source.id);

        fetchCalendars(source.id).then(err => {
          if (err != "") {
            queueNotification(
              "failure",
              `Failed to fetch calendar: ${err}`
            );
          }
        });
      }

     sources.set(fetchedSources);

      return ""
    } else {
      const json = await response.json();
      return (json ? json.error : "Could not contact the server");
    }
  } catch (e) {
    return "Unexpected error occured"
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

export const createSource = async (newSource: SourceModel): Promise<string> => {
  try {
    let formData: FormData;
    try {
      formData = getSourceFormData(newSource);
    } catch (e: any) {
      return e.message;
    }

    const response = await fetch(`/api/sources`, { method: "PUT", body: formData });
    if (response.ok) {
      const json = await response.json();
      newSource.id = json.id;
      sources.update((sources) => sources.concat(newSource));

      fetchCalendars(newSource.id).then(err => {
        if (err != "") {
          queueNotification(
            "failure",
            `Failed to fetch calendar: ${err}`
          );
        }
      });

      return "";
    } else {
      const json = await response.json();
      return (json ? json.error : "Could not contact the server");
    }
  } catch (e) {
    return "Unexpected error occured"
  }
}

export const editSource = async (modifiedSource: SourceModel): Promise<string> => {
  try {
    let formData: FormData;
    try {
      formData = getSourceFormData(modifiedSource);
    } catch (e: any) {
      return e.message;
    }

    const response = await fetch(`/api/sources/${modifiedSource.id}`, { method: "PATCH", body: formData });
    if (response.ok) {
      sources.update((sources) => sources.map((source => source.id === modifiedSource.id ? modifiedSource : source)))

      fetchCalendars(modifiedSource.id).then(err => {
        if (err != "") {
          queueNotification(
            "failure",
            `Failed to fetch calendar: ${err}`
          );
        }
      });

      return "";
    } else {
      const json = await response.json();
      return (json ? json.error : "Could not contact the server");
    }
  } catch (e) {
    return "Unexpected error occured"
  }
}

export const deleteSource = async (id: string): Promise<string> => {
  try {
    const response = await fetch(`/api/sources/${id}`, { method: "DELETE" });
    if (response.ok) {
      sources.update((sources) => sources.filter((source) => source.id !== id));
      // TODO: also remove the events from those calendars but i'll see how i do that during the rewrite/refactor of the whole GET logic
      sourceCalendars.update((sourceCalendars) => { sourceCalendars.delete(id); return sourceCalendars; });
      return "";
    } else {
      const json = await response.json();
      return (json ? json.error : "Could not contact the server");
    }
  } catch (e) {
    return "Unexpected error occured"
  }
}

export const fetchCalendars = async (id: string): Promise<string> => {
  try {
    const response = await fetch(`/api/sources/${id}/calendars`);
    if (response.ok) {
      faultySources.update((faultySources) => new Set([...faultySources].filter((faultySource) => faultySource !== id)));
      const json = await response.json() as {calendars: CalendarModel[]};
      const fetchCalendars = json.calendars;

      for (const calendar of fetchCalendars) {
        calendar.visible = isCalendarVisible(calendar.id);

        fetchEvents(calendar.id).then(err => {
          if (err != "") {
            queueNotification(
              "failure",
              `Failed to fetch events: ${err}`
            );
          }
        });
      }

      sourceCalendars.update((sourceCalendars) => sourceCalendars.set(id, fetchCalendars));

      return ""
    } else {
      faultySources.update((faultySources) => new Set(faultySources.add(id)));
      sourceCalendars.update((sourceCalendars) => { sourceCalendars.delete(id); return sourceCalendars; });
      const json = await response.json();
      return (json ? json.error : "Could not contact the server");
    }
  } catch (e) {
    return "Unexpected error occured"
  }
};

export const getCalendars = () => localCalendars;

export const fetchEvents = async (id: string): Promise<string> => {
  try {
    const response = await fetch(`/api/calendars/${id}/events`);
    if (response.ok) {
      faultyCalendars.update((faultyCalendars) => new Set([...faultyCalendars].filter((faultyCalendar) => faultyCalendar !== id)));
      const json = await response.json() as {events: EventModel[]};
      for (const event of json.events) {
        event.date.start = new Date(event.date.start);
        event.date.end = new Date(event.date.end);
      }
      calendarEvents.update((calendarEvents) => calendarEvents.set(id, json.events));
      return ""
    } else {
      faultyCalendars.update((faultyCalendars) => new Set(faultyCalendars.add(id)));
      calendarEvents.update((calendarEvents) => { calendarEvents.delete(id); return calendarEvents; });
      const json = await response.json();
      return (json ? json.error : "Could not contact the server");
    }
  } catch (e) {
    return "Unexpected error occured"
  }
};

function getEventFormData(event: EventModel): FormData {
  const formData = new FormData();
  formData.set("name", event.name);
  formData.set("desc", event.desc);
  formData.set("date_start", event.date.start.toISOString());
  formData.set("date_end", event.date.end.toISOString());
  return formData;
}

export const createEvent = async (newEvent: EventModel): Promise<string> => {
  try {
    const formData = getEventFormData(newEvent);
    // TODO: add color picker to the frontend
    formData.set("color", "#FF0000")

    const response = await fetch(`/api/calendars/${newEvent.calendar}/events/`, { method: "PUT", body: formData });
    if (response.ok) {
      const json = await response.json();
      newEvent.id = json.id;
      events.update((events) => events.concat(newEvent));

      return "";
    } else {
      const json = await response.json();
      return (json ? json.error : "Could not contact the server");
    }
  } catch (e) {
    return "Unexpected error occured"
  }
};

export const editEvent = async (modifiedEvent: EventModel): Promise<string> => {
  try {
    const formData = getEventFormData(modifiedEvent);

    const response = await fetch(`/api/events/${modifiedEvent.id}`, { method: "PATCH", body: formData });
    if (response.ok) {
      events.update((events) => events.map((event => event.id === modifiedEvent.id ? modifiedEvent : event)))
      return "";
    } else {
      const json = await response.json();
      return (json ? json.error : "Could not contact the server");
    }
  } catch (e) {
    return "Unexpected error occured"
  }
}

export const deleteEvent = async (id: string): Promise<string> => {
  try {
    const response = await fetch(`/api/events/${id}`, { method: "DELETE" });
    if (response.ok) {
      events.update((events) => events.filter((event) => event.id !== id));
      return "";
    } else {
      const json = await response.json();
      return (json ? json.error : "Could not contact the server");
    }
  } catch (e) {
    return "Unexpected error occured"
  }
}
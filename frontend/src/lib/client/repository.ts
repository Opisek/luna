import { writable } from "svelte/store";

// TODO: local storage integration
export const sources = writable([] as SourceModel[]);
export const calendars = writable([] as CalendarModel[]);
export const events = writable([] as EventModel[]);

export const faultySources = writable(new Set<string>());
export const faultyCalendars = writable(new Set<string>());

export const fetchSources = async (): Promise<string> => {
  try {
    const response = await fetch("/api/sources");
    if (response.ok) {
      sources.set(await response.json());
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
      newSource.id = await response.text();
      sources.update((sources) => sources.map((source => source.id === newSource.id ? newSource : source)))
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
      return "";
    } else {
      const json = await response.json();
      return (json ? json.error : "Could not contact the server");
    }
  } catch (e) {
    return "Unexpected error occured"
  }
}

export const fetchCalendars = async (): Promise<string> => {
  try {
    const response = await fetch("/api/calendars");
    if (response.ok) {
      const json = await response.json() as {calendars: CalendarModel[], errored: string[]};
      calendars.set(json.calendars);
      faultySources.set(new Set<string>(json.errored));
      return ""
    } else {
      const json = await response.json();
      return (json ? json.error : "Could not contact the server");
    }
  } catch (e) {
    return "Unexpected error occured"
  }
};

export const fetchEvents = async (): Promise<string> => {
  try {
    const response = await fetch("/api/events");
    if (response.ok) {
      const json = await response.json() as {events: EventModel[], errored: string[]};
      for (const event of json.events) {
        event.date.start = new Date(event.date.start);
        event.date.end = new Date(event.date.end);
      }
      events.set(json.events);
      faultyCalendars.set(new Set<string>(json.errored));
      return ""
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
    const formData = new FormData();
    formData.set("name", modifiedEvent.name);
    formData.set("desc", modifiedEvent.desc);
    formData.set("date_start", modifiedEvent.date.start.toString());
    formData.set("date_end", modifiedEvent.date.end.toString());

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
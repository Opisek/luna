import { writable } from "svelte/store";

// TODO: local storage integration
export const sources = writable([] as SourceModel[]);
export const calendars = writable([] as CalendarModel[]);
export const events = writable([] as EventModel[]);

export const fetchSources = async (): Promise<string> => {
  const response = await fetch("/api/sources");
  if (response.ok) {
    sources.set(await response.json());
    return ""
  } else {
    const json = await response.json();
    return (json ? json.error : "Could not contact the server");
  }
};

export const fetchCalendars = async (): Promise<string> => {
  const response = await fetch("/api/calendars");
  if (response.ok) {
    calendars.set(await response.json());
    return ""
  } else {
    const json = await response.json();
    return (json ? json.error : "Could not contact the server");
  }
};

export const fetchEvents = async (): Promise<string> => {
  const response = await fetch("/api/events");
  if (response.ok) {
    const fetched = await response.json()
    for (const event of fetched) {
      event.date.start = new Date(event.date.start);
      event.date.end = new Date(event.date.end);
    }
    events.set(fetched);
    return ""
  } else {
    const json = await response.json();
    return (json ? json.error : "Could not contact the server");
  }
};

export const deleteEvent = async (id: string): Promise<string> => {
  const response = await fetch(`/api/events/${id}`, { method: "DELETE" });
  if (response.ok) {
    events.update((events) => events.filter((event) => event.id !== id));
    return ""
  } else {
    const json = await response.json();
    return (json ? json.error : "Could not contact the server");
  }
}
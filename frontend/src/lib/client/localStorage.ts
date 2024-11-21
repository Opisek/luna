import { browser } from "$app/environment";

import { writable } from "svelte/store";

let localCollapsedSources = new Set<string>(browser && JSON.parse(localStorage.getItem("collapsedSources") || "[]") || []);
export const collapsedSources = writable(localCollapsedSources);

let localHiddenCalendars = new Set<string>(browser && JSON.parse(localStorage.getItem("hiddenCalendars") || "[]") || []);
export const hiddenCalendars = writable(localHiddenCalendars);

collapsedSources.subscribe(value => {
  if (browser) localStorage.setItem("collapsedSources", JSON.stringify(Array.from(value)));
});
hiddenCalendars.subscribe(value => {
  if (browser) localStorage.setItem("hiddenCalendars", JSON.stringify(Array.from(value)));
});

if (browser) {
  window.addEventListener("storage", () => {
    const newCollapsedSources = localStorage.getItem("collapsedSources");
    if (newCollapsedSources) {
      const set = new Set<string>(JSON.parse(newCollapsedSources) as string[]);
      localCollapsedSources = set;
      collapsedSources.set(set);
    }

    const newHiddenCalendars = localStorage.getItem("hiddenCalendars");
    if (newHiddenCalendars) {
      const set = new Set<string>(JSON.parse(newHiddenCalendars) as string[]);
      localHiddenCalendars = set;
      hiddenCalendars.set(set);
    }
  })
}

export const setSourceCollapse = (sourceId: string, collapsed: boolean) => {
  if (isSourceCollapsed(sourceId) === collapsed) return;

  const set = localCollapsedSources;
  if (collapsed) {
    set.add(sourceId);
  } else {
    set.delete(sourceId);
  }
  collapsedSources.set(set);
}

export const setCalendarVisibility = (calendarId: string, visible: boolean) => {
  if (isCalendarVisible(calendarId) === visible) return;

  const set = localHiddenCalendars;
  if (!visible) {
    set.add(calendarId);
  } else {
    set.delete(calendarId);
  }
  hiddenCalendars.set(set);
}

export const isSourceCollapsed = (sourceId: string) => localCollapsedSources.has(sourceId);
export const isCalendarVisible = (calendarId: string) => !localHiddenCalendars.has(calendarId);
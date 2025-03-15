import { browser } from "$app/environment";

import { writable, type Writable } from "svelte/store";

import { SubscribeableSet } from "./reactivity";

class Metadata {
  //
  // Exported Subscribeables
  //

  // Faults
  readonly faultySources: SubscribeableSet<string>;
  readonly faultyCalendars: SubscribeableSet<string>;

  // Loading
  readonly loadingSources: SubscribeableSet<string>;
  readonly loadingCalendars: SubscribeableSet<string>;
  readonly loadingData: Writable<boolean>;
  private loadingCounter;

  // Hidden / Collapsed
  readonly collapsedSources: SubscribeableSet<string>;
  readonly hiddenCalendars: SubscribeableSet<string>;

  //
  // Constructor
  //
  constructor() {
    this.faultySources = new SubscribeableSet();
    this.faultyCalendars = new SubscribeableSet();

    this.loadingSources = new SubscribeableSet(true);
    this.loadingCalendars = new SubscribeableSet(true);
    this.loadingData = writable(false);
    this.loadingCounter = 0;

    if (browser) {
      this.collapsedSources = new SubscribeableSet(false, JSON.parse(localStorage.getItem("collapsedSources") || "[]"));
      this.hiddenCalendars = new SubscribeableSet(false, JSON.parse(localStorage.getItem("hiddenCalendars") || "[]"));
    } else {
      this.collapsedSources = new SubscribeableSet();
      this.hiddenCalendars = new SubscribeableSet();
    }

    this.collapsedSources.subscribe(value => {
      if (browser) localStorage.setItem("collapsedSources", JSON.stringify(Array.from(value)));
    });
    this.hiddenCalendars.subscribe(value => {
      if (browser) localStorage.setItem("hiddenCalendars", JSON.stringify(Array.from(value)));
    });

    if (browser) {
      window.addEventListener("storage", () => {
        const newCollapsedSources = localStorage.getItem("collapsedSources");
        if (newCollapsedSources) {
          const set = new Set<string>(JSON.parse(newCollapsedSources) as string[]);
          this.collapsedSources.set(set);
        }

        const newHiddenCalendars = localStorage.getItem("hiddenCalendars");
        if (newHiddenCalendars) {
          const set = new Set<string>(JSON.parse(newHiddenCalendars) as string[]);
          this.hiddenCalendars.set(set);
        }
      })
    }
  }

  //
  // Logic
  //

  // Loading
  startLoading(): (() => void) {
    this.loadingCounter++;
    this.loadingData.set(true);

    let called = false;
    return (() => {
      if (!called) {
        called = true;
        if (--this.loadingCounter == 0) this.loadingData.set(false);
      }
    });
  }

  startLoadingSource(source: string): (() => void) {
    this.loadingSources.add(source);
    const stopLoadingInternal = this.startLoading();

    let called = false;
    return (() => {
      if (!called) {
        called = true;
        this.loadingSources.delete(source);
        stopLoadingInternal();
      }
    });
  }

  startLoadingCalendar(calendar: string): (() => void) {
    this.loadingCalendars.add(calendar);
    const stopLoadingInternal = this.startLoading();
    // TODO: start loading the calendar's source as well

    let called = false;
    return (() => {
      if (!called) {
        called = true;
        this.loadingCalendars.delete(calendar);
        stopLoadingInternal();
      }
    });
  }

  // Faults
  addFaultySource(source: string, fault: string) {
    this.faultySources.add(source);
  }

  removeFaultySource(source: string) {
    this.faultySources.delete(source);
  }

  addFaultyCalendar(calendar: string, fault: string) {
    this.faultyCalendars.add(calendar);
  }

  removeFaultyCalendar(calendar: string) {
    this.faultyCalendars.delete(calendar);
  }

  // Hidden / Collapsed
  setSourceCollapse = (sourceId: string, collapsed: boolean) => {
    if (collapsed) this.collapsedSources.add(sourceId);
    else this.collapsedSources.delete(sourceId);
  }

  setCalendarVisibility = (calendarId: string, visible: boolean) => {
    if (visible) this.hiddenCalendars.delete(calendarId);
    else this.hiddenCalendars.add(calendarId);
  }
}

let metadata: Metadata | null = null;
export function getMetadata() {
  if (metadata === null) {
    metadata = new Metadata();
  }
  return metadata;
}
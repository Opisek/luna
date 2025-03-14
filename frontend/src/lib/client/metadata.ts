import { browser } from "$app/environment";

import { writable, type Writable } from "svelte/store";

class SubscribeableSet<T> {
  private internalSet: Set<T> | Map<T, number>;
  private store: Writable<Set<T>>;

  constructor(countable: boolean = false, initial: T[] = []) {
    this.internalSet = countable ? new Map(initial.map(x => [x, 1])) : new Set(initial);
    this.store = writable(new Set(initial));
  }

  has(value: T) {
    return this.internalSet.has(value);
  }

  subscribe(callback: (value: Set<T>) => void) {
    return this.store.subscribe(callback);
  }

  add(value: T) {
    if (this.internalSet instanceof Map) {
      const count = this.internalSet.get(value) || 0;
      this.internalSet.set(value, count + 1);
      if (count === 0) this.store.update(s => { s.add(value); return s; });
    } else {
      if (this.internalSet.has(value)) return;
      this.internalSet.add(value);
      this.store.set(this.internalSet);
    }
  }

  delete(value: T) {
    if (this.internalSet instanceof Map) {
      const count = this.internalSet.get(value) || 0;
      if (count === 1) {
        this.internalSet.delete(value);
        this.store.update(s => { s.delete(value); return s; });
      } else {
        this.internalSet.set(value, count - 1);
      }
    }
    else {
      if (!this.internalSet.has(value)) return;
      this.internalSet.delete(value);
      this.store.set(this.internalSet);
    }
  }

  set(value: Set<T>) {
    if (this.internalSet instanceof Map) {
      this.internalSet = new Map(Array.from(value).map(x => [x, 1]));
    } else {
      this.internalSet = new Set(value);
    }
    this.store.set(value);
  }
}

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
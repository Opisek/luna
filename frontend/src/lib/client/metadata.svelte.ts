import { browser } from "$app/environment";
import { SvelteMap, SvelteSet } from "svelte/reactivity";
import { getRepository } from "./repository.svelte";

class Metadata {
  //
  // Exported Subscribeables
  //

  // Faults
  faultySources = $state<SvelteSet<string>>(new SvelteSet());
  faultyCalendars = $state<SvelteSet<string>>(new SvelteSet());

  // Loading
  loadingSources = $state<SvelteMap<string, number>>(new SvelteMap());
  loadingCalendars = $state<SvelteMap<string, number>>(new SvelteMap());
  loadingData = $state<boolean>(false);
  private loadingCounter;

  // Hidden / Collapsed
  collapsedSources = $state<SvelteSet<string>>(new SvelteSet());
  hiddenCalendars = $state<SvelteSet<string>>(new SvelteSet());

  //
  // Constructor
  //
  constructor() {
    this.faultySources = new SvelteSet();
    this.faultyCalendars = new SvelteSet();

    this.loadingSources = new SvelteMap();
    this.loadingCalendars = new SvelteMap();
    this.loadingData = false;
    this.loadingCounter = 0;

    if (browser) {
      this.collapsedSources = new SvelteSet(JSON.parse(localStorage.getItem("collapsedSources") || "[]"));
      this.hiddenCalendars = new SvelteSet(JSON.parse(localStorage.getItem("hiddenCalendars") || "[]"));
    } else {
      this.collapsedSources = new SvelteSet();
      this.hiddenCalendars = new SvelteSet();
    }

    if (browser) {
      window.addEventListener("storage", () => {
        const newCollapsedSources = localStorage.getItem("collapsedSources");
        if (newCollapsedSources) this.collapsedSources = new SvelteSet<string>(JSON.parse(newCollapsedSources) as string[]);

        const newHiddenCalendars = localStorage.getItem("hiddenCalendars");
        if (newHiddenCalendars) this.hiddenCalendars = new SvelteSet<string>(JSON.parse(newHiddenCalendars) as string[]);
      })
    }
  }

  //
  // Logic
  //

  // Loading
  startLoading(): (() => void) {
    this.loadingCounter++;
    this.loadingData = true;

    let called = false;
    return (() => {
      if (!called) {
        called = true;
        if (--this.loadingCounter == 0) this.loadingData = false;
      }
    });
  }

  startLoadingSource(source: string): (() => void) {
    this.loadingSources.set(source, (this.loadingSources.get(source) || 0) + 1);
    const stopLoadingInternal = this.startLoading();

    let called = false;
    return (() => {
      if (!called) {
        called = true;
        const current = this.loadingSources.get(source) || 0;
        if (current <= 1) this.loadingSources.delete(source);
        else this.loadingSources.set(source, current - 1);
        stopLoadingInternal();
      }
    });
  }

  startLoadingCalendar(calendar: string): (() => void) {
    this.loadingCalendars.set(calendar, (this.loadingCalendars.get(calendar) || 0) + 1);
    const stopLoadingInternal = this.startLoading();
    // TODO: start loading the calendar's source as well

    let called = false;
    return (() => {
      if (!called) {
        called = true;
        const current = this.loadingCalendars.get(calendar) || 0;
        if (current <= 1) this.loadingCalendars.delete(calendar);
        else this.loadingCalendars.set(calendar, current - 1);
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
    if (browser) localStorage.setItem("collapsedSources", JSON.stringify(Array.from(this.collapsedSources)));
  }

  setCalendarVisibility = (calendarId: string, visible: boolean) => {
    if (visible) this.hiddenCalendars.delete(calendarId);
    else this.hiddenCalendars.add(calendarId);
    if (browser) localStorage.setItem("collapsedSources", JSON.stringify(Array.from(this.hiddenCalendars)));
    getRepository().recalculateEvents();
  }
}

let metadata: Metadata | null = null;
export function getMetadata() {
  if (metadata === null) {
    metadata = new Metadata();
  }
  return metadata;
}
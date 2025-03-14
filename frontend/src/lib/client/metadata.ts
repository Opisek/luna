import { writable, type Writable } from "svelte/store";

class Metadata {
  //
  // Exported Subscribeables
  //

  // Faults
  readonly faultySources = writable(new Set<string>());
  readonly faultyCalendars = writable(new Set<string>());

  // Loading
  readonly loadingSources = writable(new Set<string>());
  private readonly loadingSourcesCounter = new Map<string, number>();
  readonly loadingCalendars = writable(new Set<string>());
  readonly loadingData = writable(false);

  // Hidden / Collapsed

  //
  // Logic
  //

  // Misc
  private addToSet<T>(set: Writable<Set<T>>, source: T) {
    set.update(s => {
      s.add(source);
      return s;
    });
  }

  private removeFromSet<T>(set: Writable<Set<T>>, source: T) {
    set.update(s => {
      s.delete(source);
      return s;
    });
  }

  private addToCounter<T>(map: Map<T, number>, set: Writable<Set<T>>, source: T) {
    map.set(source, (map.get(source) || 0) + 1);
    this.addToSet(set, source);
  }

  private removeFromCounter<T>(map: Map<T, number>, set: Writable<Set<T>>, source: T) {
    const count = map.get(source);
    if (count === undefined) return;
    if (count === 1) {
      map.delete(source);
      this.removeFromSet(set, source);
    } else {
      map.set(source, count - 1);
    }
  }

  // Loading
  private loadingCounter = 0;

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
    this.addToCounter(this.loadingSourcesCounter, this.loadingSources, source);
    const stopLoadingInternal = this.startLoading();

    let called = false;
    return (() => {
      if (!called) {
        called = true;
        this.removeFromCounter(this.loadingSourcesCounter, this.loadingSources, source);
        stopLoadingInternal();
      }
    });
  }

  startLoadingCalendar(calendar: string): (() => void) {
    this.addToSet(this.loadingCalendars, calendar);
    const stopLoadingInternal = this.startLoading();
    // TODO: start loading the calendar's source as well

    let called = false;
    return (() => {
      if (!called) {
        called = true;
        this.removeFromSet(this.loadingCalendars, calendar);
        stopLoadingInternal();
      }
    });
  }

  // Faults
  addFaultySource(source: string, fault: string) {
    this.addToSet(this.faultySources, source);
  }

  removeFaultySource(source: string) {
    this.removeFromSet(this.faultySources, source);
  }

  addFaultyCalendar(calendar: string, fault: string) {
    this.addToSet(this.faultyCalendars, calendar);
  }

  removeFaultyCalendar(calendar: string) {
    this.removeFromSet(this.faultyCalendars, calendar);
  }
}

let metadata: Metadata | null = null;
export function getMetadata() {
  if (metadata === null) {
    metadata = new Metadata();
  }
  return metadata;
}
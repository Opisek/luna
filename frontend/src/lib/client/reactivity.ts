import { writable, type Writable } from "svelte/store";

export class SubscribeableSet<T> {
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

export class SubscribeableArray<T> {
  private internalArray: T[];
  private store: Writable<T[]>;

  constructor(initial: T[] = []) {
    this.internalArray = initial;
    this.store = writable(initial);
  }

  has(value: T) {
    return this.internalArray.includes(value);
  }

  subscribe(callback: (value: T[]) => void) {
    return this.store.subscribe(callback);
  }

  add(value: T) {
    if (this.internalArray.includes(value)) return;
    this.internalArray.push(value);
    this.store.set(this.internalArray);
  }

  delete(value: T) {
    const index = this.internalArray.indexOf(value);
    if (index === -1) return;
    this.internalArray.splice(index, 1);
    this.store.set(this.internalArray);
  }

  set(value: T[]) {
    this.internalArray = value;
    this.store.set(value);
  }

  get(index: number) {
    return this.internalArray[index];
  }

  update(callback: (value: T[]) => T[]) {
    this.internalArray = callback(this.internalArray);
    this.store.set(this.internalArray);
  }

  getStore() {
    return this.store;
  }
}
export function compareEventsByStartDate(a: EventModel, b: EventModel): number {
  if (a.date.start === b.date.start) {
    return a.date.allDay ? b.date.allDay ? 0 : 1 : b.date.allDay ? -1 : 0;
  }
  return a.date.start.getTime() - b.date.start.getTime();
}
export function compareEventsByStartDate(a: EventModel, b: EventModel): number {
  return a.date.start.getTime() - b.date.start.getTime();
}
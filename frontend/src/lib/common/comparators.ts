export function compareEventsByStartDate(a: CalendarEventModel, b: CalendarEventModel): number {
  return a.start.getTime() - b.start.getTime();
}
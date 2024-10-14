export function compareEventsByStartDate(a: EventModel, b: EventModel): number {
  const diff = a.date.start.getTime() - b.date.start.getTime();
  if (diff === 0) {
    return a.date.allDay ? b.date.allDay ? 0 : 1 : b.date.allDay ? -1 : 0;
  }
  return diff;
}
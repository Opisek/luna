export function compareEventsByStartDate(a: EventModel, b: EventModel): number {
  const aDay = new Date(a.date.start);
  aDay.setHours(0, 0, 0, 0);
  const bDay = new Date(b.date.start);
  bDay.setHours(0, 0, 0, 0);

  const diff = aDay.getTime() - bDay.getTime();
  if (diff !== 0) return diff;

  const aMultiday = a.date.start.getDate() !== a.date.end.getDate() || a.date.start.getMonth() !== a.date.end.getMonth() || a.date.start.getFullYear() !== a.date.end.getFullYear();
  const bMultiday = b.date.start.getDate() !== b.date.end.getDate() || b.date.start.getMonth() !== b.date.end.getMonth() || b.date.start.getFullYear() !== b.date.end.getFullYear();

  if (aMultiday && !bMultiday) return -1;
  if (!a.date.allDay && bMultiday) return 1;

  const exactDiff = a.date.start.getTime() - b.date.start.getTime();
  if (exactDiff !== 0) return exactDiff;

  const aDuration = a.date.end.getTime() - a.date.start.getTime();
  const bDuration = b.date.end.getTime() - b.date.start.getTime();

  const durationDiff = bDuration - aDuration;
  if (durationDiff !== 0) return durationDiff;

  return a.name < b.name ? -1 : a.name > b.name ? 1 : 0;
}
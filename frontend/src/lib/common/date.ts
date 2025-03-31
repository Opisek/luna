import { getSettings } from "$lib/client/settings.svelte";
import { UserSettingKeys } from "../../types/settings";

export function isSameDay(a: Date, b: Date): boolean {
  return a.getFullYear() === b.getFullYear() && a.getMonth() === b.getMonth() && a.getDate() === b.getDate();
}

export function isInRange(date: Date, start: Date, end: Date): boolean {
  return date.getTime() >= start.getTime() && date.getTime() <= end.getTime();
}

export function getDayIndex(date: Date) {
  const settings = getSettings();
  const startOfWeek = settings.userSettings[UserSettingKeys.FirstDayOfWeek];
  const day = date.getDay();
  const dayIndex = (day - startOfWeek + 7) % 7;
  return dayIndex;
}
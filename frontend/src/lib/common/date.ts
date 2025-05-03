import { getSettings } from "$lib/client/data/settings.svelte";
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

// https://www.iso.org/obp/ui/#iso:std:iso:8601:-1:ed-1:v1:en
export function getWeekNumber(date: Date) {
  // Overflow check for last 3 days of december
  if (date.getMonth() === 11 && date.getDate() >= 29 && (date.getDay() + 6) % 7 < 4) {
    const januaryFirstNextYear = new Date(Date.UTC(date.getFullYear() + 1, 0, 1));
    const januaryFirstNextYearDay = (januaryFirstNextYear.getDay() + 6) % 7;
    if (januaryFirstNextYearDay < 4) return 1;
  }

  // Determine the first monday of the first week of the year
  const januaryFirst = new Date(Date.UTC(date.getFullYear(), 0, 1));
  const januaryFirstDay = (januaryFirst.getDay() + 6) % 7;
  const firstMonday = new Date(januaryFirst);
  if (januaryFirstDay < 4) firstMonday.setDate(firstMonday.getDate() - januaryFirstDay);
  else firstMonday.setDate(firstMonday.getDate() - januaryFirstDay + 7);

  // Determine the amount of days between the first monday and the date
  const dateNoTime = new Date(Date.UTC(date.getFullYear(), date.getMonth(), date.getDate()));
  const diff = dateNoTime.getTime() - firstMonday.getTime();
  const diffDays = Math.floor(diff / (1000 * 60 * 60 * 24));
  const weekNumber = Math.floor(diffDays / 7) + 1;

  return weekNumber;
}

export function getWeekMonth(weekNumber: number, year: number) {
  // Get first monday of first week of the year
  const januaryFirst = new Date(Date.UTC(year, 0, 1));
  const januaryFirstDay = (januaryFirst.getDay() + 6) % 7;
  const firstMonday = new Date(januaryFirst);
  if (januaryFirstDay < 4) firstMonday.setDate(firstMonday.getDate() - januaryFirstDay);
  else firstMonday.setDate(firstMonday.getDate() - januaryFirstDay + 7);

  // Get the thruesday of the requested week
  const thursday = new Date(firstMonday);
  thursday.setDate(thursday.getDate() + (weekNumber - 1) * 7 + 3);

  // If the majority of the days is in some month, then they include thursday
  return thursday.getMonth();
}

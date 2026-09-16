import { t } from "@sveltia/i18n";

const dayNames = [
  "sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday"
]
const monthNames = [
  "january", "february", "march", "april", "may", "june", "july", "august", "september", "october", "november", "december"
]

export function getDayName(day: number, short?: boolean): string {
  return t(`days.${short ? "short" : "full"}`, { values: { day: day } })
}
export function getWeekdayName(day: number, short?: boolean): string {
  return t(`weekdays.${short ? "short" : "full"}.${dayNames[day]}`)
}
export function getMonthName(month: number, short?: boolean): string {
  return t(`months.${short ? "short" : "full"}.${monthNames[month]}`);
}
export function getYearName(year: number, short?: boolean): string {
  return t(`years.${short ? "short" : "full"}`, { values: { year: year } });
}
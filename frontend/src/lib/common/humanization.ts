import { _ as t } from "@sveltia/i18n";

const dayNames = [
  "sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday"
]
const monthNames = [
  "january", "february", "march", "april", "may", "june", "july", "august", "september", "october", "november", "december"
]

export function getDayName(day: number, short?: boolean): string {
  return t(`calendar.weekdays.${short ? "short" : "full"}.${dayNames[day]}`);
}
export function getMonthName(month: number, short?: boolean): string {
  return t(`calendar.months.${short ? "short" : "full"}.${monthNames[month]}`);
}
const dayNames = [
  "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"
]
const monthNames = [
  "January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"
]

export function getDayName(day: number): string {
  return dayNames[day];
}
export function getMonthName(month: number): string {
  return monthNames[month];
}
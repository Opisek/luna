import dayjs from "dayjs";
import utc from "dayjs/plugin/utc";
import timezone from "dayjs/plugin/timezone";

dayjs.extend(utc);
dayjs.extend(timezone);

export function parseTimestampList(str: string): Date[] {
  const parts = str.split(":");
  const params = parts[0].split(";");
  const timestamps = parts[1].split(",");

  if (params.length * timestamps.length == 0) return [];
  // @ts-ignore
  if (!["RDATE", "EXDATE"].includes(params.shift())) return [];

  let dateType = "DATE-TIME";
  let timezoneId = "UTC";

  while (true) {
    const param = params.shift();
    if (!param) break;

    const paramParts = param?.split("=");
    if (paramParts.length != 2) continue;

    const key = paramParts[0];
    const value = paramParts[1];

    switch (key) {
      case "VALUE":
        dateType = value;
        break
      case "TZID":
        timezoneId = value;
        break
    }
  }

  const format = dateType == "DATE" ? "YYYYMMDD" : "YYYYMMDD_HHmmss";

  return timestamps.map(timestamp => dayjs(timestamp, format, timezoneId).toDate());
}

export function serializeTimestampList(listName: string, allDay: boolean, timezoneId: string, dates: Date[]): string {
  const format = allDay ? "YYYYMMDD" : "YYYYMMDD[T]HHmmss";
  return `${listName};VALUE=${allDay ? "DATE" : "DATE-TIME"};TZID=${timezoneId}:${dates.map(x => dayjs(x).tz(timezoneId).format(format))}`
}
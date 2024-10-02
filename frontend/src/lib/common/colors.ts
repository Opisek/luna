export const parseRGB = (color: string): [number, number, number] => {
  return [
    parseInt(color.substring(1, 3), 16),
    parseInt(color.substring(3, 5), 16),
    parseInt(color.substring(5, 7), 16),
  ];
}

export const serializeRGB = (rgb: [number, number, number]): string => {
  return `#${rgb[0].toString(16).padStart(2, '0')}${rgb[1].toString(16).padStart(2, '0')}${rgb[2].toString(16).padStart(2, '0')}`;
}

export const defaultEventRGB: [number, number, number] = [90, 150, 225];
export const defaultEventColor: string = serializeRGB(defaultEventRGB);
export const defaultCalendarRGB: [number, number, number] = [90, 150, 225];
export const defaultCalendarColor: string = serializeRGB(defaultCalendarRGB);
export const recommendedRGB: [number, number, number][] = [
  [216, 110, 100], // red
  [226, 165, 90], // orange
  [226, 203, 90], // yellow
  [93, 60, 62], // green
  [90, 150, 225], // blue
  [165, 90, 226], // purple
  [226, 90, 151], // pink
];
export const recommendedColors: string[] = recommendedRGB.map(serializeRGB);

export function isDark(rgb: number[]) {
  const brightness = Math.round(((rgb[0] * 299) +
                      (rgb[1] * 587) +
                      (rgb[2] * 114)) / 1000);
  return brightness <= 141;
}

// credit: https://stackoverflow.com/questions/8022885/rgb-to-hsv-color-in-javascript
export const RGBtoHSV = (rgb: number[]): [number, number, number] => {
  const r = rgb[0] / 255
  const g = rgb[1] / 255
  const b = rgb[2] / 255

  let v=Math.max(r,g,b), c=v-Math.min(r,g,b);
  let h= c && ((v==r) ? (g-b)/c : ((v==g) ? 2+(b-r)/c : 4+(r-g)/c)); 
  return [60*(h<0?h+6:h), v&&c/v, v];
}

// TODO: would prefer event.getColor() but i could not figure out how to do this without creating an additional interface or class
export const GetEventRGB = (event: EventModel | null) => {
  if (event && event.color) {
    return parseRGB(event.color);
  } else {
    return defaultEventRGB;
  }
}
export const GetEventColor = (event: EventModel | null) => {
  if (event && event.color) {
    return event.color;
  } else {
    return defaultEventColor;
  }
}
export const GetCalendarRGB = (calendar: CalendarModel | null) => {
  if (calendar && calendar.color) {
    return parseRGB(calendar.color);
  } else {
    return defaultCalendarRGB;
  }
}
export const GetCalendarColor = (calendar: CalendarModel | null) => {
  if (calendar && calendar.color) {
    return calendar.color;
  } else {
    return defaultCalendarColor;
  }
}
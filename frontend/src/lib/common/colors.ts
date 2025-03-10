export const parseRGB = (color: string): [number, number, number] => {
  return [
    parseInt(color.substring(1, 3), 16),
    parseInt(color.substring(3, 5), 16),
    parseInt(color.substring(5, 7), 16),
  ];
}

export const serializeRGB = (rgb: [number, number, number]): string => {
  return `#${rgb[0].toString(16).padStart(2, '0')}${rgb[1].toString(16).padStart(2, '0')}${rgb[2].toString(16).padStart(2, '0')}`.toUpperCase();
}

export const defaultEventRGB: [number, number, number] = [90, 150, 225];
export const defaultEventColor: string = serializeRGB(defaultEventRGB);
export const defaultCalendarRGB: [number, number, number] = [90, 150, 225];
export const defaultCalendarColor: string = serializeRGB(defaultCalendarRGB);
export const recommendedRGB: [number, number, number][] = [
  [225, 90, 90], // red
  [225, 165, 90], // orange
  [226, 203, 90], // yellow
  [151, 225, 90], // green
  [90, 150, 225], // blue
  [165, 90, 225], // purple
  [225, 90, 151], // pink
];
export const recommendedColors: string[] = recommendedRGB.map(serializeRGB);

export function isDark(rgb: number[]) {
  const brightness = Math.round(((rgb[0] * 299) +
                      (rgb[1] * 587) +
                      (rgb[2] * 114)) / 1000);
  return brightness <= 141;
}

// credit: https://gist.github.com/mjackson/5311256
export const RGBtoHSL = (rgb: number[]): [number, number, number] => {
  const r = rgb[0] / 255
  const g = rgb[1] / 255
  const b = rgb[2] / 255

  var max = Math.max(r, g, b), min = Math.min(r, g, b);
  var h, s, l = (max + min) / 2;

  if (max == min) {
    h = s = 0; // achromatic
  } else {
    var d = max - min;
    s = l > 0.5 ? d / (2 - max - min) : d / (max + min);

    switch (max) {
      case r: h = (g - b) / d + (g < b ? 6 : 0); break;
      case g: h = (b - r) / d + 2; break;
      case b: h = (r - g) / d + 4; break;
    }
  }

  if (h === undefined) {
    h = 0;
  }

  return [Math.round(h * 60), Math.round(s * 100), Math.round(l * 100)];
}
export const HSLtoRGB = (hsl: [number, number, number]): [number, number, number] => {
  const h = hsl[0] / 360;
  const s = hsl[1] / 100;
  const l = hsl[2] / 100;

  let r, g, b;

  if (s == 0) {
    r = g = b = l; // achromatic
  } else {
    function hue2rgb(p: number, q: number, t: number): number {
      if (t < 0) t += 1;
      if (t > 1) t -= 1;
      if (t < 1/6) return p + (q - p) * 6 * t;
      if (t < 1/2) return q;
      if (t < 2/3) return p + (q - p) * (2/3 - t) * 6;
      return p;
    }

    let q = l < 0.5 ? l * (1 + s) : l + s - l * s;
    let p = 2 * l - q;

    r = hue2rgb(p, q, h + 1/3);
    g = hue2rgb(p, q, h);
    b = hue2rgb(p, q, h - 1/3);
  }

  return [ Math.round(r * 255), Math.round(g * 255), Math.round(b * 255) ];
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

export const isValidColor = (color: string | null | undefined) => {
  return color !== null && color !== undefined && color.length === 7 && /^#[0-9a-fA-F]+$/.test(color);
}

export function calculateSecondaryColor(rgb: [number, number, number]): [number, number, number] {
  const hsl = RGBtoHSL(rgb);
  isDark(rgb) ? hsl[2] += 10 : hsl[2] -= 5;
  return HSLtoRGB(hsl);
}

export const GetEventHoverColor = (event: EventModel | null) => {
  if (event && event.color) {
    return serializeRGB(calculateSecondaryColor(parseRGB(event.color)));
  } else {
    return serializeRGB(calculateSecondaryColor(defaultEventRGB));
  }
}
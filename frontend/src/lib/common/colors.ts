export function isDark(rgb: number[]) {
  const brightness = Math.round(((rgb[0] * 299) +
                      (rgb[1] * 587) +
                      (rgb[2] * 114)) / 1000);
  return brightness <= 141;
}
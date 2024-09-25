export const getRedirectPage = (url: URL): string => {
  let redirectPage = url.searchParams.get("redirect");
  if (redirectPage == null || redirectPage == "") {
    redirectPage = '/';
  } else {
    redirectPage = decodeURIComponent(redirectPage);
  }
  return redirectPage;
}

export const parseRGB = (color: string): [number, number, number] => {
  return [
    parseInt(color.substring(1, 3), 16),
    parseInt(color.substring(3, 5), 16),
    parseInt(color.substring(5, 7), 16),
  ];
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
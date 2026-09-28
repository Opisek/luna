export function encodeRedirectUrl(url: URL): string {
  return encodeURIComponent(`${url.pathname}?${url.searchParams.toString()}`);
}
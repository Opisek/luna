export const getRedirectPage = (url: URL): string => {
  let redirectPage = url.searchParams.get("redirect");
  if (redirectPage == null || redirectPage == "") {
    redirectPage = '/';
  } else {
    redirectPage = decodeURIComponent(redirectPage);
  }
  return redirectPage;
}

const databaseFileIdRegex = /\/api\/files\/([a-f0-9-]{36})/;
export const getDatabaseFileIdFromUrl = (url: string): string | null => {
  if (!url) return null;
  const match = url.match(databaseFileIdRegex);
  if (match && match[1]) {
    return match[1];
  }
  return null;
}
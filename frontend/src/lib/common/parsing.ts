export const getRedirectPage = (url: URL): string => {
  let redirectPage = url.searchParams.get("redirect");
  if (redirectPage == null || redirectPage == "") {
    redirectPage = '/';
  } else {
    redirectPage = decodeURIComponent(redirectPage);
  }
  return redirectPage;
}

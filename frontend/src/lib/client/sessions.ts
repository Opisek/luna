import { browser } from "$app/environment";

export function clearSession() {
  if (!browser) return;
  localStorage.clear();
  window.location.href = "/logout";
}
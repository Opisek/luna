import { browser } from "$app/environment";
import { register, init, getLocaleFromNavigator, locales, locale, waitLocale } from "@sveltia/i18n";
import { parse } from "yaml";

register("en-US", () => import("../../lang/en-US.yaml?raw").then(m => parse(m.default)));
register("de-DE", () => import("../../lang/de-DE.yaml?raw").then(m => parse(m.default)));

init({ fallbackLocale: "en-US" });

export async function loadLanguage(userChoice: string | null | undefined) {
  await locale.set(await getCurrentLanguage(userChoice));
  await waitLocale("en-US");
  await waitLocale();
}

export async function getCurrentLanguage(userChoice: string | null | undefined) {
  if (!userChoice || !locales.includes(userChoice)) return getDefaultLanguage();
  return userChoice;
}

export async function getDefaultLanguage() {
  return (browser ? getLocaleFromNavigator() : null) ?? "en-US";
}
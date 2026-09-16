import { browser } from "$app/environment";
import { register, init, getLocaleFromNavigator, locales, locale, waitLocale, registerMessageFunction, t } from "@sveltia/i18n";
import { getLocaleDir } from "messageformat/functions";
import { parse } from "yaml";
import type { Option } from "../../types/options";
import { getSettings } from "$lib/client/data/settings.svelte";
import { UserSettingKeys } from "../../types/settings";

const languages = [ "en-US", "de-DE", "pl-PL", "ja-JP" ];

languages.forEach(x => register(x, () => import(`../../lang/${x}.yaml?raw`).then(m => parse(m.default))));

init({ fallbackLocale: "en-US" });

registerMessageFunction('date', (ctx, options, operand) => {
  let locale: Intl.LocalesArgument = ctx.locales;

  if ("format" in options && options.format === "custom") {
    const settings = getSettings().userSettings;

    const customLocale = settings[UserSettingKeys.DateLocale];
    if (customLocale !== "default") locale = customLocale;

    const customHourCycle = settings[UserSettingKeys.HourCycle];
    if (customHourCycle !== "default") options.hourCycle = customHourCycle;
  }
  delete options.format;

  // @ts-ignore
  console.log(locale);
  const dtf = new Intl.DateTimeFormat(locale, options);

  return {
    type: 'string',
    dir: getLocaleDir(dtf.resolvedOptions().locale),
    // @ts-ignore
    toString: () => dtf.format(operand),
  };
});

registerMessageFunction('number', (ctx, options, operand) => {
    // @ts-ignore
  const dtf = new Intl.NumberFormat(ctx.locales, {
    numberingSystem: options.numberingSystem ?? 'native',
    style: options.style ?? "decimal",
    useGrouping: options.useGrouping ?? false,
  });

  return {
    type: 'string',
    dir: getLocaleDir(dtf.resolvedOptions().locale),
    // @ts-ignore
    toString: () => dtf.format(operand),
  };
});

export async function loadLanguage(userChoice: string | null | undefined) {
  await locale.set(await getCurrentLanguage(userChoice));
  await waitLocale();
}

export async function getCurrentLanguage(userChoice: string | null | undefined) {
  if (!userChoice || !locales.includes(userChoice)) return getDefaultLanguage();
  return userChoice;
}

export async function getDefaultLanguage() {
  return (browser ? getLocaleFromNavigator() : null) ?? "en-US";
}

// Many locales share the same date format
// The following set has been determined experimentally
const dateLocaleOptions = [ "ar-SA", "bg-BG", "bn-IN", "de-DE", "en-AU", "en-CA", "en-GB", "en-US", "fa-IR", "fi-FI", "he-IL", "hr-HR", "hu-HU", "ja-JP", "ko-KR", "mk-MK", "mn-MN", "mr-IN", "nl-NL", "pl-PL", "pt-PT", "sk-SK", "sq-AL", "sr-RS", "te-IN", "zh-CN", "zh-HK" ];

// The following is a date in which 

export function getDateLocaleCalculationExampleDate() {
  const exampleDate = new Date("2000-01-31T00:00:00");
  exampleDate.setFullYear(new Date().getFullYear());
  return exampleDate;
}
export async function getUniqueDateLocaleOptions(): Promise<Option<string>[]> {
  const exampleDate = getDateLocaleCalculationExampleDate();

  // Map different date formats to ONE example locale that provides it
  const uniqueDateFormats = Object.fromEntries(dateLocaleOptions.map(locale => [
    new Intl.DateTimeFormat(locale, { dateStyle: "short" }).format(exampleDate),
    locale
  ]));

  // Format these into options that work in our codebase
  const uniqueOptions = Object.entries(uniqueDateFormats).map(x => ({
    value: x[1],
    name: x[0]
  }));

  // Sort and return
  return uniqueOptions.sort((a, b) => (a.name < b.name ? -1 : a.name > b.name ? 1 : 0));
}
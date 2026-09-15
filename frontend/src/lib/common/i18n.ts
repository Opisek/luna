import { browser } from "$app/environment";
import { register, init, getLocaleFromNavigator, locales, locale, waitLocale, registerMessageFunction, t } from "@sveltia/i18n";
import { getLocaleDir } from "messageformat/functions";
import { parse } from "yaml";

const languages = [ "en-US", "de-DE", "pl-PL", "fr-FR" ];

languages.forEach(x => register(x, () => import(`../../lang/${x}.yaml?raw`).then(m => parse(m.default))));
register("en-DE", () => import(`../../lang/en-US.yaml?raw`).then(m => parse(m.default))); // This is a cheat to get a DD/MM/YYYY format with English. It will be removed once a better method is developed.

init({ fallbackLocale: "en-DE" });

registerMessageFunction('weekday', (ctx, options, operand) => {
  // @ts-ignore
  const dtf = new Intl.DateTimeFormat(ctx.locales, { weekday: options.weekday ?? 'short' });

  return {
    type: 'string',
    dir: getLocaleDir(dtf.resolvedOptions().locale),
    // @ts-ignore
    toString: () => dtf.format(operand),
  };
});

registerMessageFunction('month', (ctx, options, operand) => {
    // @ts-ignore
  const dtf = new Intl.DateTimeFormat(ctx.locales, { month: options.month ?? 'short' });

  return {
    type: 'string',
    dir: getLocaleDir(dtf.resolvedOptions().locale),
    // @ts-ignore
    toString: () => dtf.format(operand),
  };
});

registerMessageFunction('ordinal', (ctx, options, operand) => {
    // @ts-ignore
  const dtf = new Intl.DateTimeFormat(ctx.locales);
  const locale = getLocaleDir(dtf.resolvedOptions().locale)

  return {
    type: 'string',
    dir: getLocaleDir(dtf.resolvedOptions().locale),
    // @ts-ignore
    toString: () => t(`numbers.ordinal.${options.order ?? "normal"}`, { values: { num: operand }, locale: locale }),
  };
});

export async function loadLanguage(userChoice: string | null | undefined) {
  await locale.set(await getCurrentLanguage(userChoice));
  await waitLocale("en-DE");
  await waitLocale();
}

export async function getCurrentLanguage(userChoice: string | null | undefined) {
  if (!userChoice || !locales.includes(userChoice)) return getDefaultLanguage();
  return userChoice;
}

export async function getDefaultLanguage() {
  return (browser ? getLocaleFromNavigator() : null) ?? "en-DE";
}
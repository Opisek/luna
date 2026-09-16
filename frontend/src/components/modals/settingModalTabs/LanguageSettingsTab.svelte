<script lang="ts">
  import type { Settings } from "../../../lib/client/data/settings.svelte";
  import type { Option } from "../../../types/options";
  import { UserSettingKeys } from "../../../types/settings";
  import SelectInput from "../../forms/SelectInput.svelte";
  import { locale, t } from "@sveltia/i18n";
  import { loadLanguage } from "$lib/common/i18n";

  interface Props {
    settings: Settings;
    languages: Option<string>[];
    dateLocales: Option<string>[];
  }

  let {
    settings,
    languages,
    dateLocales
  }: Props = $props();

  const exampleTwentythree = new Date("1990-01-01T23:00");
  const exampleMidnight = new Date("1990-01-01T24:00");
  let defaultHourCycle = $derived.by(() => {
    const format = Intl.DateTimeFormat(locale.current, {
      numberingSystem: "arabic",
      timeStyle: "short",
    });
    return `h${
      Number.parseInt(format.format(exampleTwentythree).substring(0,2))
      +
      (format.format(exampleMidnight).startsWith("0") ? 0 : 1)
    }`;
  });
  let effectiveDateLocale = $derived.by(() => {
    let loc = locale.current;
    const userLoc = settings.userSettings[UserSettingKeys.DateLocale];
    if (userLoc !== "default") loc = userLoc;
    return loc;
  });
  let defaultHourCycleOption: Option<string> = $derived({
    name: t("settings.language.time.default", {
      values: {
        default: t(`settings.language.time.options.${defaultHourCycle}`, {
          values: {
            midnight: Intl.DateTimeFormat(effectiveDateLocale, { timeStyle: "short" }).format(exampleMidnight),
          },
        })
      }
    }),
    value: "default",
  });
  let hourCycleOptions = $derived(["h11", "h12", "h23", "h24"].map(x => ({
    value: x,
    name: t(`settings.language.time.options.${x}`, {
      values: {
        midnight: Intl.DateTimeFormat(effectiveDateLocale, { hourCycle: x as ("h11" | "h12" | "h23" | "h24"), timeStyle: "short" }).format(exampleMidnight),
      },
    }),
  })));
  let combinedHourCycleOptions = $derived([defaultHourCycleOption].concat(hourCycleOptions));
</script>

<SelectInput
  name={UserSettingKeys.Language}
  placeholder={t("settings.language.display.label")}
  bind:value={settings.userSettings[UserSettingKeys.Language]}
  options={languages}
  click={(l) => { loadLanguage(l) }}
/>

<SelectInput
  name={UserSettingKeys.DateLocale}
  placeholder={t("settings.language.date.label")}
  bind:value={settings.userSettings[UserSettingKeys.DateLocale]}
  options={dateLocales}
/>

<SelectInput
  name={UserSettingKeys.HourCycle}
  placeholder={t("settings.language.time.label")}
  bind:value={settings.userSettings[UserSettingKeys.HourCycle]}
  options={combinedHourCycleOptions}
/>
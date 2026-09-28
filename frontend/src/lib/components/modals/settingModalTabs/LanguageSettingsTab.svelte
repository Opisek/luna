<script lang="ts">
  import type { Settings } from "$lib/scripts/client/data/settings.svelte";
  import type { Option } from "$lib/types/options";
  import { UserSettingKeys } from "$lib/types/settings";
  import SelectInput from "../../forms/SelectInput.svelte";
  import { locale, t } from "@sveltia/i18n";
  import { loadLanguage } from "$lib/scripts/common/i18n";

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

  const exampleMidnight = new Date("1990-01-01T24:00");

  let effectiveDateLocale = $derived.by(() => {
    let loc = locale.current;
    const userLoc = settings.userSettings[UserSettingKeys.DateLocale];
    if (userLoc !== "default") loc = userLoc;
    return loc;
  });

  let defaultHourCycle = $derived(new Intl.DateTimeFormat(effectiveDateLocale, { hour: "numeric" }).resolvedOptions().hourCycle || "h23");
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
  
  let defaultDigitCounts = $derived.by(() => {
    const opts = new Intl.DateTimeFormat(effectiveDateLocale, { hour: "numeric" }).resolvedOptions();
    return {
      hour: (opts.hour === "2-digit" ? 2 : 1),
      minute: (opts.minute === "2-digit" ? 2 : 1),
      second: (opts.second === "2-digit" ? 2 : 1),
    };
  })
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

<!--<ToggleInput
  name={UserSettingKeys.ForceTwoTimeDigits}
  description={t("settings.language.digits.label")}
  bind:value={settings.userSettings[UserSettingKeys.ForceTwoTimeDigits]}
/>-->
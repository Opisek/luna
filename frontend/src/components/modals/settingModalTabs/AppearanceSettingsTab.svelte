<script lang="ts">
  import type { Settings } from "../../../lib/client/data/settings.svelte";
  import type { Option } from "../../../types/options";
  import { UserSettingKeys } from "../../../types/settings";
  import SelectInput from "../../forms/SelectInput.svelte";
  import SliderInput from "../../forms/SliderInput.svelte";
  import ToggleInput from "../../forms/ToggleInput.svelte";
  import SectionDivider from "../../layout/SectionDivider.svelte";
  import { number, t } from "@sveltia/i18n";
  import { loadLanguage } from "$lib/common/i18n";

  interface Props {
    settings: Settings;
    lightThemes: Option<string>[];
    darkThemes: Option<string>[];
    fonts: Option<string>[];
    languages: Option<string>[];
  }

  let {
    settings,
    lightThemes,
    darkThemes,
    fonts,
    languages,
  }: Props = $props();
</script>

<SectionDivider title={t("settings.appearance.calendar.subtitle")}/>
<ToggleInput
  name={UserSettingKeys.DisplayAllDayEventsFilled}
  description={t("settings.appearance.calendar.fill.allDay")}
  bind:value={settings.userSettings[UserSettingKeys.DisplayAllDayEventsFilled]}
/>
<ToggleInput
  name={UserSettingKeys.DisplayNonAllDayEventsFilled}
  description={t("settings.appearance.calendar.fill.nonAllDay")}
  bind:value={settings.userSettings[UserSettingKeys.DisplayNonAllDayEventsFilled]}
/>
<ToggleInput
  name={UserSettingKeys.DisplaySmallCalendar}
  description={t("settings.appearance.calendar.small")}
  bind:value={settings.userSettings[UserSettingKeys.DisplaySmallCalendar]}
/>
<ToggleInput
  name={UserSettingKeys.DynamicCalendarRows}
  description={t("settings.appearance.calendar.dynamic.main")}
  bind:value={settings.userSettings[UserSettingKeys.DynamicCalendarRows]}
/>
{#if settings.userSettings[UserSettingKeys.DisplaySmallCalendar]}
  <ToggleInput
    name={UserSettingKeys.DynamicSmallCalendarRows}
    description={t("settings.appearance.calendar.dynamic.small")}
    bind:value={settings.userSettings[UserSettingKeys.DynamicSmallCalendarRows]}
  />
{/if}
<ToggleInput
  name={UserSettingKeys.DisplayWeekNumbers}
  description={t("settings.appearance.calendar.week.numbers")}
  bind:value={settings.userSettings[UserSettingKeys.DisplayWeekNumbers]}
/>
<SelectInput
  name={UserSettingKeys.FirstDayOfWeek}
  placeholder={t("settings.appearance.calendar.week.first")}
  bind:value={settings.userSettings[UserSettingKeys.FirstDayOfWeek]}
  options={[
    { name: t("weekdays.full.monday"), value: 1 },
    { name: t("weekdays.full.tuesday"), value: 2 },
    { name: t("weekdays.full.wednesday"), value: 3 },
    { name: t("weekdays.full.thursday"), value: 4 },
    { name: t("weekdays.full.friday"), value: 5 },
    { name: t("weekdays.full.saturday"), value: 6 },
    { name: t("weekdays.full.sunday"), value: 0 }
  ]}
/>
<SectionDivider title={t("settings.appearance.site.subtitle")}/>
<SelectInput
  name={UserSettingKeys.Language}
  placeholder={t("settings.appearance.site.language")}
  bind:value={settings.userSettings[UserSettingKeys.Language]}
  options={languages}
  click={(l) => { loadLanguage(l) }}
/>
<ToggleInput
  name={UserSettingKeys.AppearenceFrostedGlass}
  description={t("settings.appearance.site.frosted")}
  bind:value={settings.userSettings[UserSettingKeys.AppearenceFrostedGlass]}
/>
<ToggleInput
  name={UserSettingKeys.DisplayRoundedCorners}
  description={t("settings.appearance.site.rounded")}
  bind:value={settings.userSettings[UserSettingKeys.DisplayRoundedCorners]}
/>
<ToggleInput
  name={UserSettingKeys.UseTextButtons}
  description={t("settings.appearance.site.buttons")}
  bind:value={settings.userSettings[UserSettingKeys.UseTextButtons]}
/>
<SliderInput
  name={UserSettingKeys.UiScaling}
  title={t("settings.appearance.site.scaling")}
  bind:value={settings.userSettings[UserSettingKeys.UiScaling]}
  min={0.5}
  max={1.5}
  step={0.05}
  detentTransform={(value) => number(value, { format: "percent" })}
/>
<SelectInput
  name={UserSettingKeys.ThemeLight}
  placeholder={t("settings.appearance.site.theme.light")}
  bind:value={settings.userSettings[UserSettingKeys.ThemeLight]}
  options={lightThemes}
/>
<SelectInput
  name={UserSettingKeys.ThemeDark}
  placeholder={t("settings.appearance.site.theme.dark")}
  bind:value={settings.userSettings[UserSettingKeys.ThemeDark]}
  options={darkThemes}
/>
<ToggleInput
  name={UserSettingKeys.ThemeSynchronize}
  description={t("settings.appearance.site.theme.sync")}
  bind:value={settings.userSettings[UserSettingKeys.ThemeSynchronize]}
/>
<SelectInput
  name={UserSettingKeys.FontText}
  placeholder={t("settings.appearance.site.font.text")}
  bind:value={settings.userSettings[UserSettingKeys.FontText]}
  options={fonts}
/>
<SelectInput
  name={UserSettingKeys.FontTime}
  placeholder={t("settings.appearance.site.font.mono")}
  bind:value={settings.userSettings[UserSettingKeys.FontTime]}
  options={fonts}
/>
<SectionDivider title={t("settings.appearance.animations.subtitle")}/>
<SliderInput
  name={UserSettingKeys.AnimationDuration}
  title={t("settings.appearance.animations.duration.display")}
  info={t("settings.appearance.animations.duration.info")}
  bind:value={settings.userSettings[UserSettingKeys.AnimationDuration]}
  min={0}
  max={2}
  step={0.1}
  detentTransform={(value) => number(value, { format: "percent" })}
/>
{#if settings.userSettings[UserSettingKeys.AnimationDuration] != 0}
  <ToggleInput
    name={UserSettingKeys.AnimateCalendarSwipe}
    description={t("settings.appearance.animations.calendar.main")}
    bind:value={settings.userSettings[UserSettingKeys.AnimateCalendarSwipe]}
  />
  {#if settings.userSettings[UserSettingKeys.DisplaySmallCalendar]}
    <ToggleInput
      name={UserSettingKeys.AnimateSmallCalendarSwipe}
      description={t("settings.appearance.animations.calendar.small")}
      bind:value={settings.userSettings[UserSettingKeys.AnimateSmallCalendarSwipe]}
    />
  {/if}
  <ToggleInput
    name={UserSettingKeys.AnimateMonthSelectionSwipe}
    description={t("settings.appearance.animations.month")}
    bind:value={settings.userSettings[UserSettingKeys.AnimateMonthSelectionSwipe]}
  />
{/if}
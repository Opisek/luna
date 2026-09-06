<script lang="ts">
  import { Frequency, RRule, Weekday, type ByWeekday, type Options } from "rrule";
  import SelectInput from "./SelectInput.svelte";
  import ToggleInput from "./ToggleInput.svelte";
  import DateTimeInput from "./DateTimeInput.svelte";
  import { untrack } from "svelte";
  import SelectButtons from "./SelectButtons.svelte";
  import SelectButtonsMulti from "./SelectButtonsMulti.svelte";
  import { t } from "@sveltia/i18n";
  import { getDayName, getMonthName } from "$lib/common/humanization";
  import { getSettings } from "$lib/client/data/settings.svelte";
  import { UserSettingKeys } from "../../types/settings";
  import SelectInputMulti from "./SelectInputMulti.svelte";
  import TextInput from "./TextInput.svelte";
  import RadioInput from "./RadioInput.svelte";
  import type { Option } from "../../types/options";
  import NumberInput from "./NumberInput.svelte";

  interface Props {
    options: Partial<Options>;
    dtstart: Date;
    allDay: boolean;
    editable: boolean;
    simple?: boolean;
  }

  let {
    options = $bindable(),
    dtstart,
    allDay,
    editable,
    simple = false,
  }: Props = $props();

  const settings = getSettings();

  let byValsEnabled = $state({
    byMonthLimit: false,
    byYearDayLimit: false,
    byMonthDayLimit: false,
    byDayLimit: false,
    byHourLimit: false,
    byMinuteLimit: false,
    bySecondLimit: false,
    bySetPosLimit: false,
    byMonthExpand: false,
    byWeekNoExpand: false,
    byYearDayExpand: false,
    byMonthDayExpand: false,
    byDayExpand: false,
    byHourExpand: false,
    byMinuteExpand: false,
    bySecondExpand: false,
  })

  let byValsSemantics = $derived({
    byMonthLimit: options.freq != RRule.YEARLY,
    byYearDayLimit: options.freq !== undefined && [RRule.SECONDLY, RRule.MINUTELY, RRule.HOURLY].includes(options.freq),
    byMonthDayLimit: options.freq !== undefined && [RRule.SECONDLY, RRule.MINUTELY, RRule.HOURLY, RRule.DAILY].includes(options.freq),
    byDayLimit: options.freq !== undefined && [RRule.SECONDLY, RRule.MINUTELY, RRule.HOURLY, RRule.DAILY].includes(options.freq) || (options.freq == RRule.MONTHLY && byValsEnabled.byMonthDayExpand) || (options.freq == RRule.YEARLY && (byValsEnabled.byYearDayExpand || byValsEnabled.byMonthDayExpand)),
    byHourLimit: !allDay && options.freq !== undefined && [RRule.SECONDLY, RRule.MINUTELY, RRule.HOURLY].includes(options.freq),
    byMinuteLimit: !allDay && options.freq !== undefined && [RRule.SECONDLY, RRule.MINUTELY].includes(options.freq),
    bySecondLimit: !allDay && options.freq !== undefined && [RRule.SECONDLY].includes(options.freq),
    byMonthExpand: options.freq == RRule.YEARLY,
    byWeekNoExpand: options.freq == RRule.YEARLY,
    byYearDayExpand: options.freq == RRule.YEARLY,
    byMonthDayExpand: options.freq !== undefined && [RRule.MONTHLY, RRule.YEARLY].includes(options.freq),
    byDayExpand: options.freq == RRule.WEEKLY || (options.freq == RRule.MONTHLY && !byValsEnabled.byMonthDayExpand) || (options.freq == RRule.YEARLY && !byValsEnabled.byYearDayExpand && !byValsEnabled.byMonthDayExpand),
    byHourExpand: !allDay && options.freq !== undefined && [RRule.DAILY, RRule.WEEKLY, RRule.MONTHLY, RRule.YEARLY].includes(options.freq),
    byMinuteExpand: !allDay && options.freq !== undefined && [RRule.HOURLY, RRule.DAILY, RRule.WEEKLY, RRule.MONTHLY, RRule.YEARLY].includes(options.freq),
    bySecondExpand: !allDay && options.freq != RRule.SECONDLY,
    byDayNumeric: options.freq == RRule.MONTHLY || (options.freq == RRule.YEARLY && !byValsEnabled.byWeekNoExpand)
  })

  let byVals = $state({
    byMonthLimit: [] as number[],
    byYearDayLimit: [] as number[],
    byMonthDayLimit: [] as number[],
    byDayLimit: [] as string[],
    byHourLimit: [] as number[],
    byMinuteLimit: [] as number[],
    bySecondLimit: [] as number[],
    bySetPosLimit: [] as number[],
    byMonthExpand: [] as number[],
    byWeekNoExpand: [] as number[],
    byYearDayExpand: [] as number[],
    byMonthDayExpand: [] as number[],
    byDayExpand: [] as string[],
    byHourExpand: [] as number[],
    byMinuteExpand: [] as number[],
    bySecondExpand: [] as number[],
  })

  const weekdays = [RRule.SU, RRule.MO, RRule.TU, RRule.WE, RRule.TH, RRule.FR, RRule.SA];

  // Initialize values
  $effect(() => {
    options;
    untrack(() => {
      // Initialize BYMONTH
      if (options.bymonth == null) {
        byValsEnabled.byMonthExpand = false;
        byValsEnabled.byMonthLimit = false;
      } else {
        const arr = Array.isArray(options.bymonth) ? options.bymonth : [ options.bymonth ];
        if (byValsSemantics.byMonthExpand) {
          byValsEnabled.byMonthExpand = true;
          byVals.byMonthExpand = arr;
        } else if (byValsSemantics.byMonthLimit) {
          byValsEnabled.byMonthLimit = true;
          byVals.byMonthLimit = arr;
        }
      }

      // Initialize BYWEEKNO
      if (options.byweekno == null) {
        byValsEnabled.byWeekNoExpand = false;
      } else {
        const arr = Array.isArray(options.byweekno) ? options.byweekno : [ options.byweekno ];
        if (byValsSemantics.byWeekNoExpand) {
          byValsEnabled.byWeekNoExpand = true;
          byVals.byWeekNoExpand = arr;
        }
      }

      // Initialize BYYEARDAY
      if (options.byyearday == null) {
        byValsEnabled.byYearDayExpand = false;
        byValsEnabled.byYearDayLimit = false;
      } else {
        const arr = Array.isArray(options.byyearday) ? options.byyearday : [ options.byyearday ];
        if (byValsSemantics.byYearDayExpand) {
          byValsEnabled.byYearDayExpand = true;
          byVals.byYearDayExpand = arr;
        } else if (byValsSemantics.byYearDayLimit) {
          byValsEnabled.byYearDayLimit = true;
          byVals.byYearDayLimit = arr;
        }
      }

      // Initialize BYMONTHDAY
      if (options.bymonthday == null) {
        byValsEnabled.byMonthDayExpand = false;
        byValsEnabled.byMonthDayLimit = false;
      } else {
        const arr = Array.isArray(options.bymonthday) ? options.bymonthday : [ options.bymonthday ];
        if (byValsSemantics.byMonthDayExpand) {
          byValsEnabled.byMonthDayExpand = true;
          byVals.byMonthDayExpand = arr;
        } else if (byValsSemantics.byMonthDayLimit) {
          byValsEnabled.byMonthDayLimit = true;
          byVals.byMonthDayLimit = arr;
        }
      }

      // Initialize BYDAY
      if (options.byweekday == null && options.bynweekday == null) {
        byValsEnabled.byDayExpand = false;
        byValsEnabled.byDayLimit = false;
      } else {
        const arr1 = options.byweekday == null ? [] : (Array.isArray(options.byweekday) ? options.byweekday : [ options.byweekday ]);
        const arr2 = options.bynweekday == null ? [] : options.bynweekday.flatMap(x => weekdays[x[0]].nth(x[1]));
        const arr = arr1.concat(arr2).map(x => x.toString());
        if (byValsSemantics.byDayExpand) {
          byValsEnabled.byDayExpand = true;
          byVals.byDayExpand = arr;
        } else if (byValsSemantics.byDayLimit) {
          byValsEnabled.byDayLimit = true;
          byVals.byDayLimit = arr;
        }
      }

      // Initialize BYHOUR
      if (options.byhour == null) {
        byValsEnabled.byHourExpand = false;
        byValsEnabled.byHourLimit = false;
      } else {
        const arr = Array.isArray(options.byhour) ? options.byhour : [ options.byhour ];
        if (byValsSemantics.byHourExpand) {
          byValsEnabled.byHourExpand = true;
          byVals.byHourExpand = arr;
        } else if (byValsSemantics.byHourLimit) {
          byValsEnabled.byHourLimit = true;
          byVals.byHourLimit = arr;
        }
      }

      // Initialize BYMINUTE
      if (options.byminute == null) {
        byValsEnabled.byMinuteExpand = false;
        byValsEnabled.byMinuteLimit = false;
      } else {
        const arr = Array.isArray(options.byminute) ? options.byminute : [ options.byminute ];
        if (byValsSemantics.byMinuteExpand) {
          byValsEnabled.byMinuteExpand = true;
          byVals.byMinuteExpand = arr;
        } else if (byValsSemantics.byMinuteLimit) {
          byValsEnabled.byMinuteLimit = true;
          byVals.byMinuteLimit = arr;
        }
      }

      // Initialize BYSECOND
      if (options.bysecond == null) {
        byValsEnabled.bySecondExpand = false;
        byValsEnabled.bySecondLimit = false;
      } else {
        const arr = Array.isArray(options.bysecond) ? options.bysecond : [ options.bysecond ];
        if (byValsSemantics.bySecondExpand) {
          byValsEnabled.bySecondExpand = true;
          byVals.bySecondExpand = arr;
        } else if (byValsSemantics.bySecondLimit) {
          byValsEnabled.bySecondLimit = true;
          byVals.bySecondLimit = arr;
        }
      }

      // Initialize BYSETPOS
      if (options.bysetpos == null) {
        byValsEnabled.bySetPosLimit = false;
      } else {
        const arr = Array.isArray(options.bysetpos) ? options.bysetpos : [ options.bysetpos ];
        if (byValsSemantics.bySecondLimit) {
          byValsEnabled.bySetPosLimit = true;
          byVals.bySetPosLimit = arr;
        }
      }

      // Initialize COUNT/UNTIL
      if (options.count !== null && options.count !== undefined) {
        endType = "count";
      } else if (options.until !== null && options.until !== undefined) {
        endType = "date";
      } else {
        endType = "forever";
      }
    })
  })

  let yearlyType: "date" | "day" = $state("date");
  let endType: "date" | "count" | "forever" = $state("forever");

  $effect(() => {
    endType;
    untrack(() => {
      switch (endType) {
        case "date":
          options.count = null;
          if (options.until == null) options.until = new Date();
          break
        case "count":
          options.until = null;
          if (options.count == null) options.count = 1;
          break
        case "forever":
          options.count = null;
          options.until = null;
          break
      }
    })
  });

  // Set WKST
  $effect(() => {
    const wkst = settings.userSettings[UserSettingKeys.FirstDayOfWeek];
    untrack(() => {
      options.wkst = weekdays[wkst];
    })
  })

  // Set DTSTART
  //$effect(() => {
  //  dtstart;
  //  untrack(() => {
  //    options.dtstart = dtstart;
  //  })
  //})

  // Set BYMONTH
  $effect(() => {
    byValsSemantics.byMonthExpand,
    byValsSemantics.byMonthLimit,
    byValsEnabled.byMonthExpand,
    byValsEnabled.byMonthLimit,
    byVals.byMonthExpand,
    byVals.byMonthLimit
    untrack(() => {
      if (byValsSemantics.byMonthExpand && byValsEnabled.byMonthExpand) options.bymonth = byVals.byMonthExpand;
      else if (byValsSemantics.byMonthLimit && byValsEnabled.byMonthLimit) options.bymonth = byVals.byMonthLimit;
      else options.bymonth = null;
    })
  })

  // Set BYWEEKNO
  $effect(() => {
    byValsSemantics.byWeekNoExpand,
    byValsEnabled.byWeekNoExpand,
    byVals.byWeekNoExpand,
    untrack(() => {
      if (byValsSemantics.byWeekNoExpand && byValsEnabled.byWeekNoExpand) options.byweekno = byVals.byWeekNoExpand;
      else options.byweekno = null;
    })
  })

  // Set BYYEARDAY
  $effect(() => {
    byValsSemantics.byYearDayExpand,
    byValsSemantics.byYearDayLimit,
    byValsEnabled.byYearDayExpand,
    byValsEnabled.byYearDayLimit,
    byVals.byYearDayExpand,
    byVals.byYearDayLimit
    untrack(() => {
      if (byValsSemantics.byYearDayExpand && byValsEnabled.byYearDayExpand) options.byyearday = byVals.byYearDayExpand;
      else if (byValsSemantics.byYearDayLimit && byValsEnabled.byYearDayLimit) options.byyearday = byVals.byYearDayLimit;
      else options.byyearday = null;
    })
  })

  // Set BYMONTHDAY
  $effect(() => {
    byValsSemantics.byMonthDayExpand,
    byValsSemantics.byMonthDayLimit,
    byValsEnabled.byMonthDayExpand,
    byValsEnabled.byMonthDayLimit,
    byVals.byMonthDayExpand,
    byVals.byMonthDayLimit
    untrack(() => {
      options.bynmonthday = null;
      if (byValsSemantics.byMonthDayExpand && byValsEnabled.byMonthDayExpand) options.bymonthday = byVals.byMonthDayExpand;
      else if (byValsSemantics.byMonthDayLimit && byValsEnabled.byMonthDayLimit) options.bymonthday = byVals.byMonthDayLimit;
      else options.bymonthday = null;
    })
  })

  // Set BYDAY
  const nWeekdayRegex = /(?<num>(\+|-)?\d+)?(?<day>MO|TU|WE|TH|FR|SA|SU)/;
  function stringToNWeekday(str: string): Weekday {
    const matches = nWeekdayRegex.exec(str);
    if (matches == null || !matches.groups || !("day" in matches.groups)) return RRule.MO;
    const day = matches.groups["day"];
    const baseDay = { "MO": RRule.MO, "TU": RRule.TU, "WE": RRule.WE, "TH": RRule.TH, "FR": RRule.FR, "SA": RRule.SA, "SU": RRule.SU }[day];
    if (!baseDay) return RRule.MO;
    if (!("num" in matches.groups)) return baseDay;
    const num = Number.parseInt(matches.groups["num"]);
    if (num == 0) return baseDay;
    return baseDay.nth(num);
  }

  $effect(() => {
    byValsSemantics.byDayExpand,
    byValsSemantics.byDayLimit,
    byValsEnabled.byDayExpand,
    byValsEnabled.byDayLimit,
    byVals.byDayExpand,
    byVals.byDayLimit
    untrack(() => {
      options.bynweekday = null;
      if (byValsSemantics.byDayExpand && byValsEnabled.byDayExpand) options.byweekday = byVals.byDayExpand.map(x => stringToNWeekday(x));
      else if (byValsSemantics.byDayLimit && byValsEnabled.byDayLimit) options.byweekday = byVals.byDayLimit.map(x => stringToNWeekday(x));
      else options.byweekday = null;
    })
  })

  // Set BYHOUR
  $effect(() => {
    byValsSemantics.byHourExpand,
    byValsSemantics.byHourLimit,
    byValsEnabled.byHourExpand,
    byValsEnabled.byHourLimit,
    byVals.byHourExpand,
    byVals.byHourLimit
    untrack(() => {
      if (byValsSemantics.byHourExpand && byValsEnabled.byHourExpand) options.byhour = byVals.byHourExpand;
      else if (byValsSemantics.byHourLimit && byValsEnabled.byHourLimit) options.byhour = byVals.byHourLimit;
      else options.byhour = null;
    })
  })

  // Set BYMINUTE
  $effect(() => {
    byValsSemantics.byMinuteExpand,
    byValsSemantics.byMinuteLimit,
    byValsEnabled.byMinuteExpand,
    byValsEnabled.byMinuteLimit,
    byVals.byMinuteExpand,
    byVals.byMinuteLimit
    untrack(() => {
      if (byValsSemantics.byMinuteExpand && byValsEnabled.byMinuteExpand) options.byminute = byVals.byMinuteExpand;
      else if (byValsSemantics.byMinuteLimit && byValsEnabled.byMinuteLimit) options.byminute = byVals.byMinuteLimit;
      else options.byminute = null;
    })
  })
  
  // Set BYSECOND
  $effect(() => {
    byValsSemantics.bySecondExpand,
    byValsSemantics.bySecondLimit,
    byValsEnabled.bySecondExpand,
    byValsEnabled.bySecondLimit,
    byVals.bySecondExpand,
    byVals.bySecondLimit
    untrack(() => {
      if (byValsSemantics.bySecondExpand && byValsEnabled.bySecondExpand) options.bysecond = byVals.bySecondExpand;
      else if (byValsSemantics.bySecondLimit && byValsEnabled.bySecondLimit) options.bysecond = byVals.bySecondLimit;
      else options.bysecond = null;
    })
  })

  // Presets
  enum RecurrencePreset {
    Daily,
    DailyWeekdays,
    Weekly,
    MonthlyDate,
    MonthlyDay,
    MonthlyDayReverse,
    YearlyDate,
    YearlyMonthDay,
    YearlyMonthDayReverse,
  }
  let nthMonthDay = $derived(Math.ceil(dtstart.getDate() / 7))
  let nthReverseMonthDay = $derived.by(() => {
    const lastDayOfMonth = new Date(dtstart);
    lastDayOfMonth.setMonth(lastDayOfMonth.getMonth() + 1);
    lastDayOfMonth.setDate(0);
    return Math.max(Math.ceil((lastDayOfMonth.getDate() - dtstart.getDate() + 1) / 7), 1);
  })
  let recurrencePresetNameParameters = $derived({ values: { interval: options.interval || 1, date: dtstart, day: dtstart.getDate(), nth: nthMonthDay, nthrev: nthReverseMonthDay } });
  let recurrencePresetNames: Record<RecurrencePreset, string> = $derived({
    [RecurrencePreset.Daily]: t("recurrence.presets.daily", recurrencePresetNameParameters),
    [RecurrencePreset.DailyWeekdays]: t("recurrence.presets.dailyweekdays", recurrencePresetNameParameters),
    [RecurrencePreset.Weekly]: t("recurrence.presets.weekly", recurrencePresetNameParameters),
    [RecurrencePreset.MonthlyDate]: t("recurrence.presets.monthlydate", recurrencePresetNameParameters),
    [RecurrencePreset.MonthlyDay]: t("recurrence.presets.monthlyday", recurrencePresetNameParameters),
    [RecurrencePreset.MonthlyDayReverse]: t("recurrence.presets.monthlydayreverse", recurrencePresetNameParameters),
    [RecurrencePreset.YearlyDate]: t("recurrence.presets.yearlydate", recurrencePresetNameParameters),
    [RecurrencePreset.YearlyMonthDay]: t("recurrence.presets.yearlymonthday", recurrencePresetNameParameters),
    [RecurrencePreset.YearlyMonthDayReverse]: t("recurrence.presets.yearlymonthdayreverse", recurrencePresetNameParameters),
  })
  let applicableRecurrencePresets: Record<Frequency, RecurrencePreset[]> = {
    [RRule.SECONDLY]: [],
    [RRule.MINUTELY]: [],
    [RRule.HOURLY]: [],
    [RRule.DAILY]: [ RecurrencePreset.Daily, RecurrencePreset.DailyWeekdays ],
    [RRule.WEEKLY]: [ RecurrencePreset.Weekly ],
    [RRule.MONTHLY]: [ RecurrencePreset.MonthlyDate, RecurrencePreset.MonthlyDay, RecurrencePreset.MonthlyDayReverse ],
    [RRule.YEARLY]: [ RecurrencePreset.YearlyDate, RecurrencePreset.YearlyMonthDay, RecurrencePreset.YearlyMonthDayReverse ],
  }
  let recurrencePresets: Record<RecurrencePreset, Partial<Options>> = $derived({
    [RecurrencePreset.Daily]: { dtstart: dtstart, freq: RRule.DAILY },
    [RecurrencePreset.DailyWeekdays]: { dtstart: dtstart, freq: RRule.DAILY, byweekday: [ RRule.MO, RRule.TU, RRule.WE, RRule.TH, RRule.FR ] },
    [RecurrencePreset.Weekly]: { dtstart: dtstart, freq: RRule.WEEKLY },
    [RecurrencePreset.MonthlyDate]: { dtstart: dtstart, freq: RRule.MONTHLY, bymonthday: dtstart.getDate() },
    [RecurrencePreset.MonthlyDay]: { dtstart: dtstart, freq: RRule.MONTHLY, byweekday: weekdays[dtstart.getDay()].nth(nthMonthDay) },
    [RecurrencePreset.MonthlyDayReverse]: { dtstart: dtstart, freq: RRule.MONTHLY, byweekday: weekdays[dtstart.getDay()].nth(-nthReverseMonthDay) },
    [RecurrencePreset.YearlyDate]: { dtstart: dtstart, freq: RRule.YEARLY, bymonth: dtstart.getMonth() + 1, bymonthday: dtstart.getDate() },
    [RecurrencePreset.YearlyMonthDay]: { dtstart: dtstart, freq: RRule.YEARLY, bymonth: dtstart.getMonth() + 1, byweekday: weekdays[dtstart.getDay()].nth(nthMonthDay) },
    [RecurrencePreset.YearlyMonthDayReverse]: { dtstart: dtstart, freq: RRule.YEARLY, bymonth: dtstart.getMonth() + 1, byweekday: weekdays[dtstart.getDay()].nth(-nthReverseMonthDay) },
  });
  let chosenRecurrencePreset: RecurrencePreset | null = $state(null);
  let previousDtStart: Date | null = $state(null);

  function applyPreset(chosen: RecurrencePreset | null) {
    if (chosen == null) return;

    options.bymonth = null;
    options.byweekno = null;
    options.byyearday = null;
    options.bymonthday = null;
    options.bynmonthday = null;
    options.byweekday = null;
    options.bynweekday = null;
    options.byhour = null;
    options.byminute = null;
    options.bysecond = null;
    options.bysetpos = null;

    for (const [key, value] of Object.entries(recurrencePresets[chosen])) {
      //@ts-ignore
      options[key] = value;
    }
  }

  const presetRelevantOptions = [
    "freq", "bysetpos", "bymonth", "bymonthday", "bynmonthday", "byyearday", "byweekno", "byweekday", "bynweekday", "byhour", "byminute", "bysecond", "byeaster"
  ]

  // Map recurrence to preset
  $effect(() => {
    if (previousDtStart != dtstart && chosenRecurrencePreset != null) {
      applyPreset(chosenRecurrencePreset);
      previousDtStart = dtstart;
      return;
    }

    const normalizedCurrent = new RRule(options);

    const matches = Object
      .entries(recurrencePresets)
      .filter(preset => {
        const normalizedPreset = new RRule(preset[1]);
        //@ts-ignore
        return presetRelevantOptions.every(optionKey => JSON.stringify(normalizedCurrent.options[optionKey]) === JSON.stringify(normalizedPreset.options[optionKey]));
      });

    chosenRecurrencePreset = (matches.length == 0 ? null : matches[0][0]) as (RecurrencePreset | null);
    previousDtStart = dtstart;
  });
</script>

<!-- FREQ -->
<SelectInput bind:value={options.freq} name="recurrence_freq" placeholder={t("recurrence.props.freq.name")} showLabel={true} options={
  (simple || allDay ? [] : [
    { value: RRule.SECONDLY, name: t("recurrence.props.freq.vals.secondly") },
    { value: RRule.MINUTELY, name: t("recurrence.props.freq.vals.minutely") },
    { value: RRule.HOURLY, name: t("recurrence.props.freq.vals.hourly") },
  ]).concat([
    { value: RRule.DAILY, name: t("recurrence.props.freq.vals.daily") },
    { value: RRule.WEEKLY, name: t("recurrence.props.freq.vals.weekly") },
    { value: RRule.MONTHLY, name: t("recurrence.props.freq.vals.monthly") },
    { value: RRule.YEARLY, name: t("recurrence.props.freq.vals.yearly") },
  ]
)} editable={editable} />

<NumberInput
  placeholder={t("recurrence.props.interval.name")}
  name="recurrence_interval"
  min={1}
  bind:value={options.interval}
/>

<!--
{#if options.freq == RRule.YEARLY}
  <SelectButtons bind:value={yearlyType} name="recurrence_freq_yearly_type" placeholder="By" options={[
    { value: "date", name: "By day number" },
    { value: "day", name: "By week day" },
  ]} editable={editable} />
{/if}
-->

{#if simple}
  {@const presets = options.freq === undefined ? [] : applicableRecurrencePresets[options.freq]}
  {@const mappedPresets = presets.map(x => ({
    name: recurrencePresetNames[x],
    value: x,
  }))}
  {@const mappedPresetsAndCustom = (mappedPresets as Option<RecurrencePreset | null>[]).concat(chosenRecurrencePreset !== null ? [] : [{
    name: "Custom",
    value: chosenRecurrencePreset,
  }])}
  {#if mappedPresetsAndCustom.length > 1}
    <RadioInput
      bind:value={chosenRecurrencePreset}
      name="recurrence_preset"
      options={mappedPresetsAndCustom}
      onClick={applyPreset}
    />
  {/if}
{:else}
  <!-- BYMONTH -->
  {#if byValsSemantics.byMonthLimit }
    <ToggleInput
      bind:value={byValsEnabled.byMonthLimit}
      name="recurrence_bymonth_limit_enable" 
      description={t("recurrence.props.by.month.limit.enable")}
    />
    {#if byValsEnabled.byMonthLimit}
      <SelectButtonsMulti
        bind:values={byVals.byMonthLimit}
        name="recurrence_bymonth_limit"
        placeholder={t("recurrence.props.by.month.limit.name")}
        options={ [...Array(12).keys()].map(x => ({ value: x+1, name: getMonthName(x, true) })) }
        onClick={() => { byVals.byMonthLimit = $state.snapshot(byVals.byMonthLimit); }}
      />
    {/if}
  {/if}
  {#if byValsSemantics.byMonthExpand }
    <ToggleInput
      bind:value={byValsEnabled.byMonthExpand}
      name="recurrence_bymonth_expand_enable" 
      description={t("recurrence.props.by.month.expand.enable")}
    />
    {#if byValsEnabled.byMonthExpand}
      <SelectButtonsMulti
        bind:values={byVals.byMonthExpand}
        name="recurrence_bymonth_expand"
        placeholder={t("recurrence.props.by.month.expand.name")}
        options={ [...Array(12).keys()].map(x => ({ value: x+1, name: getMonthName(x, true) })) }
        onClick={() => { byVals.byMonthExpand = $state.snapshot(byVals.byMonthExpand); }}
      />
    {/if}
  {/if}

  <!-- BYWEEKNO -->
  {#if byValsSemantics.byWeekNoExpand }
    <ToggleInput
      bind:value={byValsEnabled.byWeekNoExpand}
      name="recurrence_byweekno_expand_enable" 
      description={t("recurrence.props.by.weekno.expand.enable")}
    />
    {#if byValsEnabled.byWeekNoExpand}
      <SelectInputMulti
        bind:values={byVals.byWeekNoExpand}
        name="recurrence_byweekno_expand"
        placeholder={t("recurrence.props.by.weekno.expand.name")}
        options={
          [...Array(53).keys()].map(x => ({ value: x+1, name: t("numbers.ordinal.normal", { values: { num: x + 1 } }) })).concat(
          [...Array(53).keys()].map(x => ({ value: -(x+1), name: t("numbers.ordinal.reverse", { values: { num: x + 1 } }) })))
        }
        click={() => { byVals.byWeekNoExpand = $state.snapshot(byVals.byWeekNoExpand); }}
      />
    {/if}
  {/if}

  <!-- BYYEARDAY -->
  {#if byValsSemantics.byYearDayLimit }
    <ToggleInput
      bind:value={byValsEnabled.byYearDayLimit}
      name="recurrence_byyearday_limit_enable" 
      description={t("recurrence.props.by.yearday.limit.enable")}
    />
    {#if byValsEnabled.byYearDayLimit}
      <SelectInputMulti
        bind:values={byVals.byMonthDayLimit}
        name="recurrence_byyearday_limit"
        placeholder={t("recurrence.props.by.year.limit.name")}
        options={
          [...Array(366).keys()].map(x => ({ value: x+1, name: t("numbers.ordinal.normal", { values: { num: x + 1 } }) })).concat(
          [...Array(366).keys()].map(x => ({ value: -(x+1), name: t("numbers.ordinal.reverse", { values: { num: x + 1 } }) })))
        }
        click={() => { byVals.byYearDayLimit = $state.snapshot(byVals.byYearDayLimit); }}
      />
    {/if}
  {/if}
  {#if byValsSemantics.byYearDayExpand }
    <ToggleInput
      bind:value={byValsEnabled.byYearDayExpand}
      name="recurrence_byyearday_limit_enable" 
      description={t("recurrence.props.by.yearday.expand.enable")}
    />
    {#if byValsEnabled.byYearDayExpand}
      <SelectInputMulti
        bind:values={byVals.byMonthDayExpand}
        name="recurrence_byyearday_expand"
        placeholder={t("recurrence.props.by.yearday.expand.name")}
        options={
          [...Array(366).keys()].map(x => ({ value: x+1, name: t("numbers.ordinal.normal", { values: { num: x + 1 } }) })).concat(
          [...Array(366).keys()].map(x => ({ value: -(x+1), name: t("numbers.ordinal.reverse", { values: { num: x + 1 } }) })))
        }
        click={() => { byVals.byYearDayExpand = $state.snapshot(byVals.byYearDayExpand); }}
      />
    {/if}
  {/if}

  <!-- BYMONTHDAY -->
  {#if byValsSemantics.byMonthDayLimit }
    <ToggleInput
      bind:value={byValsEnabled.byMonthDayLimit}
      name="recurrence_bymonthday_limit_enable" 
      description={t("recurrence.props.by.monthday.limit.enable")}
    />
    {#if byValsEnabled.byMonthDayLimit}
      <SelectInputMulti
        bind:values={byVals.byMonthDayLimit}
        name="recurrence_bymonthday_limit"
        placeholder={t("recurrence.props.by.monthday.limit.name")}
        options={
          [...Array(31).keys()].map(x => ({ value: x+1, name: t("numbers.ordinal.normal", { values: { num: x + 1 } }) })).concat(
          [...Array(31).keys()].map(x => ({ value: -(x+1), name: t("numbers.ordinal.reverse", { values: { num: x + 1 } }) })))
        }
        click={() => { byVals.byMonthDayLimit = $state.snapshot(byVals.byMonthDayLimit); }}
      />
    {/if}
  {/if}
  {#if byValsSemantics.byMonthDayExpand }
    <ToggleInput
      bind:value={byValsEnabled.byMonthDayExpand}
      name="recurrence_bymonthday_expand_enable" 
      description={t("recurrence.props.by.monthday.expand.enable")}
    />
    {#if byValsEnabled.byMonthDayExpand}
      <SelectInputMulti
        bind:values={byVals.byMonthDayExpand}
        name="recurrence_bymonthday_expand"
        placeholder={t("recurrence.props.by.monthday.expand.name")}
        options={
          [...Array(31).keys()].map(x => ({ value: x+1, name: t("numbers.ordinal.normal", { values: { num: x + 1 } }) })).concat(
          [...Array(31).keys()].map(x => ({ value: -(x+1), name: t("numbers.ordinal.reverse", { values: { num: x + 1 } }) })))
        }
        click={() => { byVals.byMonthDayExpand = $state.snapshot(byVals.byMonthDayExpand); }}
      />
    {/if}
  {/if}

  <!-- BYDAY -->
  {#if byValsSemantics.byDayLimit }
    <ToggleInput
      bind:value={byValsEnabled.byDayLimit}
      name="recurrence_byday_limit_enable" 
      description={t(`recurrence.props.by.day.${byValsSemantics.byDayNumeric ? "numeric" : "dayname"}.limit.enable`)}
    />
    {#if byValsEnabled.byDayLimit}
      {#if byValsSemantics.byDayNumeric}
        <SelectInputMulti
          bind:values={byVals.byDayLimit}
          name="recurrence_byday_limit"
          placeholder={t("recurrence.props.by.day.numeric.limit.name")}
          options={
            [RRule.SU, RRule.MO, RRule.TU, RRule.WE, RRule.TH, RRule.FR, RRule.SA]
              .map((x, i) => ({ value: x, name: getDayName(i, false), index: (i + 7 - settings.userSettings[UserSettingKeys.FirstDayOfWeek]) % 7 }))
              .sort((a, b) =>  a.index-b.index)
              .flatMap(x =>
                [{ value: x.value, name: `Every ${x.name}`}].concat(
                  [...Array(4).keys()].map(y => ({ value: x.value.nth(y+1), name: `${t("numbers.ordinal.normal", { values: { num: y + 1 } })} ${x.name}` }))
                ).concat(
                  [...Array(4).keys()].map(y => ({ value: x.value.nth(-(y+1)), name: `${t("numbers.ordinal.reverse", { values: { num: y + 1 } })} ${x.name}` }))
                )
              )
              .map(x => ({ value: x.value.toString(), name: x.name }))
          }
          click={() => { byVals.byDayLimit = $state.snapshot(byVals.byDayLimit); }}
        />
      {:else}
        <SelectButtonsMulti
          bind:values={byVals.byDayLimit}
          name="recurrence_byday_limit"
          placeholder={t("recurrence.props.by.day.dayname.limit.name")}
          options={
            [RRule.SU, RRule.MO, RRule.TU, RRule.WE, RRule.TH, RRule.FR, RRule.SA]
              .map((x, i) => ({ value: x, name: getDayName(i, true), index: (i + 7 - settings.userSettings[UserSettingKeys.FirstDayOfWeek]) % 7 }))
              .sort((a, b) =>  a.index-b.index)
              .map(x => ({ value: x.value.toString(), name: x.name }))
          }
          onClick={() => { byVals.byDayLimit = $state.snapshot(byVals.byDayLimit); }}
        />
      {/if}
    {/if}
  {/if}
  {#if byValsSemantics.byDayExpand }
    <ToggleInput
      bind:value={byValsEnabled.byDayExpand}
      name="recurrence_byday_expand_enable" 
      description={t(`recurrence.props.by.day.${byValsSemantics.byDayNumeric ? "numeric" : "dayname"}.expand.enable`)}
    />
    {#if byValsEnabled.byDayExpand}
      {#if byValsSemantics.byDayNumeric}
        <SelectInputMulti
          bind:values={byVals.byDayExpand}
          name="recurrence_byday_expand"
          placeholder={t("recurrence.props.by.day.numeric.expand.name")}
          options={
            [RRule.SU, RRule.MO, RRule.TU, RRule.WE, RRule.TH, RRule.FR, RRule.SA]
              .map((x, i) => ({ value: x, name: getDayName(i, false), index: (i + 7 - settings.userSettings[UserSettingKeys.FirstDayOfWeek]) % 7 }))
              .sort((a, b) =>  a.index-b.index)
              .flatMap(x =>
                [{ value: x.value, name: `Every ${x.name}`}].concat(
                  [...Array(4).keys()].map(y => ({ value: x.value.nth(y+1), name: `${t("numbers.ordinal.normal", { values: { num: y + 1 } })} ${x.name}` }))
                ).concat(
                  [...Array(4).keys()].map(y => ({ value: x.value.nth(-(y+1)), name: `${t("numbers.ordinal.reverse", { values: { num: y + 1 } })} ${x.name}` }))
                )
              )
              .map(x => ({ value: x.value.toString(), name: x.name }))
          }
          click={() => { byVals.byDayExpand = $state.snapshot(byVals.byDayExpand); }}
        />
      {:else}
        <SelectButtonsMulti
          bind:values={byVals.byDayExpand}
          name="recurrence_byday_expand"
          placeholder={t("recurrence.props.by.day.dayname.expand.name")}
          options={
            [RRule.SU, RRule.MO, RRule.TU, RRule.WE, RRule.TH, RRule.FR, RRule.SA]
              .map((x, i) => ({ value: x, name: getDayName(i, true), index: (i + 7 - settings.userSettings[UserSettingKeys.FirstDayOfWeek]) % 7 }))
              .sort((a, b) =>  a.index-b.index)
              .map(x => ({ value: x.value.toString(), name: x.name }))
          }
          onClick={() => { byVals.byDayExpand = $state.snapshot(byVals.byDayExpand); }}
        />
      {/if}
    {/if}
  {/if}

  {#if !allDay}
    <!-- BYHOUR -->
    {#if byValsSemantics.byHourLimit }
      <ToggleInput
        bind:value={byValsEnabled.byHourLimit}
        name="recurrence_byhour_limit_enable" 
        description={t("recurrence.props.by.hour.limit.enable")}
      />
      {#if byValsEnabled.byHourLimit}
        <SelectInputMulti
          bind:values={byVals.byHourLimit}
          name="recurrence_byhour_limit"
          placeholder={t("recurrence.props.by.hour.limit.name")}
          options={
            [...Array(24).keys()].map(x => ({ value: x, name: t("numbers.ordinal.normal", { values: { num: x } }) }))
          }
          click={() => { byVals.byHourLimit = $state.snapshot(byVals.byHourLimit); }}
        />
      {/if}
    {/if}
    {#if byValsSemantics.byHourExpand }
      <ToggleInput
        bind:value={byValsEnabled.byHourExpand}
        name="recurrence_byhour_expand_enable" 
        description={t("recurrence.props.by.hour.expand.enable")}
      />
      {#if byValsEnabled.byHourExpand}
        <SelectInputMulti
          bind:values={byVals.byHourExpand}
          name="recurrence_byhour_expand"
          placeholder={t("recurrence.props.by.hour.expand.name")}
          options={
            [...Array(24).keys()].map(x => ({ value: x, name: t("numbers.ordinal.normal", { values: { num: x } }) }))
          }
          click={() => { byVals.byHourExpand = $state.snapshot(byVals.byHourExpand); }}
        />
      {/if}
    {/if}

    <!-- BYMINUTE -->
    {#if byValsSemantics.byMinuteLimit }
      <ToggleInput
        bind:value={byValsEnabled.byMinuteLimit}
        name="recurrence_byminute_limit_enable" 
        description={t("recurrence.props.by.minute.limit.enable")}
      />
      {#if byValsEnabled.byMinuteLimit}
        <SelectInputMulti
          bind:values={byVals.byMinuteLimit}
          name="recurrence_byminute_limit"
          placeholder={t("recurrence.props.by.minute.limit.name")}
          options={
            [...Array(60).keys()].map(x => ({ value: x, name: t("numbers.ordinal.normal", { values: { num: x } }) }))
          }
          click={() => { byVals.byMinuteLimit = $state.snapshot(byVals.byMinuteLimit); }}
        />
      {/if}
    {/if}
    {#if byValsSemantics.byMinuteExpand }
      <ToggleInput
        bind:value={byValsEnabled.byMinuteExpand}
        name="recurrence_byminute_expand_enable" 
        description={t("recurrence.props.by.minute.expand.enable")}
      />
      {#if byValsEnabled.byMinuteExpand}
        <SelectInputMulti
          bind:values={byVals.byMinuteExpand}
          name="recurrence_byminute_expand"
          placeholder={t("recurrence.props.by.minute.expand.name")}
          options={
            [...Array(60).keys()].map(x => ({ value: x, name: t("numbers.ordinal.normal", { values: { num: x } }) }))
          }
          click={() => { byVals.byMinuteExpand = $state.snapshot(byVals.byMinuteExpand); }}
        />
      {/if}
    {/if}

    <!-- BYSECOND -->
    {#if byValsSemantics.bySecondLimit }
      <ToggleInput
        bind:value={byValsEnabled.bySecondLimit}
        name="recurrence_bysecond_limit_enable" 
        description={t("recurrence.props.by.second.limit.enable")}
      />
      {#if byValsEnabled.bySecondLimit}
        <SelectInputMulti
          bind:values={byVals.bySecondLimit}
          name="recurrence_bysecond_limit"
          placeholder={t("recurrence.props.by.second.limit.name")}
          options={
            [...Array(60).keys()].map(x => ({ value: x, name: t("numbers.ordinal.normal", { values: { num: x } }) }))
          }
          click={() => { byVals.bySecondLimit = $state.snapshot(byVals.bySecondLimit); }}
        />
      {/if}
    {/if}
    {#if byValsSemantics.bySecondExpand }
      <ToggleInput
        bind:value={byValsEnabled.bySecondExpand}
        name="recurrence_bysecond_expand_enable" 
        description={t("recurrence.props.by.second.expand.enable")}
      />
      {#if byValsEnabled.bySecondExpand}
        <SelectInputMulti
          bind:values={byVals.bySecondExpand}
          name="recurrence_bysecond_expand"
          placeholder={t("recurrence.props.by.second.expand.name")}
          options={
            [...Array(60).keys()].map(x => ({ value: x, name: t("numbers.ordinal.normal", { values: { num: x } }) }))
          }
          click={() => { byVals.bySecondExpand = $state.snapshot(byVals.bySecondExpand); }}      
        />
      {/if}
    {/if}
  {/if}
{/if}

<SelectButtons bind:value={endType} name="recurrence_duration" placeholder={t("recurrence.props.duration.name")} options={[
  { value: "date", name: t("recurrence.props.duration.vals.date.name") },
  { value: "count", name: t("recurrence.props.duration.vals.count.name") },
  { value: "forever", name: t("recurrence.props.duration.vals.forever") },
]} editable={editable} />

{#if endType == "date"}
  <DateTimeInput
    bind:value={options.until}
    editable={editable}
    placeholder={t("recurrence.props.duration.vals.date.param")}
    name="recurrence_until" 
    allDay={allDay}
  />
{:else if endType == "count"}
  <NumberInput
    placeholder={t("recurrence.props.duration.vals.count.param")}
    name="recurrence_count"
    min={1}
    bind:value={options.count}
  />
{/if}

{#if !simple}
  <TextInput
    placeholder={t("recurrence.props.rrule")}
    name="recurrence_rrule"
    value={options ? RRule.optionsToString(options).split("RRULE:")[1] : ""}
    onChange={(x) => {
      const parts = x.split("RRULE:");
      const ruleStr = parts[parts.length - 1];
      const newOptions = RRule.fromString(ruleStr).origOptions;
      //@ts-ignore
      options = newOptions;
    }}
  />
{/if}
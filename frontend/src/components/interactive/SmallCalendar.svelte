<script lang="ts">
  import { NoOp } from "$lib/client/placeholders";
  import { focusIndicator } from "$lib/client/decoration";
  import { getDayIndex, isSameDay } from "$lib/common/date";
  import { UserSettingKeys } from "../../types/settings";
  import { getSettings } from "$lib/client/settings.svelte";

  interface Props {
    date: Date;
    onDayClick?: (date: Date) => any;
    smaller?: boolean;
  }

  let {
    date = $bindable(new Date()),
    onDayClick = NoOp,
    smaller = false,
  }: Props = $props();

  const settings = getSettings();

  let today = new Date();

  let [days, amountOfRows] = $derived((() => {
    // Date calculation
    const firstMonthDay = new Date(date.getFullYear(), date.getMonth(), 1);
    const lastMonthDay = new Date(date.getFullYear(), date.getMonth() + 1, 0);
    const firstDayOfWeek = getDayIndex(firstMonthDay);

    //amountOfRows = ;
    const amountOfRows = 
      settings.userSettings[UserSettingKeys.DynamicSmallCalendarRows] ?
      Math.ceil((lastMonthDay.getDate() + firstDayOfWeek) / 7)
      : 6;

    const firstViewDay = new Date(firstMonthDay);
    firstViewDay.setDate(firstMonthDay.getDate() - firstDayOfWeek);
    const lastViewDay = new Date(firstMonthDay);
    lastViewDay.setDate(firstMonthDay.getDate() + 7 * amountOfRows - 1);

    // Fill
    const days = [];

    const dateIterator = new Date(firstViewDay);

    for (let i = 0; i < 7 * amountOfRows; i++) {
      days.push(new Date(dateIterator));
      dateIterator.setDate(dateIterator.getDate() + 1);
    }

    return [days, amountOfRows];
  })());
</script>

<style lang="scss">
  @use "../../styles/animations.scss";
  @use "../../styles/colors.scss";
  @use "../../styles/dimensions.scss";
  @use "../../styles/text.scss";

  div.calendar {
    display: grid;
    grid-template-columns: repeat(7, 1fr);
    gap: dimensions.$gapSmall; 
  }

  div.smaller {
    font-size: text.$fontSizeSmall;
    gap: dimensions.$gapSmaller; 
  }

  button.day {
    all: unset;
    display: flex;
    justify-content: center;
    align-items: center;
    border-radius: dimensions.$borderRadiusSmall;
    color: colors.$foregroundSecondary;
    background-color: colors.$backgroundSecondary;
    padding: dimensions.$gapSmaller;
    cursor: pointer;
    user-select: none;
    position: relative;
    overflow: hidden;
  }

  button.day.sunday {
    color: colors.$foregroundSunday;
  }

  button.day.today {
    background-color: colors.$backgroundAccent;
    color: colors.$foregroundAccent;
  }

  button.day.otherMonth {
    opacity: 0.5;
  }
</style>

<div class="calendar" class:smaller={smaller} style="grid-template-rows: repeat({amountOfRows}, 1fr)">
  {#each days as day}
    <button
      class="day"
      class:sunday={day.getDay() == 0}
      class:today={isSameDay(day, today)}
      class:otherMonth={day.getMonth() != date.getMonth()}
      type="button"
      onclick={() => (onDayClick(day))}
      use:focusIndicator
    >
      {day.getDate()}
    </button>
  {/each}
</div>
<script lang="ts">
  import { NoOp } from "$lib/client/placeholders";
  import { focusIndicator } from "$lib/client/decoration";

  interface Props {
    date: Date;
    onDayClick?: (date: Date) => any;
  }

  let {
    date = $bindable(new Date()),
    onDayClick = NoOp
  }: Props = $props();

  let [days, amountOfRows] = $derived((() => {
    // Date calculation
    const firstMonthDay = new Date(date.getFullYear(), date.getMonth(), 1);
    //const lastMonthDay = new Date(year, month + 2, 0);
    const firstDayOfWeek = (firstMonthDay.getDay() + 6) % 7;

    //amountOfRows = Math.ceil((lastMonthDay.getDate() + firstDayOfWeek) / 7);
    const amountOfRows = 6;

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
  @import "../../styles/animations.scss";
  @import "../../styles/colors.scss";
  @import "../../styles/dimensions.scss";

  div.calendar {
    display: grid;
    grid-template-columns: repeat(7, 1fr);
    gap: $gapSmall; 
  }

  button.day {
    all: unset;
    display: flex;
    justify-content: center;
    align-items: center;
    border-radius: $borderRadiusSmall;
    color: $foregroundSecondary;
    background-color: $backgroundSecondary;
    padding: $gapSmaller;
    cursor: pointer;
    user-select: none;
    position: relative;
    overflow: hidden;
  }

  button.day.sunday {
    color: $foregroundSunday;
  }

  button.day.otherMonth {
    opacity: 0.5;
  }
</style>

<div class="calendar" style="grid-template-rows: repeat({amountOfRows}, 1fr)">
  {#each days as day, i}
    <button
      class="day"
      class:sunday={day.getDay() == 0}
      class:otherMonth={day.getMonth() != date.getMonth()}
      type="button"
      onclick={() => (onDayClick(day))}
      use:focusIndicator
    >
      {day.getDate()}
    </button>
  {/each}
</div>
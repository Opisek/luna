<script lang="ts">
  export let month: number;
  export let year: number;

  export let onDayClick: (date: Date) => any = () => {};

  let days: Date[] = [];
  let amountOfRows: number = 0;

  $: ((month: number, year: number) => {
    // Date calculation
    const firstMonthDay = new Date(year, month, 1);
    const lastMonthDay = new Date(year, month + 1, 0);
    const firstDayOfWeek = (firstMonthDay.getDay() + 6) % 7;

    //amountOfRows = Math.ceil((lastMonthDay.getDate() + firstDayOfWeek) / 7);
    amountOfRows = 6;

    const firstViewDay = new Date(firstMonthDay);
    firstViewDay.setDate(firstMonthDay.getDate() - firstDayOfWeek);
    const lastViewDay = new Date(firstMonthDay);
    lastViewDay.setDate(firstMonthDay.getDate() + 7 * amountOfRows - 1);

    // Fill
    days = [];

    const dateIterator = new Date(firstViewDay);

    for (let i = 0; i < 7 * amountOfRows; i++) {
      days.push(new Date(dateIterator));
      dateIterator.setDate(dateIterator.getDate() + 1);
    }
  })(month, year);
</script>

<style lang="scss">
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
    padding: $paddingTiny;
    cursor: pointer;
    user-select: none;
  }

  button.day.sunday {
    color: $foregroundSunday;
  }

  button.day.otherMonth {
    opacity: 0.5;
  }
</style>

<div class="calendar" style="grid-template-rows: repeat({amountOfRows}, 1fr)">
  {#each days as day}
    <button
      class="day"
      class:sunday={day.getDay() == 0}
      class:otherMonth={day.getMonth() != month}
      on:click={() => (onDayClick(day))}
    >
      {day.getDate()}
    </button>
  {/each}
</div>
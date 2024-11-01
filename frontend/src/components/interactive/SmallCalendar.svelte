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

  let clickedDay = -1;
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
    padding: $paddingTiny;
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

  div.focus {
    background-color: $backgroundAccent;
    height: 100%;
    width: $borderActiveWidth;
    position: absolute;
    left: 0;
    top: 0;
    transform: translateX(-100%);
    transition: transform $animationSpeedFast linear;
  }

  button.day:focus:not(.click) > div.focus {
    transform: translateX(0);
  }
</style>

<div class="calendar" style="grid-template-rows: repeat({amountOfRows}, 1fr)">
  {#each days as day, i}
    <button
      class="day"
      class:sunday={day.getDay() == 0}
      class:otherMonth={day.getMonth() != month}
      class:click={clickedDay === i}
      type="button"
      on:click={() => (onDayClick(day))}
      on:mousedown={() => (clickedDay = i)}
      on:focusout={() => {if (clickedDay === i) clickedDay = -1;}}
    >
      {day.getDate()}
      <div class="focus"></div>
    </button>
  {/each}
</div>
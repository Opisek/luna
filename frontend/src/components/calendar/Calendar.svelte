<script lang="ts">
    import { getDayName } from "../../lib/common/humanization";
    import Day from "./Day.svelte";

  export let month: number;
  export let year: number;

  let days: Date[] = [];

  $: ((month: number, year: number) => {
    const firstDay = new Date(year, month, 1);
    const lastDay = new Date(year, month + 1, 0);
    const firstDayOfWeek = (firstDay.getDay() + 6) % 7;

    console.log(firstDay, lastDay, firstDayOfWeek);

    const amountOfRows = Math.ceil((lastDay.getDate() + firstDayOfWeek - 1) / 7);

    const iterator = new Date(firstDay);
    iterator.setDate(firstDay.getDate() - firstDayOfWeek);

    days = [];
    for (let i = 0; i < 7 * amountOfRows; i++) {
      days.push(new Date(iterator));
      iterator.setDate(iterator.getDate() + 1);
    }
  })(month, year);
</script>

<style lang="scss">
  @import "../../styles/dimensions.scss";

  div.calendar {
    display: flex;
    flex-direction: column;
    gap: $gap;
    width: 100%;
    height: 100%;
  }

  div.weekdays {
    display: grid;
    grid-template-columns: repeat(7, 1fr);
    gap: $gap;
  }
  div.weekday {
    text-align: center;
  }

  div.days {
    display: grid;
    grid-template-columns: repeat(7, 1fr);
    gap: $gap;
    flex-grow: 1;
  }
</style>

<div class="calendar">
  <div class="weekdays">
    {#each Array(7) as _, weekDay}
      <div class="weekday">
        {getDayName(weekDay)}
      </div>
    {/each}
  </div>
  <div class="days">
    {#each days as day}
      <Day
        day={day.getDate()}
        dayOfWeek={day.getDay()}
        isCurrentMonth={day.getMonth() === month} 
      >
      </Day>
    {/each}
  </div>
</div>
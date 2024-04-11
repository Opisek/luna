<script lang="ts">
    import { getDayName } from "../../lib/common/humanization";
    import Day from "./Day.svelte";

  export let month: number;
  export let year: number;

  let firstDay: Date;
  let days: Date[] = [];

  $: ((month: number, year: number) => {
    firstDay = new Date(year, month, 1);
    let firstDayOfWeek = firstDay.getDay() - 1;
    if (firstDayOfWeek === -1) firstDayOfWeek = 6;
    console.log(firstDayOfWeek);
    firstDay.setDate(firstDay.getDate() - firstDayOfWeek);

    days = [];
    for (let i = 0; i < 7 * 5; i++) {
      days.push(new Date(firstDay));
      firstDay.setDate(firstDay.getDate() + 1);
    }
  })(month, year);
</script>

<style lang="scss">
  div.calendar {
    display: grid;
    grid-template-columns: repeat(7, 1fr);
    grid-template-rows: auto repeat(5, 1fr);
    gap: 1rem;
    width: 100%;
    height: 100%;
  }
</style>

<div class="calendar">
  {#each Array(7) as _, weekDay}
    <div>{getDayName(weekDay)}</div>
  {/each}
  {#each days as day}
    <Day
      day={day.getDate()}
      dayOfWeek={day.getDay()}
      isCurrentMonth={day.getMonth() === month} 
    >
    </Day>
  {/each}
</div>
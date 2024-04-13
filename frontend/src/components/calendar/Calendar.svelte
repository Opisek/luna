<script lang="ts">
  import Day from "./Day.svelte";
  import { getDayName } from "../../lib/common/humanization";

  export let month: number;
  export let year: number;
  export let events: CalendarEventModel[];

  let days: Date[] = [];
  let amountOfRows: number = 0;

  $: ((month: number, year: number) => {
    const firstDay = new Date(year, month, 1);
    const lastDay = new Date(year, month + 1, 0);
    const firstDayOfWeek = (firstDay.getDay() + 6) % 7;

    amountOfRows = Math.ceil((lastDay.getDate() + firstDayOfWeek - 1) / 7);

    const iterator = new Date(firstDay);
    iterator.setDate(firstDay.getDate() - firstDayOfWeek);

    days = [];
    for (let i = 0; i < 7 * amountOfRows; i++) {
      days.push(new Date(iterator));
      iterator.setDate(iterator.getDate() + 1);
    }
  })(month, year);

  let calendarEventMap: Map<string, CalendarEventModel[]> = new Map();

  $: ((events: CalendarEventModel[]) => {
    calendarEventMap = new Map();
    events.forEach(event => {
      for (const iterator = new Date(event.start); iterator < event.end; iterator.setDate(iterator.getDate() + 1)) {
        const key = iterator.toISOString().split("T")[0];
        const arr = calendarEventMap.get(key);
        if (arr) arr.push(event);
        else calendarEventMap.set(key, [ event ]);
      }
    });
  })(events);
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
  <div class="days" style="grid-template-rows: repeat({amountOfRows}, 1fr)">
    {#each days as day}
      <Day
        day={day.getDate()}
        dayOfWeek={day.getDay()}
        isCurrentMonth={day.getMonth() === month} 
        events={calendarEventMap.get(day.toISOString().split("T")[0]) || []}
      >
      </Day>
    {/each}
  </div>
</div>
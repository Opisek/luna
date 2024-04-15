<script lang="ts">
  import Day from "./Day.svelte";
  import { getDayName } from "../../lib/common/humanization";
  import { compareEventsByStartDate } from "../../lib/common/comparators";

  export let month: number;
  export let year: number;
  export let events: CalendarEventModel[];

  let currentlyHoveredEvent: CalendarEventModel | null = null;
  let currentlyClickedEvent: CalendarEventModel | null = null;

  let days: Date[] = [];
  let amountOfRows: number = 0;
  let processedEvents: (CalendarEventModel | null)[][] = [];

  $: ((month: number, year: number, events: CalendarEventModel[]) => {
    // Date calculation
    const firstMonthDay = new Date(year, month, 1);
    const lastMonthDay = new Date(year, month + 1, 0);
    const firstDayOfWeek = (firstMonthDay.getDay() + 6) % 7;

    amountOfRows = Math.ceil((lastMonthDay.getDate() + firstDayOfWeek - 1) / 7);

    const firstViewDay = new Date(firstMonthDay);
    firstViewDay.setDate(firstMonthDay.getDate() - firstDayOfWeek);
    const lastViewDay = new Date(firstMonthDay);
    lastViewDay.setDate(firstMonthDay.getDate() + 7 * amountOfRows - 1);

    // Event pre-processing
    const filteredEvents = events.sort(compareEventsByStartDate).filter(e => e.start >= firstViewDay && e.end < lastViewDay);

    // Fill
    days = [];
    processedEvents = [];

    const dateIterator = new Date(firstViewDay);
    let eventIterator = 0;

    for (let i = 0; i < 7 * amountOfRows; i++) {
      // Copy events from previous day and remove whichever are over
      const dayEvents =
        i == 0
          ? []
          : processedEvents[i - 1]
            .map(
              e => e === null || e.end <= dateIterator
                ? null
                : e
            );
      
      // Fit new events in fitting slots
      let emptyIterator = 0;
      while (eventIterator < filteredEvents.length && filteredEvents[eventIterator].start <= dateIterator) {
        while (emptyIterator < dayEvents.length && dayEvents[emptyIterator] != null) emptyIterator++;
        if (emptyIterator < dayEvents.length) dayEvents[emptyIterator] = filteredEvents[eventIterator];
        else dayEvents.push(filteredEvents[eventIterator]);
        emptyIterator++;
        eventIterator++;
      }

      // Remove unnecessary nulls
      while(dayEvents.length > 0 && dayEvents[dayEvents.length - 1] == null) dayEvents.pop();

      processedEvents.push(dayEvents);
      days.push(new Date(dateIterator));
      dateIterator.setDate(dateIterator.getDate() + 1);
    }
  })(month, year, events);

  function eventClick(event: CalendarEventModel) {
    window.alert(event.title);
  }
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
    {#each days as day, i}
      <Day
        date={day}
        isCurrentMonth={day.getMonth() === month} 
        events={processedEvents[i]}
        isFirstDay={i == 0}
        isLastDay={i == days.length - 1}
        bind:currentlyHoveredEvent={currentlyHoveredEvent}
        bind:currentlyClickedEvent={currentlyClickedEvent}
        clickCallback={eventClick}
      >
      </Day>
    {/each}
  </div>
</div>
<script lang="ts">
  import Day from "./Day.svelte";

  import { compareEventsByStartDate } from "$lib/common/comparators";
  import { getDayName } from "$lib/common/humanization";

  import { getContext, setContext } from "svelte";
  import { writable } from "svelte/store";
  import { isSameDay } from "$lib/common/date";
  import { fade, fly } from "svelte/transition";

  interface Props {
    date: Date;
    view: "month" | "week" | "day";
    events: EventModel[];
  }

  let {
    date,
    view,
    events,
  }: Props = $props();

  let today = new Date();

  let currentlyClickedEvent = writable(null);
  let currentlyHoveredEvent = writable(null);
  setContext("currentlyHoveredEvent", currentlyHoveredEvent);
  setContext("currentlyClickedEvent", currentlyClickedEvent);

  let startDate = $derived(
    (() => {
      switch (view) {
        case "month":
          return new Date(date.getFullYear(), date.getMonth(), 1);
        case "week":
          return new Date(date.getFullYear(), date.getMonth(), date.getDate() - ((date.getDay() + 6) % 7));
        case "day":
          return new Date(date.getFullYear(), date.getMonth(), date.getDate());
      }
    })()
  );
  let endDate = $derived(
    (() => {
      switch (view) {
        case "month":
          return new Date(date.getFullYear(), date.getMonth() + 1, 0);
        case "week":
          return new Date(date.getFullYear(), date.getMonth(), date.getDate() - ((date.getDay() + 6) % 7) + 7);
        case "day":
          return new Date(date.getFullYear(), date.getMonth(), date.getDate() + 1);
      }
    })()
  );

  let [days, amountOfRows, processedEvents] = $derived((() => {
      // Date calculation
      const firstDayOfWeek = (startDate.getDay() + 6) % 7;

      const amountOfColumns = view === "day" ? 1 : 7;
      const amountOfRows = view === "month" ? Math.ceil((endDate.getDate() + firstDayOfWeek) / amountOfColumns) : 1;

      const firstViewDay = new Date(startDate);
      if (view === "month") firstViewDay.setDate(startDate.getDate() - firstDayOfWeek);
      const lastViewDay = new Date(endDate);
      if (view === "month") lastViewDay.setDate(firstViewDay.getDate() + amountOfColumns * amountOfRows - 1);
      const nextViewDay = new Date(lastViewDay);
      nextViewDay.setDate(nextViewDay.getDate() + 1);

      // Event pre-processing
      const filteredEvents = events.filter(e => e.date.end.getTime() >= firstViewDay.getTime() && e.date.start.getTime() < nextViewDay.getTime());
      filteredEvents.sort(compareEventsByStartDate);

      // Fill
      const days: Date[] = [];
      const processedEvents: (EventModel | null)[][] = [];

      const dateIterator = new Date(firstViewDay);
      let eventIterator = 0;

      // Long events from previous view should be added to the current view
      const pastViewEvents = [];
      while (eventIterator < filteredEvents.length && filteredEvents[eventIterator].date.start.getTime() < dateIterator.getTime()) {
        pastViewEvents.push(filteredEvents[eventIterator]);
        eventIterator++;
      }

      for (let i = 0; i < amountOfColumns * amountOfRows; i++) {
        // Copy events from previous day and remove whichever are over
        const dayEvents =
          i == 0
            ? pastViewEvents
            : processedEvents[i - 1]
              .map(
                e => e === null || e.date.end.getTime() <= dateIterator.getTime()
                  ? null
                  : e
              );
              
        
        days.push(new Date(dateIterator));
        dateIterator.setDate(dateIterator.getDate() + 1);
        
        // Fit new events in fitting slots
        let emptyIterator = 0;
        while (
          eventIterator < filteredEvents.length &&
          filteredEvents[eventIterator].date.start.getTime() < dateIterator.getTime() &&
          filteredEvents[eventIterator].date.start.getTime() >= days[days.length - 1].getTime()
        ) {
          while (emptyIterator < dayEvents.length && dayEvents[emptyIterator] != null) emptyIterator++;
          if (emptyIterator < dayEvents.length) dayEvents[emptyIterator] = filteredEvents[eventIterator];
          else dayEvents.push(filteredEvents[eventIterator]);
          emptyIterator++;
          eventIterator++;
        }

        // Remove unnecessary nulls
        while(dayEvents.length > 0 && dayEvents[dayEvents.length - 1] == null) dayEvents.pop();
        processedEvents.push(dayEvents);
      }

      return [days, amountOfRows, processedEvents];
  })());

  let containerHeight: number = $state(0);
  // TODO: figure out how to do this without hard-codeds
  // 9: gap between events
  // 27: height of an event
  let maxEvents: number = $derived(containerHeight === 0 ? 0 : Math.max(Math.floor((containerHeight + 9) / 27), 0));

  /* Show more */
  let showDateModal: ((date: Date, events: (EventModel | null)[]) => any) = getContext("showDateModal");
  function showMore(date: Date, events: (EventModel | null)[]) {
    showDateModal(date, events);
  }
</script>

<style lang="scss">
  @use "../../styles/dimensions.scss";

  div.calendar {
    display: flex;
    flex-direction: column;
    gap: dimensions.$gapSmall;
    width: 100%;
    height: 100%;
    overflow: hidden;
  }

  div.weekdays {
    display: grid;
    gap: dimensions.$gapSmall;
  }
  //div.weekdays-day {
  //  grid-template-columns: repeat(1, 1fr);
  //}
  div.weekday {
    text-align: center;
  }

  div.days {
    display: grid;
    gap: 0;
    flex-grow: 1;
  }
  
  div.columns-month,
  div.columns-week {
    grid-template-columns: repeat(7, 1fr);
  }
  div.columns-day {
    grid-template-columns: repeat(1, 1fr);
  }
</style>

<div class="calendar">
  <div
    class="weekdays"
    class:columns-month={view === "month"}
    class:columns-week={view === "week"}
    class:columns-day={view === "day"}
  >
    {#if view === "month" || view === "week"}
      {#each Array(7) as _, weekDay}
        <div class="weekday">
          {getDayName(weekDay)}
        </div>
      {/each}
    {:else}
      <div class="weekday">
        {getDayName((date.getDay() + 6) % 7)}
      </div>
    {/if}
  </div>
  <div
    class="days"
    style="grid-template-rows: repeat({amountOfRows}, 1fr)"
    class:columns-month={view === "month"}
    class:columns-week={view === "week"}
    class:columns-day={view === "day"}
  >
    {#each days as day, i}
      <Day
        date={day}
        isCurrentMonth={day.getMonth() === date.getMonth()} 
        events={processedEvents[i]}
        isFirstDay={i == 0}
        isToday={isSameDay(day, today)}
        maxEvents={maxEvents}
        bind:containerHeight={containerHeight}
        view={view}
        showMore={showMore}
      >
      </Day>
    {/each}
  </div>
</div>
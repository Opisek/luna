<script lang="ts">
  import Day from "./Day.svelte";
  import EventModal from "../modals/EventModal.svelte";

  import { EmptyEvent, NoOp } from "$lib/client/placeholders";
  import { compareEventsByStartDate } from "$lib/common/comparators";
  import { getDayName } from "$lib/common/humanization";
  import { untrack } from "svelte";

  interface Props {
    month: number;
    year: number;
    events: EventModel[];
  }

  let {
    month,
    year,
    events,
  }: Props = $props();

  let currentlyHoveredEvent: EventModel = $state(EmptyEvent);
  let currentlyClickedEvent: EventModel = $state(EmptyEvent);

  let showModal = $state(NoOp);

  let [days, amountOfRows, processedEvents] = $derived((() => {
      // Date calculation
      const firstMonthDay = new Date(year, month, 1);
      const lastMonthDay = new Date(year, month + 1, 0);
      const firstDayOfWeek = (firstMonthDay.getDay() + 6) % 7;

      const amountOfRows = Math.ceil((lastMonthDay.getDate() + firstDayOfWeek) / 7);

      const firstViewDay = new Date(firstMonthDay);
      firstViewDay.setDate(firstMonthDay.getDate() - firstDayOfWeek);
      const lastViewDay = new Date(firstMonthDay);
      lastViewDay.setDate(firstMonthDay.getDate() + 7 * amountOfRows - 1);

      // Event pre-processing
      const filteredEvents = events.filter(e => e.date.start.getTime() >= firstViewDay.getTime() && e.date.end.getTime() < lastViewDay.getTime());
      filteredEvents.sort(compareEventsByStartDate);

      // Fill
      const days: Date[] = [];
      const processedEvents: (EventModel | null)[][] = [];

      const dateIterator = new Date(firstViewDay);
      let eventIterator = 0;

      for (let i = 0; i < 7 * amountOfRows; i++) {
        // Copy events from previous day and remove whichever are over
        const dayEvents =
          i == 0
            ? []
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

  let clickedEvent: EventModel = $state(EmptyEvent);
  function eventClick(event: EventModel) {
    if (!event) return
    clickedEvent = event;
    setTimeout(() => showModal(), 0);
  }

  let containerHeight: number = $state(0);
  let maxEvents: number = $state(0);
  function calculateMaxEvents(height: number) {
    if (height == 0) {
      maxEvents = 0;
      return;
    }
  }
  $effect(() => {
    calculateMaxEvents(containerHeight);
  });
</script>

<style lang="scss">
  @import "../../styles/dimensions.scss";

  div.calendar {
    display: flex;
    flex-direction: column;
    gap: $gapSmall;
    width: 100%;
    height: 100%;
  }

  div.weekdays {
    display: grid;
    grid-template-columns: repeat(7, 1fr);
    gap: $gapSmall;
  }
  div.weekday {
    text-align: center;
  }

  div.days {
    display: grid;
    grid-template-columns: repeat(7, 1fr);
    gap: 0;
    flex-grow: 1;
    // TODO: figure out proper height
    height: 90%;
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
        maxEvents={maxEvents}
        bind:containerHeight={containerHeight}
      >
      </Day>
    {/each}
  </div>
</div>

{#if clickedEvent.id}
  <EventModal bind:showModal event={clickedEvent}/>
{/if}
<script lang="ts">
  import Day from "./Day.svelte";

  import { compareEventsByStartDate } from "$lib/common/comparators";
  import { getDayName } from "$lib/common/humanization";

  import { getContext, setContext } from "svelte";
  import { writable } from "svelte/store";
  import { getDayIndex, getWeekNumber, getWeekMonth, isSameDay } from "$lib/common/date";
  import { getSettings } from "$lib/client/settings.svelte";
  import { UserSettingKeys } from "../../types/settings";
  import { svelteFlyInHorizontal, svelteFlyOutHorizontal } from "../../lib/client/animations";

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

  const settings = getSettings();

  let today = new Date();

  let currentlyClickedEvent = writable(null);
  let currentlyHoveredEvent = writable(null);
  setContext("currentlyHoveredEvent", currentlyHoveredEvent);
  setContext("currentlyClickedEvent", currentlyClickedEvent);

  /* Animation */
  let viewIteration = $state(0);
  let currentDate = $state(new Date(date));
  let flyDirection = $state("left");
  setContext("flyDirection", () => flyDirection);
  $effect(() => {
    if (date.getTime() === currentDate.getTime()) return;
    flyDirection = currentDate.getTime() <= date.getTime() ? "left" : "right";
    viewIteration++;
    currentDate = new Date(date);
  });

  /* View Calculation */

  let startDate = $derived(
    (() => {
      switch (view) {
        case "month":
          return new Date(currentDate.getFullYear(), currentDate.getMonth(), 1);
        case "week":
          return new Date(currentDate.getFullYear(), currentDate.getMonth(), currentDate.getDate() - getDayIndex(currentDate));
        case "day":
          return new Date(currentDate.getFullYear(), currentDate.getMonth(), currentDate.getDate());
      }
    })()
  );
  let endDate = $derived(
    (() => {
      switch (view) {
        case "month":
          return new Date(currentDate.getFullYear(), currentDate.getMonth() + 1, 0);
        case "week":
          return new Date(currentDate.getFullYear(), currentDate.getMonth(), currentDate.getDate() - getDayIndex(currentDate) + 7);
        case "day":
          return new Date(currentDate.getFullYear(), currentDate.getMonth(), currentDate.getDate() + 1);
      }
    })()
  );

  let [days, amountOfRows, processedEvents] = $derived.by(() => {
    // Date calculation
    const firstDayOfWeek = getDayIndex(startDate);

    const amountOfColumns = view === "day" ? 1 : 7;
    const amountOfRows = 
      view === "month" ?
      settings.userSettings[UserSettingKeys.DynamicCalendarRows] ?
      Math.ceil((endDate.getDate() + firstDayOfWeek) / amountOfColumns)
      : 6
      : 1;

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
  });

  let containerHeight: number = $state(0);
  // TODO: figure out how to do this without hard-coded values
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
  @use "../../styles/colors.scss";

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
    margin: 0 dimensions.$gapSmaller;
  }
  div.weekday {
    text-align: center;
  }
  div.weekdays.padded {
    padding-left: calc(1.7em + dimensions.$gapSmaller);
  }

  div.monthAnimation {
    width: 100%;
    height: 100%;
    overflow: hidden;
    position: relative;
  }
  
  div.days.animate {
    position: absolute;
  }

  div.days {
    width: 100%;
    height: 100%;
    display: grid;
    gap: 0;
    flex-grow: 1;
    padding: 0;
    margin: 0;
  }

  div.columns-month,
  div.columns-week {
    grid-template-columns: repeat(7, 1fr);
  }
  div.columns-day {
    grid-template-columns: repeat(1, 1fr);
  }

  div.weekNumbersWrapper {
    display: flex;
    flex-direction: row;
    gap: dimensions.$gapSmaller;
    flex-grow: 1;
  }

  div.weekNumbers {
    display: grid;
    flex-direction: column;
    width: 1.7em;
  }
  
  div.weekNumber {
    display: flex;
    justify-content: center;
    align-items: center;
    background-color: colors.$backgroundSecondary;
    margin: dimensions.$gapSmaller 0;
    border-radius: dimensions.$borderRadiusSmall;
  }

  div.weekNumber.otherMonth {
    opacity: 0.5;
  }
</style>

<div class="calendar">
  <div
    class="weekdays"
    class:columns-month={view === "month"}
    class:columns-week={view === "week"}
    class:columns-day={view === "day"}
    class:padded={settings.userSettings[UserSettingKeys.DisplayWeekNumbers]}
  >
    {#if view === "month" || view === "week"}
      {#each Array(7) as _, weekDay}
        <div class="weekday">
          {getDayName((weekDay + settings.userSettings[UserSettingKeys.FirstDayOfWeek]) % 7)}
        </div>
      {/each}
    {:else}
      <div class="weekday">
        {getDayName(date.getDay())}
      </div>
    {/if}
  </div>

  {#if settings.userSettings[UserSettingKeys.DisplayWeekNumbers]}
    <div class="weekNumbersWrapper">
      <div
        class="weekNumbers"
      >
        {#each Array(amountOfRows) as _, i}
        {@const weekNumber = getWeekNumber(days[7 * i] || new Date())}
        <div class="weekNumber" class:otherMonth={view !== "day" && getWeekMonth(weekNumber, (days[7 * i] || new Date()).getFullYear()) !== date.getMonth()}>
          {weekNumber}
        </div>
        {/each}
      </div>
      {@render daysGridAnimated()}
    </div>
  {:else}
    {@render daysGridAnimated()}
  {/if}
</div>


{#snippet daysGridAnimated()}
  {#if settings.userSettings[UserSettingKeys.AnimateCalendarSwipe]}
    <div class="monthAnimation">
      {#each [ days ] as currentDays (viewIteration)}
        {@render grid(currentDays, amountOfRows, true)}
      {/each}
    </div>
  {:else}
    {@render grid(days, amountOfRows, false)}
  {/if}
{/snippet}

{#snippet grid(days: Date[], amountOfRows: number, animate: boolean)}
  <div
    class="days"
    class:animate={animate}
    style="grid-template-rows: repeat({amountOfRows}, 1fr)"
    class:columns-month={view === "month"}
    class:columns-week={view === "week"}
    class:columns-day={view === "day"}
    in:svelteFlyInHorizontal={{duration: animate ? 500 : 0}}
    out:svelteFlyOutHorizontal={{duration: animate ? 500 : 0}}
  >
    {#each days as day, i}
      <Day
        date={day}
        isCurrentMonth={day.getMonth() === currentDate.getMonth()} 
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
{/snippet}
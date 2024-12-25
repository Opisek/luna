<script lang="ts">
  import { PlusIcon } from "lucide-svelte";

  import Calendar from "../components/calendar/Calendar.svelte";
  import CalendarEntry from "../components/calendar/CalendarEntry.svelte";
  import EventModal from "../components/modals/EventModal.svelte";
  import Horizontal from "../components/layout/Horizontal.svelte";
  import IconButton from "../components/interactive/IconButton.svelte";
  import MonthSelection from "../components/interactive/MonthSelection.svelte";
  import SelectButtons from "../components/forms/SelectButtons.svelte";
  import SourceEntry from "../components/calendar/SourceEntry.svelte";
  import SourceModal from "../components/modals/SourceModal.svelte";
  import Title from "../components/layout/Title.svelte";

  import { afterNavigate } from "$app/navigation";
  import { browser } from "$app/environment";

  import { NoOp } from "$lib/client/placeholders";
  import { calendars, events, fetchAllEvents, fetchSources, sources } from "$lib/client/repository";
  import { queueNotification } from "$lib/client/notifications";

  import { setContext, untrack } from "svelte";

  /* View */
  let view: "month" | "week" | "day" = $state("month");

  /* Fetched data */
  let localSources: SourceModel[] = $state([]);
  let localCalendars: CalendarModel[] = $state([]);
  let localEvents: EventModel[] = $state([]);

  let sourceCalendars: Map<string, number[]> = $state(new Map());
  let calendarEvents: Map<string, number[]> = new Map();

  /* Fetching logic */
  let loaded: boolean = $state(false);

  const today = new Date();
  let date = $state(today);
  let rangeStart: Date = $derived(new Date(date.getFullYear(), date.getMonth() - 1, 1));
  let rangeEnd: Date = $derived(new Date(date.getFullYear(), date.getMonth() + 1, 0));

  function getRangeFromStorage() {
    const storedDate = browser ? sessionStorage.getItem("selectedDate") : null;
    date = storedDate === null ? today : new Date(storedDate);
  }

  if (browser) {
    getRangeFromStorage();
    loaded = true;
  }

  afterNavigate(() => {
    getRangeFromStorage();
  });

  (async () => {
    if (!browser) return;

    fetchSources().catch(err => {
      queueNotification(
        "failure",
        `Failed to fetch sources: ${err.message}`
      );
    });

    events.subscribe((newEvents) => {
      localEvents = newEvents;

      calendarEvents = new Map();
      localEvents.forEach((event, i) => {
        if (calendarEvents.has(event.calendar)) {
          // @ts-ignore typescript says that this might be undefined despite the check above
          calendarEvents.get(event.calendar).push(i);
        } else {
          calendarEvents.set(event.calendar, [ i ]);
        }
      });
    });

    calendars.subscribe((newCalendars) => {
      localCalendars = newCalendars;

      sourceCalendars = new Map();
      localCalendars.forEach((calendar, i) => {
        if (sourceCalendars.has(calendar.source)) {
          // @ts-ignore typescript says that this might be undefined despite the check above
          sourceCalendars.get(calendar.source).push(i);
        } else {
          sourceCalendars.set(calendar.source, [ i ]);
        }
      });
    });

    sources.subscribe((newSources) => {
      localSources = newSources;
    });
  })();

  /* Month selection logic */
  $effect(() => {
    ((date: Date, loaded: boolean) => {
      untrack(() => {
        if (!browser || !loaded) return;

        sessionStorage.setItem("selectedDate", date.toString());

        fetchAllEvents(rangeStart, rangeEnd); // TODO: what actually has to be refetched will be moved to an event manager singleton
        
        //if (lastDayPreviousMonth < rangeStart) {
        //  const lastStart = rangeStart;
        //  rangeStart = new Date(year, month - 2, 1);
        //  fetchAllEvents(rangeStart, lastStart);
        //}

        //if (firstDayNextMonth > rangeEnd) {
        //  const lastEnd = rangeEnd;
        //  rangeEnd = new Date(year, month + 3, 0);
        //  fetchAllEvents(lastEnd, rangeEnd);
        //}
      });
    })(date, loaded);
  });

  /* Single instance modal logic */
  let showNewSourceModal: () => any = $state(NoOp);

  let showSourceModal: (source: SourceModel) => any = $state(NoOp);
  const showSourceModalInternal = (source: SourceModel) => { return showSourceModal(source); };
  setContext("showSourceModal", showSourceModalInternal);

  //let showNewCalendarModal: () => any = $state(NoOp);
  //let showCalendarModal: () => any = $state(NoOp);

  let showNewEventModal: (date: Date) => any = $state(NoOp);
  const showNewEventModalInternal = (date: Date) => { return showNewEventModal(date); };
  setContext("showNewEventModal", showNewEventModalInternal);

  let showEventModal: (event: EventModel) => any = $state(NoOp);
  const showEventModalInternal = (event: EventModel) => { return showEventModal(event); };
  setContext("showEventModal", showEventModalInternal);
</script>

<style lang="scss">
  @import "../styles/dimensions.scss";

  main {
    width: 100%;
    height: 100%;
    display: flex;
    flex-direction: column;
    gap: $gap;
  }

  div.wrapper {
    display: flex;
    flex-direction: row;
    gap: $gap;
    padding: $gap;
    height: 100%;
    width: 100%;
  }
  
  aside {
    display: flex;
    flex-direction: column;
    gap: $gap;
    min-width: 10em;
    width: 20vw;
    max-width: 20em;
    overflow: hidden;
  }

  div.sources {
    flex-grow: 1;
    display: flex;
    flex-direction: column;
    gap: $gap;
  }

  div.toprow {
    display: flex;
    flex-direction: row;
    gap: $gapSmall;
    justify-content: space-between;
    margin: 0 $gapSmaller;
    align-items: center;
  }
</style>

<SourceModal bind:showCreateModal={showNewSourceModal} bind:showModal={showSourceModal}/>
<!--<CalendarModal bind:showCreateModal={showNewCalendarModal} bind:showModal={showCalendarModal}}/>-->
<EventModal bind:showCreateModal={showNewEventModal} bind:showModal={showEventModal}/>

<div class="wrapper">
  <aside>
    <Title>Luna</Title>

    <!-- SmallCalendar put here only for testing purposes but might consider leaving it if it can serve some useful purpose -->
    <!--<SmallCalendar year={selectedYear} month={selectedMonth}/>-->

    <div class="sources">
      {@render sourceEntries(localSources)}
    </div>
    <Horizontal position="center">
      <IconButton click={showNewSourceModal}>
        <PlusIcon/>
      </IconButton>
    </Horizontal>
  </aside>
  <main>
    <div class="toprow">
      <MonthSelection bind:date granularity={view} />
        <SelectButtons
          name="layout"
          compact={true}
          bind:value={view}
          options={[
            { value: "day", name: "Day"},
            { value: "week", name: "Week"},
            { value: "month", name: "Month"},
          ]}
        />
    </div>
      <Calendar
        date={date}
        view={view}
        events={localEvents}
      />
  </main>
</div>

{#snippet sourceEntries(sources: SourceModel[])}
  {#each sources as source, i}
    <SourceEntry bind:source={localSources[i]}/>
    {#if (!source.collapsed)}
      {@render calendarEntries(sourceCalendars.get(source.id) || [])} 
    {/if}
  {/each}
{/snippet}

{#snippet calendarEntries(calendarIndices: number[])}
  {#each calendarIndices as i}
    <CalendarEntry bind:calendar={localCalendars[i]}/>
  {/each}
{/snippet}
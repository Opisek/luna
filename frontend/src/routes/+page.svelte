<script lang="ts">
  import { PlusIcon, RefreshCw } from "lucide-svelte";

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

  import { afterNavigate, beforeNavigate } from "$app/navigation";
  import { browser } from "$app/environment";

  import { NoOp } from "$lib/client/placeholders";
  import { calendars, events, getAllEvents, getSources, invalidateCache, loadingData, sources } from "$lib/client/repository";
  import { queueNotification } from "$lib/client/notifications";

  import { setContext, untrack } from "svelte";

  /* Constants */
  let autoRefreshInterval = 1000 * 60; // 1 minute

  /* View */
  let view: "month" | "week" | "day" = $state("month");

  /* Fetched data */
  let localSources: SourceModel[] = $state([]);
  let localCalendars: CalendarModel[] = $state([]);
  let localEvents: EventModel[] = $state([]);

  let sourceCalendars: Map<string, number[]> = $state(new Map());
  let calendarEvents: Map<string, number[]> = new Map();

  /* Fetching logic */
  let pageLoaded: boolean = $state(false);

  let isLoading: boolean = $state(false);
  let loaderAnimation = $state(false);
  loadingData.subscribe((loadingData) => {
    isLoading = loadingData;
    if (isLoading) loaderAnimation = true;
  });

  const today = new Date();
  let date = $state(today);

  function getRangeFromStorage() {
    if (pageLoaded) return;
    const storedDate = browser ? sessionStorage.getItem("selectedDate") : null;
    date = storedDate === null ? today : new Date(storedDate);
    pageLoaded = true;
  }

  afterNavigate(() => {
    getRangeFromStorage();
  });

  beforeNavigate(() => {
    pageLoaded = false;
  });

  (async () => {
    if (!browser) return;

    getSources().catch(NoOp);

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

  let spooledRefresh: ReturnType<typeof setTimeout>;
  function refresh(date: Date, force = false) {
    sessionStorage.setItem("selectedDate", date.toString());

    const rangeStart = new Date(date);
    rangeStart.setHours(0, 0, 0, 0);
    const rangeEnd = new Date(date);
    rangeEnd.setHours(23, 59, 59, 999);

    switch (view) {
      case "month":
        rangeStart.setDate(1);
        rangeEnd.setMonth(rangeEnd.getMonth() + 1);
        rangeEnd.setDate(0);
        break;
      case "week":
        rangeStart.setDate(date.getDate() - ((date.getDay() + 6) % 7));
        rangeEnd.setDate(rangeStart.getDate() + 7);
        break;
      case "day":
      default:
    }

    getAllEvents(rangeStart, rangeEnd, force).catch((err) => {
      queueNotification("failure", `Failed to fetch events: ${err.message}`);
    });

    clearTimeout(spooledRefresh);
    spooledRefresh = setTimeout(() => {
      refresh(date);
    }, autoRefreshInterval);
  }

  function forceRefresh() {
    invalidateCache();
    refresh(date, true);
  }

  $effect(() => {
    ((date: Date, loaded: boolean) => {
      untrack(() => {
        if (!browser) return;
        if (!loaded) {
          getRangeFromStorage();
          return;
        }
        refresh(date);
      });
    })(date, pageLoaded);
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
  @use "../styles/animations.scss";
  @use "../styles/dimensions.scss";

  main {
    width: 100%;
    height: 100%;
    display: flex;
    flex-direction: column;
    gap: dimensions.$gap;
  }

  div.wrapper {
    display: flex;
    flex-direction: row;
    gap: dimensions.$gap;
    padding: dimensions.$gap;
    height: 100%;
    width: 100%;
  }
  
  aside {
    display: flex;
    flex-direction: column;
    gap: dimensions.$gap;
    min-width: 10em;
    width: 20vw;
    max-width: 20em;
    overflow: hidden;
  }

  div.sources {
    flex-grow: 1;
    display: flex;
    flex-direction: column;
    gap: dimensions.$gap;
  }

  div.toprow {
    display: flex;
    flex-direction: row;
    gap: dimensions.$gapSmall;
    justify-content: space-between;
    margin: 0 dimensions.$gapSmaller;
    align-items: center;
  }

  span.refreshButtonWrapper {
    display: flex;
    align-items: center;
    justify-content: center;
  }

  span.spin {
    animation: spin animations.$animationSpeedSlow animations.$cubic infinite forwards;
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
      <Horizontal position="right" width="auto">
        <IconButton click={forceRefresh}>
          <span class="refreshButtonWrapper" class:spin={loaderAnimation} onanimationiteration={() => { if (!isLoading) loaderAnimation = false; }}>
            <RefreshCw size={20}/>
          </span>
        </IconButton>
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
      </Horizontal>
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
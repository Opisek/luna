<script lang="ts">
  import { Github } from "svelte-simples"
  import { PlusIcon, RefreshCw, Settings, WifiOff } from "lucide-svelte";
  import { setContext, untrack } from "svelte";

  import Calendar from "../components/calendar/Calendar.svelte";
  import CalendarEntry from "../components/calendar/CalendarEntry.svelte";
  import CalendarModal from "../components/modals/CalendarModal.svelte";
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

  import SmallCalendar from "../components/interactive/SmallCalendar.svelte";
  import { NoOp } from "$lib/client/placeholders";
  import { getMetadata } from "$lib/client/metadata";
  import { getRepository } from "$lib/client/repository";
  import { queueNotification } from "$lib/client/notifications";
  import { getConnectivity, Reachability } from "$lib/client/connectivity";
  import Button from "../components/interactive/Button.svelte";

  /* Reachability */
  let reachability: Reachability = $state(Reachability.Database);

  /* Constants */
  let autoRefreshInterval = 1000 * 60; // 1 minute

  /* View */
  let view: "month" | "week" | "day" = $state("month");

  /* Fetched data */
  let localSources: SourceModel[] = $state([]);
  let isCollapsed: boolean[] = $state([]);
  let localCalendars: CalendarModel[] = $state([]);
  let localEvents: EventModel[] = $state([]);

  let sourceCalendars: Map<string, number[]> = $state(new Map());
  let calendarEvents: Map<string, number[]> = new Map();

  /* Fetching logic */
  let pageLoaded: boolean = $state(false);

  let isLoading: boolean = $state(false);
  let loaderAnimation = $state(false);
  getMetadata().loadingData.subscribe((loadingData) => {
    isLoading = loadingData;
    if (isLoading) loaderAnimation = true;
  });

  let today = $state(new Date());
  let date = $state(new Date());
  let todayInRange = $state(true);

  function seeToday() {
    today = new Date();
    date = new Date(today);
    todayInRange = true;
  }

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

  getRepository().getSources().catch(NoOp);

  getRepository().events.subscribe((newEvents) => {
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

  getRepository().calendars.subscribe((newCalendars) => {
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

  getRepository().sources.subscribe((newSources) => {
    isCollapsed = newSources.map((source) => getMetadata().collapsedSources.has(source.id));
    localSources = newSources;
  });

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

    todayInRange = rangeStart.getTime() <= today.getTime() && today.getTime() <= rangeEnd.getTime();

    getRepository().getAllEvents(rangeStart, rangeEnd, force).catch((err) => {
      queueNotification("failure", `Failed to fetch events: ${err.message}`);
    });

    getConnectivity().check().then((res) => reachability = res);

    clearTimeout(spooledRefresh);
    spooledRefresh = setTimeout(() => {
      refresh(date);
    }, autoRefreshInterval);
  }

  function forceRefresh() {
    getRepository().invalidateCache();
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

  getMetadata().collapsedSources.subscribe((collapsed) => {
    isCollapsed = localSources.map((source) => collapsed.has(source.id));
  });

  /* Single instance modal logic */
  let showNewSourceModal: () => any = $state(NoOp);

  let showSourceModal: (source: SourceModel) => any = $state(NoOp);
  const showSourceModalInternal = (source: SourceModel) => { return showSourceModal(source); };
  setContext("showSourceModal", showSourceModalInternal);

  let showNewCalendarModal: () => any = $state(NoOp);

  let showCalendarModal: (calendar: CalendarModel) => any = $state(NoOp);
  const showCalendarModalInternal = (calendar: CalendarModel) => { return showCalendarModal(calendar); };
  setContext("showCalendarModal", showCalendarModalInternal);

  let showNewEventModal: (date: Date) => any = $state(NoOp);
  const showNewEventModalInternal = (date: Date) => { return showNewEventModal(date); };
  setContext("showNewEventModal", showNewEventModalInternal);

  let showEventModal: (event: EventModel) => any = $state(NoOp);
  const showEventModalInternal = (event: EventModel) => { return showEventModal(event); };
  setContext("showEventModal", showEventModalInternal);
</script>

<style lang="scss">
  @use "../styles/animations.scss";
  @use "../styles/colors.scss";
  @use "../styles/dimensions.scss";

  :global(body) {
    display: flex;
    flex-direction: row;
    //display: grid;
    //grid-template-columns: auto 1fr;
    ////grid-template-rows: 1fr auto;
    ////grid-template-areas:
    ////  "aside main"
    ////  "aside footer";
    //grid-template-rows: auto;
    //grid-template-areas: "aside main";
  }

  main {
    width: 100%;
    height: 100%;
    display: flex;
    flex-direction: column;
    gap: dimensions.$gapLarge;
    grid-area: main;
  }
  
  aside {
    display: flex;
    flex-direction: column;
    gap: dimensions.$gapLarge;
    min-width: 10em;
    width: 20vw;
    max-width: 20em;
    grid-area: aside;
  }

  div.sources {
    flex-grow: 1;
    display: flex;
    flex-direction: column;
    gap: dimensions.$gapLarge;
    overflow: auto;
    margin: -(dimensions.$gapSmall);
    padding: dimensions.$gapSmall;
  }

  div.toprow {
    display: flex;
    flex-direction: row;
    gap: dimensions.$gapSmall;
    justify-content: space-between;
    margin: 0 dimensions.$gapSmaller;
    align-items: center;
  }

  span.reachability {
    color: colors.$backgroundFailure;
    align-items: center;
    display: flex;
    flex-direction: row;
    justify-content: center;
    gap: dimensions.$gapSmall;
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
<CalendarModal bind:showCreateModal={showNewCalendarModal} bind:showModal={showCalendarModal}/>
<EventModal bind:showCreateModal={showNewEventModal} bind:showModal={showEventModal}/>

<aside>
  <Title>Luna</Title>

  <SmallCalendar date={date} smaller={true} onDayClick={(newDate) => date=newDate}></SmallCalendar>

  <div class="sources">
    {@render sourceEntries(localSources)}
  </div>

  <Horizontal position="center">
    <IconButton href="/settings">
      <Settings/>
    </IconButton>
    <IconButton click={showNewSourceModal}>
      <PlusIcon/>
    </IconButton>
    <IconButton href="https://github.com/Opisek/luna">
      <Github/>
    </IconButton>
  </Horizontal>
</aside>

<main>
  <div class="toprow">
    <MonthSelection bind:date granularity={view} />
    <Horizontal position="justify" width="auto">
      {#if reachability != Reachability.Database}
        <span class="reachability">
          {#if reachability == Reachability.Backend}
            The database cannot be reached
          {:else if reachability == Reachability.Frontend}
            The backend server cannot be reached
          {:else if reachability == Reachability.None}
            The frontend server cannot be reached
          {:else}
            Unknown network error
          {/if}
          <WifiOff size={20}/>
        </span>
      {/if}
      
      <IconButton click={forceRefresh}>
        <span class="refreshButtonWrapper" class:spin={loaderAnimation} onanimationiteration={() => { if (!isLoading) loaderAnimation = false; }}>
          <RefreshCw size={20}/>
        </span>
      </IconButton>

      {#if !todayInRange}
        <Button onClick={seeToday} compact={true}>
          Today
        </Button>
      {/if}

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

{#snippet sourceEntries(sources: SourceModel[])}
  {#each sources as source, i}
    <SourceEntry bind:source={localSources[i]}/>
    {#if !isCollapsed[i]}
      {@render calendarEntries(sourceCalendars.get(source.id) || [])}
    {/if}
  {/each}
{/snippet}

{#snippet calendarEntries(calendarIndices: number[])}
  {#each calendarIndices as i}
    <CalendarEntry bind:calendar={localCalendars[i]}/>
  {/each}
{/snippet}
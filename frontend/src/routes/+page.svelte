<script lang="ts">
  import { Github, Jsonwebtokens } from "svelte-simples"
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
  import { getMetadata } from "$lib/client/metadata.svelte";
  import { getRepository } from "$lib/client/repository.svelte";
  import { queueNotification } from "$lib/client/notifications";
  import { getConnectivity, Reachability } from "$lib/client/connectivity";
  import Button from "../components/interactive/Button.svelte";
  import DayViewModal from "../components/modals/DayViewModal.svelte";
  import { getDayIndex, isInRange } from "../lib/common/date";
  import { compareEventsByStartDate } from "../lib/common/comparators";
  import SourceWizardModal from "../components/modals/SourceWizardModal.svelte";
  import SettingsModal from "../components/modals/SettingsModal.svelte";
  import { getSettings } from "$lib/client/settings.svelte";
  import { UserSettingKeys } from "../types/settings";
  import ThemeToggle from "../components/interactive/ThemeToggle.svelte";
  import { ColorKeys } from "../types/colors";

  /* Singletons */
  const settings = getSettings();
  const metadata = getMetadata();
  const repository = getRepository();

  /* Reachability */
  let reachability: Reachability = $state(Reachability.Database);

  /* Constants */
  let autoRefreshInterval = 1000 * 60; // 1 minute

  /* View */
  let view: "month" | "week" | "day" = $state("month");

  /* View logic */
  let today = $state(new Date());
  let date = $state(new Date());

  // Svelte bug... These don't work like they're supposed to.
  //let rangeStart = $derived.by(() => {
  //  console.log("deriving range start");
  //  const tmp = new Date(date);
  //  tmp.setHours(0, 0, 0, 0);
  //  switch (view) {
  //    case "month":
  //      tmp.setDate(1);
  //      break;
  //    case "week":
  //      tmp.setDate(date.getDate() - ((date.getDay() + 6) % 7));
  //      break;
  //    case "day":
  //    default:
  //  }
  //  console.log("range start changed", date, tmp)
  //  return tmp;
  //});
  //let rangeEnd = $derived.by(() => {
  //  console.log("deriving range end");
  //  const tmp = new Date(date);
  //  tmp.setHours(23, 59, 59, 999);
  //  switch (view) {
  //    case "month":
  //      tmp.setMonth(tmp.getMonth() + 1);
  //      tmp.setDate(0);
  //      break;
  //    case "week":
  //      tmp.setDate(rangeStart.getDate() + 7);
  //      break;
  //    case "day":
  //    default:
  //  }
  //  console.log("range end changed", date, tmp)
  //  return tmp;
  //});
  //let todayInRange = $derived(isInRange(today, rangeStart, rangeEnd));

  function getVisibleRange(date: Date, view: "month" | "week" | "day"): { start: Date, end: Date } {
    const rangeStart = new Date(date);
    const rangeEnd = new Date(date);
    rangeStart.setHours(0, 0, 0, 0);
    rangeEnd.setHours(23, 59, 59, 999);
    switch (view) {
      case "month":
        rangeStart.setDate(1);
        rangeStart.setDate(rangeStart.getDate() - getDayIndex(date));
        rangeEnd.setMonth(rangeEnd.getMonth() + 1);
        rangeEnd.setDate(0);
        rangeEnd.setDate(rangeEnd.getDate() + 7 - getDayIndex(date));
        break;
      case "week":
        rangeStart.setDate(date.getDate() - getDayIndex(date));
        rangeEnd.setDate(rangeStart.getDate() + 7);
        break;
      case "day":
      default:
    }
    return { start: rangeStart, end: rangeEnd };
  }
  let todayInRange = $derived.by(() => {
    const range = getVisibleRange(date, view);
    return isInRange(today, range.start, range.end)
  });

  function seeToday() {
    today = new Date();
    date = new Date(today);
  }

  function getRangeFromStorage() {
    if (pageLoaded) return;

    const storedDate = browser ? sessionStorage.getItem("selectedDate") : null;
    date = storedDate === null ? today : new Date(storedDate);

    const storedView = browser ? sessionStorage.getItem("selectedView") : "month";
    //@ts-ignore
    view = storedView && [ "month", "week", "day" ].includes(storedView) ? storedView : "month";

    pageLoaded = true;
  }

  function smallCalendarClick(clickedDate: Date) {
    const range = getVisibleRange(date, view);
    if (isInRange(clickedDate, range.start, range.end) && clickedDate.getMonth() === date.getMonth()) {
      showDateModal(clickedDate, repository.events
        .filter((event) => event.date.start.getTime() <= clickedDate.getTime() + 24 * 60 * 60 * 1000 && event.date.end.getTime() >= clickedDate.getTime())
        .sort(compareEventsByStartDate)
      );
    } else {
      date = clickedDate;
    }
  }

  /* Fetching logic */
  let pageLoaded: boolean = $state(false);
  let isLoading: boolean = $derived(getMetadata().loadingData);
  let loaderAnimation = $state(false);
  $effect(() => {
    if (isLoading) loaderAnimation = true;
  })

  afterNavigate(() => {
    getConnectivity().check().then((res) => reachability = res);
    getRangeFromStorage();
    refresh();
  });

  beforeNavigate((args) => {
    if (args.to === null) return;
    pageLoaded = false;
    clearTimeout(spooledRefresh);
    spooledRefresh = undefined; 
  });

  getRepository().getSources().catch(NoOp);

  let spooledRefresh: (ReturnType<typeof setTimeout> | undefined) = $state(undefined);
  function refresh(force = false) {
    sessionStorage.setItem("selectedDate", date.toString());
    sessionStorage.setItem("selectedView", view);

    const range = getVisibleRange(date, view);

    getRepository().getAllEvents(range.start, range.end, force).catch((err) => {
      queueNotification(ColorKeys.Danger, `Failed to fetch events: ${err.message}`);
    });

    if (force) getConnectivity().check().then((res) => reachability = res);

    clearTimeout(spooledRefresh);
    spooledRefresh = setTimeout(() => {
      refresh();
    }, autoRefreshInterval);
  }

  function forceRefresh() {
    getRepository().invalidateCache();
    refresh(true);
  }

  $effect(() => {
    ((date: Date, view: "month" | "week" | "day", loaded: boolean) => {
      untrack(() => {
        if (!browser) return;
        if (!loaded) {
          getRangeFromStorage();
          return;
        }
        refresh();
      });
    })(date, view, pageLoaded);
  });

  /* Single instance modal logic */
  let showSourceWizardModal: () => any = $state(NoOp);

  let showNewSourceModal: () => any = $state(NoOp);
  const showNewSourceModalInternal = () => { return showNewSourceModal(); };
  setContext("showNewSourceModal", showNewSourceModalInternal);

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

  let showDateModal: (date: Date, events: (EventModel | null)[]) => any = $state(NoOp);
  const showDateModalInternal = (date: Date, events: (EventModel | null)[]) => { return showDateModal(date, events); };
  setContext("showDateModal", showDateModalInternal);

  let showSettingsModal: () => any = $state(NoOp);
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

<SourceWizardModal bind:showModal={showSourceWizardModal}/>
<SourceModal bind:showCreateModal={showNewSourceModal} bind:showModal={showSourceModal}/>
<CalendarModal bind:showCreateModal={showNewCalendarModal} bind:showModal={showCalendarModal}/>
<EventModal bind:showCreateModal={showNewEventModal} bind:showModal={showEventModal}/>
<DayViewModal bind:showModal={showDateModal}/>
<SettingsModal bind:showModal={showSettingsModal}/>

<aside>
  <Title>Luna</Title>

  {#if settings.userSettings[UserSettingKeys.DisplaySmallCalendar]}
    <SmallCalendar date={date} smaller={true} onDayClick={(clickedDate) => smallCalendarClick(clickedDate)}></SmallCalendar>
  {/if}

  <div class="sources">
    {@render sourceEntries(repository.sources)}
  </div>

  <Horizontal position="center">
    <IconButton click={showSettingsModal}>
      <Settings/>
    </IconButton>
    <IconButton click={showSourceWizardModal}>
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
            The database cannot be reached.
          {:else if reachability == Reachability.Frontend}
            The backend server cannot be reached.
          {:else if reachability == Reachability.None}
            The frontend server cannot be reached.
          {:else if reachability == Reachability.Incompatible}
            The frontend server and the backend server are not compatible.
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

      {#if !settings.userSettings[UserSettingKeys.ThemeSynchronize]}
        <ThemeToggle/>
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
      events={repository.events}
    />
</main>

{#snippet sourceEntries(sources: SourceModel[])}
  {#each sources as source, i}
    <SourceEntry bind:source={repository.sources[i]}/>
    {#if !metadata.collapsedSources.has(repository.sources[i].id)}
      {@render calendarEntries(repository.calendars.filter(cal => cal.source === source.id) || [])}
    {/if}
  {/each}
{/snippet}

{#snippet calendarEntries(calendars: CalendarModel[])}
  {#each calendars as cal}
    {@const index = repository.calendars.findIndex((calendar) => calendar.id === cal.id)}
    <CalendarEntry bind:calendar={repository.calendars[index]}/>
  {/each}
{/snippet}
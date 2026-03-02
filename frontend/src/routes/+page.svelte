<script lang="ts">
  import { Github } from "svelte-simples"
  import { Copyleft, PlusIcon, RefreshCw, Settings, WifiOff } from "lucide-svelte";
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
  import { getMetadata } from "$lib/client/data/metadata.svelte";
  import { getRepository } from "$lib/client/data/repository.svelte";
  import { queueNotification } from "$lib/client/notifications";
  import { getConnectivity, Reachability } from "$lib/client/data/connectivity.svelte";
  import Button from "../components/interactive/Button.svelte";
  import DayViewModal from "../components/modals/DayViewModal.svelte";
  import { getDayIndex, isInRange } from "../lib/common/date";
  import { compareEventsByStartDate } from "../lib/common/comparators";
  import SourceWizardModal from "../components/modals/SourceWizardModal.svelte";
  import SettingsModal from "../components/modals/SettingsModal.svelte";
  import { getSettings } from "$lib/client/data/settings.svelte";
  import { UserSettingKeys } from "../types/settings";
  import ThemeToggle from "../components/interactive/ThemeToggle.svelte";
  import { ColorKeys } from "../types/colors";
  import { page } from "$app/state";
  import CreditsModal from "../components/modals/CreditsModal.svelte";

  /* Singletons */
  const settings = getSettings();
  const metadata = getMetadata();
  const repository = getRepository();
  const connectivity = getConnectivity();

  /* Constants */
  let autoRefreshInterval = 1000 * 60; // 1 minute

  /* View logic */
  let view: "month" | "week" | "day" = $derived.by(() => {
    const stored = page.url.searchParams.get("view");
    if (!stored || !["month", "week", "day"].includes(stored)) return "month"
    return stored as "month" | "week" | "day";
  });

  let today = $state(new Date());
  let date = $derived.by(() => {
    const stored = page.url.searchParams.get("date");
    if (!stored) return new Date(today);
    const parsed = new Date(stored);
    return parsed;
  });

  $effect(() => {
    const url = new URL(window.location.toString());

    url.searchParams.set("view", view);
    url.searchParams.set("date", date.toISOString().split("T")[0]);

    history.replaceState(history.state, '', url);
  })

  function getVisibleRange(date: Date, view: "month" | "week" | "day"): { start: Date, end: Date } {
    const rangeStart = new Date(date);
    const rangeEnd = new Date(date);
    rangeStart.setHours(0, 0, 0, 0);
    rangeEnd.setHours(23, 59, 59, 999);
    switch (view) {
      case "month":
        rangeStart.setDate(1);
        rangeStart.setDate(rangeStart.getDate() - getDayIndex(rangeStart));
        rangeEnd.setMonth(rangeEnd.getMonth() + 1);
        rangeEnd.setDate(0);
        rangeEnd.setDate(rangeEnd.getDate() + 6 - getDayIndex(rangeEnd));
        break;
      case "week":
        rangeStart.setDate(date.getDate() - getDayIndex(rangeStart));
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
  let isLoading: boolean = $derived(getMetadata().loadingData);
  let loaderAnimation = $state(false);
  $effect(() => {
    if (isLoading) loaderAnimation = true;
  })

  afterNavigate(() => {
    refresh();
  });

  beforeNavigate((args) => {
    if (args.to === null) return;
    clearTimeout(spooledRefresh);
    spooledRefresh = undefined; 
  });

  getRepository().getSources().catch(NoOp);

  let spooledRefresh: (ReturnType<typeof setTimeout> | undefined) = $state(undefined);
  function refresh(force = false) {
    const range = getVisibleRange(date, view);

    getRepository().getAllEvents(range.start, range.end, force).catch((err) => {
      queueNotification(ColorKeys.Danger, `Failed to fetch events: ${err.message}`);
    });

    connectivity.check();

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
    ((date: Date, view: "month" | "week" | "day") => {
      untrack(() => {
        if (!browser) return;
        refresh();
      });
    })(date, view);
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

  let showCreditsModal: () => any = $state(NoOp);
</script>

<style lang="scss">
  @use "../styles/animations.scss";
  @use "../styles/colors.scss";
  @use "../styles/dimensions.scss";
  @use "../styles/text.scss";

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

  span.copyright {
    color: color-mix(in srgb, colors.$foregroundPrimary 50%, transparent);
    font-size: text.$fontSizeSmall;
    text-align: center;
    margin-top: -(dimensions.$gapSmall);
  }
</style>

<SourceWizardModal bind:showModal={showSourceWizardModal}/>
<SourceModal bind:showCreateModal={showNewSourceModal} bind:showModal={showSourceModal}/>
<CalendarModal bind:showCreateModal={showNewCalendarModal} bind:showModal={showCalendarModal}/>
<EventModal bind:showCreateModal={showNewEventModal} bind:showModal={showEventModal}/>
<DayViewModal bind:showModal={showDateModal}/>
<SettingsModal bind:showModal={showSettingsModal}/>
<CreditsModal bind:showModal={showCreditsModal}/>

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
    <IconButton click={showCreditsModal}>
      <Copyleft/>
    </IconButton>
  </Horizontal>

  <span class="copyright">
    Copyright © 2025 Kacper Darowski "Opisek" 
    Licensed under TBD 
  </span>
</aside>

<main>
  <div class="toprow">
    <MonthSelection bind:date granularity={view} />
    <Horizontal position="justify" width="auto">
      {#if connectivity.reachable != Reachability.Database}
        <span class="reachability">
          {#if connectivity.reachable == Reachability.Backend}
            The database cannot be reached.
          {:else if connectivity.reachable == Reachability.Frontend}
            The backend server cannot be reached.
          {:else if connectivity.reachable == Reachability.None}
            The frontend server cannot be reached.
          {:else if connectivity.reachable == Reachability.Incompatible}
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
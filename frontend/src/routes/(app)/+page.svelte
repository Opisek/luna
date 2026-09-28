<script lang="ts">
  import { Copyleft, PlusIcon, RefreshCw, Settings, WifiOff } from "lucide-svelte";
  import { setContext, untrack } from "svelte";

  import Calendar from "$lib/components/calendar/Calendar.svelte";
  import CalendarEntry from "$lib/components/calendar/CalendarEntry.svelte";
  import CalendarModal from "$lib/components/modals/CalendarModal.svelte";
  import EventModal from "$lib/components/modals/EventModal.svelte";
  import Horizontal from "$lib/components/layout/Horizontal.svelte";
  import IconButton from "$lib/components/interactive/IconButton.svelte";
  import MonthSelection from "$lib/components/interactive/MonthSelection.svelte";
  import SelectButtons from "$lib/components/forms/SelectButtons.svelte";
  import SourceEntry from "$lib/components/calendar/SourceEntry.svelte";
  import SourceModal from "$lib/components/modals/SourceModal.svelte";
  import Title from "$lib/components/layout/Title.svelte";

  import { afterNavigate, beforeNavigate, replaceState } from "$app/navigation";
  import { browser } from "$app/environment";

  import SmallCalendar from "$lib/components/interactive/SmallCalendar.svelte";
  import { AsyncNoOp, NoOp } from "$lib/scripts/client/placeholders";
  import { getMetadata } from "$lib/scripts/client/data/metadata.svelte";
  import { getRepository } from "$lib/scripts/client/data/repository.svelte";
  import { queueNotification } from "$lib/scripts/client/notifications";
  import { getConnectivity, Reachability } from "$lib/scripts/client/data/connectivity.svelte";
  import Button from "$lib/components/interactive/Button.svelte";
  import DayViewModal from "$lib/components/modals/DayViewModal.svelte";
  import { getDayIndex, isInRange } from "$lib/scripts/common/date";
  import { compareEventsByStartDate } from "$lib/scripts/common/comparators";
  import SourceWizardModal from "$lib/components/modals/SourceWizardModal.svelte";
  import SettingsModal from "$lib/components/modals/SettingsModal.svelte";
  import { getSettings } from "$lib/scripts/client/data/settings.svelte";
  import { UserSettingKeys } from "$lib/types/settings";
  import ThemeToggle from "$lib/components/interactive/ThemeToggle.svelte";
  import { ColorKeys } from "$lib/types/colors";
  import { page } from "$app/state";
  import CreditsModal from "$lib/components/modals/CreditsModal.svelte";
  import CreatePopup from "$lib/components/popups/CreatePopup.svelte";
  import { locale, t } from "@sveltia/i18n";
  import AffectedRecurrencesModal from "$lib/components/modals/AffectedRecurrencesModal.svelte";

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

    replaceState(url, page.state);
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
  let showSourceWizardModalInternal: () => Promise<SourceModel> = $state(Promise.reject);
  const showSourceWizardModal = () => { return showSourceWizardModalInternal(); };

  let showNewSourceModalInternal: () => any = $state(NoOp);
  const showNewSourceModal = () => { return showNewSourceModalInternal(); };
  setContext("showNewSourceModal", showNewSourceModal);

  let showSourceModalInternal: (initial?: SourceModel, anchor?: HTMLElement) => Promise<SourceModel> = $state(Promise.reject);
  const showSourceModal = (source?: SourceModel, anchor?: HTMLElement) => { return showSourceModalInternal(source, anchor); };
  setContext("showSourceModal", showSourceModal);

  let showCalendarModalInternal: (initial?: CalendarModel, anchor?: HTMLElement) => any = $state(NoOp);
  const showCalendarModal = (calendar?: CalendarModel, anchor?: HTMLElement) => { return showCalendarModalInternal(calendar, anchor); };
  setContext("showCalendarModal", showCalendarModal);

  let showEventModalInternal: (initial?: EventModel, date?: Date, anchor?: HTMLElement) => any = $state(NoOp);
  const showEventModal = (initial?: EventModel, date?: Date, anchor?: HTMLElement) => { return showEventModalInternal(initial, date, anchor); };
  setContext("showEventModal", showEventModal);

  let selectAffectedRecurrencesModalInternal: (edit: boolean) => Promise<"all" | "thisandfuture" | "this"> = $state(Promise.reject);
  const selectAffectedRecurrences = (edit: boolean) => { return selectAffectedRecurrencesModalInternal(edit); };
  setContext("selectAffectedRecurrences", selectAffectedRecurrences);

  let showDateModalInternal: (date: Date, events: (EventModel | null)[]) => any = $state(NoOp);
  const showDateModal = (date: Date, events: (EventModel | null)[]) => { return showDateModalInternal(date, events); };
  setContext("showDateModal", showDateModal);

  let showSettingsModalInternal: () => any = $state(NoOp);
  const showSettingsModal = () => { return showSettingsModalInternal(); }

  let showCreditsModalInternal: () => any = $state(NoOp);
  const showCreditsModal = () => { return showCreditsModalInternal(); };

  let showCreatePopup: () => any = $state(NoOp);

  let addButton: HTMLElement | undefined = $state();
</script>

<style lang="scss">
  @use "$lib/styles/animations.scss";
  @use "$lib/styles/colors.scss";
  @use "$lib/styles/dimensions.scss";
  @use "$lib/styles/text.scss";

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
    min-width: 10rem;
    width: 20vw;
    max-width: 20rem;
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

<SourceWizardModal bind:showModal={showSourceWizardModalInternal}/>
<SourceModal bind:showModal={showSourceModalInternal}/>
<CalendarModal bind:showModal={showCalendarModalInternal}/>
<EventModal bind:showModal={showEventModalInternal}/>
<AffectedRecurrencesModal bind:showModal={selectAffectedRecurrencesModalInternal}/>
<DayViewModal bind:showModal={showDateModalInternal}/>
<SettingsModal bind:showModal={showSettingsModalInternal}/>
<CreditsModal bind:showModal={showCreditsModalInternal}/>

<aside>
  <Title>{t("branding.name")}</Title>

  {#if settings.userSettings[UserSettingKeys.DisplaySmallCalendar]}
    <SmallCalendar date={date} smaller={true} onDayClick={(clickedDate) => smallCalendarClick(clickedDate)}></SmallCalendar>
  {/if}

  <div
    class="sources"
    aria-live="polite"
    aria-relevant="all"
    aria-busy={isLoading}
    role="tree"
  >
    {@render sourceEntries(getRepository().sources)}
  </div>

  <Horizontal position="center">
    <IconButton onClick={showSettingsModal} alt={t("button.settings")}>
      <Settings/>
    </IconButton>
    <IconButton bind:button={addButton} onClick={() => showCreatePopup().catch(NoOp)} alt={t("button.add.generic")}>
      <PlusIcon/>
    </IconButton>
    <CreatePopup
      bind:showPopup={showCreatePopup}
      anchor={addButton}
      addSource={showSourceWizardModal}
      addCalendar={showCalendarModal}
      addEvent={showEventModal}
    />
    <IconButton onClick={showCreditsModal} alt={t("button.credits")}>
      <Copyleft/>
    </IconButton>
  </Horizontal>

  <span class="copyright">
    {t("branding.copyright")}
  </span>
</aside>

<main>
  <div class="toprow">
    <MonthSelection bind:date granularity={view} />
    <Horizontal position="justify" width="auto">
      {#if connectivity.reachable != Reachability.Database}
        <span class="reachability">
          {#if connectivity.reachable == Reachability.Backend}
            {t("reachability.database")}
          {:else if connectivity.reachable == Reachability.Frontend}
            {t("reachability.backend")}
          {:else if connectivity.reachable == Reachability.None}
            {t("reachability.frontend")}
          {:else if connectivity.reachable == Reachability.Incompatible}
            {t("reachability.compatibility")}
          {:else}
            {t("reachability.unknown")}
          {/if}
          <WifiOff size={20}/>
        </span>
      {/if}

      <IconButton onClick={forceRefresh} alt={t("button.refresh")}>
        <span class="refreshButtonWrapper" class:spin={loaderAnimation} onanimationiteration={() => { if (!isLoading) loaderAnimation = false; }}>
          <RefreshCw size={20}/>
        </span>
      </IconButton>

      {#if !todayInRange}
        <Button onClick={seeToday} compact={true}>
          {t("scope.today")}
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
          { value: "day", name: t("scope.day")},
          { value: "week", name: t("scope.week")},
          { value: "month", name: t("scope.month")},
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
  {#each sources as source (source.id)}
    {@const index = repository.sources.findIndex((x) => x.id === source.id)}
    {@const calendars = repository.calendars.filter(cal => cal.source === source.id) || []}
    <SourceEntry bind:source={repository.sources[index]} calendars={calendars} idSeed="main"/>
    {#if !metadata.collapsedSources.has(repository.sources[index].id)}
        {@render calendarEntries(calendars)}
    {/if}
  {/each}
{/snippet}

{#snippet calendarEntries(calendars: CalendarModel[])}
  {#each calendars as calendar (calendar.id)}
    {@const index = repository.calendars.findIndex((x) => x.id === calendar.id)}
    <CalendarEntry bind:calendar={repository.calendars[index]} idSeed="main"/>
  {/each}
{/snippet}
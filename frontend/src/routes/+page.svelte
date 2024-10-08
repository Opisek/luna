<script lang="ts">
  import Calendar from "../components/calendar/Calendar.svelte";
  import IconButton from "../components/interactive/IconButton.svelte";

  import { browser } from "$app/environment";
  import SourceEntry from "../components/calendar/SourceEntry.svelte";
  import { calendars, events, fetchAllEvents, fetchSources, sources } from "$lib/client/repository";
  import { queueNotification } from "$lib/client/notifications";
  import CalendarEntry from "../components/calendar/CalendarEntry.svelte";
  import Title from "../components/layout/Title.svelte";
  import Horizontal from "../components/layout/Horizontal.svelte";
  import { PlusIcon } from "lucide-svelte";
  import SourceModal from "../components/modals/SourceModal.svelte";
  import MonthSelection from "../components/interactive/MonthSelection.svelte";

  let localSources: SourceModel[] = [];
  let localCalendars: CalendarModel[] = [];
  let sourceCalendars: Map<string, CalendarModel[]> = new Map();
  let localEvents: EventModel[] = [];
  let calendarEvents: Map<string, EventModel[]> = new Map();

  let showNewSourceModal: () => any;
  let newSource: SourceModel;

  function createNewSource() {
    newSource = {
      id: "",
      name: "",
      type: "caldav",
      settings: {},
      auth_type: "none",
      auth: {},
      collapsed: false
    };
    setTimeout(showNewSourceModal, 0);
  }

  const currentYear = new Date().getFullYear();
  const currentMonth = new Date().getMonth() 
  let selectedYear = currentYear;
  let selectedMonth = currentMonth;

  let rangeStart = new Date(selectedYear, selectedMonth - 1, 1);
  let rangeEnd = new Date(selectedYear, selectedMonth + 2, 0);

  (async () => {
    if (!browser) return;

    fetchSources().then(err => {
      if (err != "") {
        queueNotification(
          "failure",
          `Failed to fetch sources: ${err}`
        );
      }
    });

    events.subscribe((newEvents) => {
      localEvents = newEvents;

      calendarEvents = new Map();
      localEvents.forEach((event) => {
        if (calendarEvents.has(event.calendar)) {
          // @ts-ignore typescript says that this might be undefined despite the check above
          calendarEvents.get(event.calendar).push(event);
        } else {
          calendarEvents.set(event.calendar, [ event ]);
        }
      });
    });

    calendars.subscribe((newCalendars) => {
      localCalendars = newCalendars;

      sourceCalendars = new Map();
      localCalendars.forEach((calendar) => {
        if (sourceCalendars.has(calendar.source)) {
          // @ts-ignore typescript says that this might be undefined despite the check above
          sourceCalendars.get(calendar.source).push(calendar);
        } else {
          sourceCalendars.set(calendar.source, [ calendar ]);
        }
      });
    });

    sources.subscribe((newSources) => {
      localSources = newSources;
    });
  })();

  $: ((year: number, month: number) => {
    const firstDayNextMonth = new Date(year, month + 1, 1);
    const lastDayPreviousMonth = new Date(year, month, 0);

    if (lastDayPreviousMonth < rangeStart) {
      const lastStart = rangeStart;
      rangeStart = new Date(year, month - 2, 1);
      fetchAllEvents(rangeStart, lastStart);
    }

    if (firstDayNextMonth > rangeEnd) {
      const lastEnd = rangeEnd;
      rangeEnd = new Date(year, month + 3, 0);
      fetchAllEvents(lastEnd, rangeEnd);
    }
  })(selectedYear, selectedMonth);
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
</style>

<div class="wrapper">
  <aside>
    <Title>Luna</Title>

    <!-- SmallCalendar put here only for testing purposes but might consider leaving it if it can serve some useful purpose -->
    <!--<SmallCalendar year={selectedYear} month={selectedMonth}/>-->

    <div class="sources">
      {#each localSources as source}
        <SourceEntry bind:source={source}/>
        {#if (!source.collapsed)}
          {#each sourceCalendars.get(source.id) || [] as calendar}
            <CalendarEntry calendar={calendar}/>
          {/each}
        {/if}
      {/each}
    </div>
    <Horizontal position="center">
      <IconButton click={createNewSource}>
        <PlusIcon/>
      </IconButton>
      <SourceModal bind:showCreateModal={showNewSourceModal} source={newSource}/>
    </Horizontal>
  </aside>
  <main>
    <MonthSelection bind:month={selectedMonth} bind:year={selectedYear}/>
    <Calendar
      year={selectedYear}
      month={selectedMonth}
      events={localEvents}
    />
  </main>
</div>
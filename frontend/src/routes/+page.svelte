<script lang="ts">
  import Calendar from "../components/calendar/Calendar.svelte";
  import IconButton from "../components/interactive/IconButton.svelte";
  import LeftIcon from "lucide-svelte/icons/chevron-left";
  import RightIcon from "lucide-svelte/icons/chevron-right";

  import { browser } from "$app/environment";
  import { getMonthName } from "../lib/common/humanization";
  import CalendarEntry from "../components/interactive/CalendarEntry.svelte";
  import SourceRow from "../components/calendar/SourceRow.svelte";

  let sources: SourceModel[] = [];
  let calendars: CalendarModel[] = [];
  let sourceCalendars: Map<string, CalendarModel[]> = new Map();
  let events: EventModel[] = [];
  let calendarEvents: Map<string, EventModel[]> = new Map();

  const currentYear = new Date().getFullYear();
  const currentMonth = new Date().getMonth() 
  let selectedYear = currentYear;
  let selectedMonth = currentMonth;

  function previousMonth() {
    selectedMonth--;
    if (selectedMonth === -1) {
      selectedMonth = 11;
      selectedYear--;
    }
  }

  function nextMonth() {
    if (selectedMonth === 11) selectedYear++;
    selectedMonth = (selectedMonth + 1) % 12;
  }

  async function fetchSources(): Promise<SourceModel[]> {
    const response = await fetch("/api/sources");
    if (response.ok) return response.json();
    else {
      console.log("Failed to fetch sources");
      console.log(response);
      return [];
    }
  }

  async function fetchCalendars(): Promise<CalendarModel[]> {
    const response = await fetch("/api/calendars");
    if (response.ok) return response.json();
    else {
      console.log("Failed to fetch calendars");
      console.log(response);
      return [];
    }
  }

  async function fetchEvents(): Promise<EventModel[]> {
    const response = await fetch("/api/events");
    if (response.ok) {
      const events = await response.json();
      for (const event of events) {
        event.date.start = new Date(event.date.start);
        event.date.end = new Date(event.date.end);
      }
      return events;
    } else {
      console.log("Failed to fetch events");
      console.log(response);
      return [];
    }
  }

  (async () => {
    if (!browser) return;
    sources = await fetchSources();

    calendars = await fetchCalendars();
    sourceCalendars = new Map();
    calendars.forEach((calendar) => {
      if (sourceCalendars.has(calendar.source)) {
        // @ts-ignore typescript says that this might be undefined despite the check above
        sourceCalendars.get(calendar.source).push(calendar);
      } else {
        sourceCalendars.set(calendar.source, [ calendar ]);
      }
    });

    events = await fetchEvents();
  })();
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

  div.monthSelection {
    display: flex;
    flex-direction: row;
    gap: $gapSmall;
    align-items: center;
  }

  div.wrapper {
    display: flex;
    flex-direction: row;
    gap: $gap;
    padding: $gap;
    height: 100%;
    width: 100%;
  }
</style>

<div class="wrapper">
  <aside>
    <h1>Calendars</h1>
    {#each sources as source}
      <SourceRow source={source}/>
      <!--<SectionTitle title={source.name} />-->
      {#each sourceCalendars.get(source.id) || [] as calendar}
        <CalendarEntry calendar={calendar}/>
      {/each}
    {/each}
  </aside>
  <main>
    <div class="monthSelection">
      <IconButton callback={previousMonth}>
        <LeftIcon/>
      </IconButton>
      <IconButton callback={nextMonth}>
        <RightIcon/>
      </IconButton>
      <span class="monthLabel">
        {`${getMonthName(selectedMonth)} ${selectedYear}`}
      </span>
    </div>

    <Calendar
      year={selectedYear}
      month={selectedMonth}
      events={events}
    />
  </main>
</div>
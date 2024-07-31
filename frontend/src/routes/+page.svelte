<script lang="ts">
  import Calendar from "../components/calendar/Calendar.svelte";
  import IconButton from "../components/interactive/IconButton.svelte";
  import LeftIcon from "lucide-svelte/icons/chevron-left";
  import RightIcon from "lucide-svelte/icons/chevron-right";

  import { browser } from "$app/environment";
  import { getMonthName } from "../lib/common/humanization";
    import CalendarEntry from "../components/interactive/CalendarEntry.svelte";

  let calendars: CalendarModel[] = [];

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

  async function fetchCalendars(): Promise<CalendarModel[]> {
    const response = await fetch("/api/calendars");
    if (response.ok) return response.json();
    else {
      console.log("Failed to fetch calendars");
      console.log(response);
      return [];
    }
  }

  (async () => {
    if (!browser) return;
    calendars = await fetchCalendars();
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
    {#each calendars as calendar}
      <CalendarEntry calendar={calendar}/>
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
      year={2024}
      month={selectedMonth}
      events={[
        {
          title: "Event 1",
          start: new Date(2024, 3, 1),
          end: new Date(2024, 3, 3),
          allDay: false,
          color: "red"
        },
        {
          title: "Event 2",
          start: new Date(2024, 3, 2),
          end: new Date(2024, 3, 4),
          allDay: false,
          color: "blue"
        },
        {
          title: "Event 3",
          start: new Date(2024, 3, 3),
          end: new Date(2024, 3, 5),
          allDay: false,
          color: "green"
        },
        {
          title: "Event 4",
          start: new Date(2024, 3, 4),
          end: new Date(2024, 3, 6),
          allDay: false,
          color: "yellow"
        },
        {
          title: "Event 5",
          start: new Date(2024, 3, 7),
          end: new Date(2024, 3, 9),
          allDay: false,
          color: "yellow"
        }
      ]}
    />
  </main>
</div>
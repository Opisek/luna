<script>
  import Calendar from "../components/calendar/Calendar.svelte";
  import IconButton from "../components/interactive/IconButton.svelte";
  import { getMonthName } from "../lib/common/humanization";
  import LeftIcon from "lucide-svelte/icons/chevron-left";
  import RightIcon from "lucide-svelte/icons/chevron-right";

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
</script>

<style lang="scss">
  @import "../styles/dimensions.scss";

  div {
    display: flex;
    flex-direction: row;
    gap: 1em;
  }
  main {
    width: 100vw;
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
</style>

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
        start: new Date(2024, 0, 1),
        end: new Date(2024, 0, 3),
        allDay: false,
        color: "red"
      },
      {
        title: "Event 2",
        start: new Date(2024, 0, 2),
        end: new Date(2024, 0, 4),
        allDay: false,
        color: "blue"
      }
    ]}
  />
</main>
<script>
  import Calendar from "../components/calendar/Calendar.svelte";
    import { getMonthName } from "../lib/common/humanization";

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
  div {
    display: flex;
    flex-direction: row;
    gap: 1em;
  }
</style>

<div>
  <button on:click={previousMonth}>←</button>
  <h1>
    {getMonthName(selectedMonth) + (currentYear === selectedYear ? "" : ` ${selectedYear}`)}
  </h1>
  <button on:click={nextMonth}>→</button>
</div>
<Calendar year={2024} month={selectedMonth}/>
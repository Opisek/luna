<script lang="ts">
  import MonthSelection from "../interactive/MonthSelection.svelte";
  import SmallCalendar from "../interactive/SmallCalendar.svelte";
  import Modal from "./Modal.svelte";

  export let date: Date;

  let currentYear: number;
  let currentMonth: number;

  export const showModal = () => {
    currentYear = date.getFullYear();
    currentMonth = date.getMonth();
    showModalInternal();
  };

  let showModalInternal: () => any;
  let hideModalInternal: () => any;

  function dateSelected(selectedDate: Date) {
    // keep the time of day the same:
    selectedDate.setHours(date.getHours(), date.getMinutes(), date.getSeconds(), date.getMilliseconds());
    date = selectedDate;
    hideModalInternal();
  }
</script>

<style lang="scss">

</style>

<Modal title="Pick Date" bind:showModal={showModalInternal} bind:hideModal={hideModalInternal}>
  <MonthSelection bind:month={currentMonth} bind:year={currentYear} />
  <SmallCalendar year={currentYear} month={currentMonth} onDayClick={dateSelected} />
</Modal>
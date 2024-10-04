<script lang="ts">
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
</script>

<style lang="scss">

</style>

<Modal title="Pick Date" bind:showModal={showModalInternal} bind:hideModal={hideModalInternal}>
  <SmallCalendar year={currentYear} month={currentMonth} />
</Modal>
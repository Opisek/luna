<script lang="ts">
  import Modal from "./Modal.svelte";

  export let date: Date;

  export const showModal = () => {
    showModalInternal();
  };

  let showModalInternal: () => any;
  let hideModalInternal: () => any;

  function dateSelected(selectedDate: Date) {
    selectedDate.setDate(date.getDate());
    selectedDate.setMonth(date.getMonth());
    selectedDate.setFullYear(date.getFullYear());
    date = selectedDate;
    hideModalInternal();
  }
</script>

<style lang="scss">

</style>

<Modal title="Pick Time" bind:showModal={showModalInternal} bind:hideModal={hideModalInternal}>

</Modal>
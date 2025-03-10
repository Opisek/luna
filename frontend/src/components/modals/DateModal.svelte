<script lang="ts">
  import Modal from "./Modal.svelte";
  import MonthSelection from "../interactive/MonthSelection.svelte";
  import SmallCalendar from "../interactive/SmallCalendar.svelte";

  import { NoOp } from "$lib/client/placeholders";

  interface Props {
    date: Date;
    onChange?: (date: Date) => void;
    showModal?: () => any;
    hideModal?: () => any;
  }

  let {
    date = $bindable(new Date()),
    onChange = NoOp,
    showModal = $bindable(),
    hideModal = $bindable()
  }: Props = $props();

  let showModalInternal: () => any = $state(NoOp);
  let hideModalInternal: () => any = $state(NoOp);

  showModal = () => {
    showModalInternal();
  };

  hideModal = () => {
    hideModalInternal();
  };

  function dateSelected(selectedDate: Date) {
    // keep the time of day the same:
    selectedDate.setHours(date.getHours(), date.getMinutes(), date.getSeconds(), date.getMilliseconds());
    date = selectedDate;
    hideModalInternal();
    onChange(date);
  }
</script>

<Modal title="Pick Date" bind:showModal={showModalInternal} bind:hideModal={hideModalInternal}>
  <MonthSelection bind:date />
  <SmallCalendar bind:date onDayClick={dateSelected} />
</Modal>
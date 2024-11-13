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
    date = $bindable(),
    onChange = NoOp,
    showModal = $bindable(),
    hideModal = $bindable()
  }: Props = $props();

  let currentYear: number = $state(0);
  let currentMonth: number = $state(0);

  let showModalInternal: () => any = $state(NoOp);
  let hideModalInternal: () => any = $state(NoOp);

  showModal = () => {
    currentYear = date.getFullYear();
    currentMonth = date.getMonth();
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

<Modal title="Pick Date" bind:showModal={showModalInternal} bind:hideModal={hideModal}>
  <MonthSelection bind:month={currentMonth} bind:year={currentYear} />
  <SmallCalendar year={currentYear} month={currentMonth} onDayClick={dateSelected} />
</Modal>
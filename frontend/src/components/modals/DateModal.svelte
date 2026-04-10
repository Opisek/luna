<script lang="ts">
  import Modal from "./Modal.svelte";
  import MonthSelection from "../interactive/MonthSelection.svelte";
  import SmallCalendar from "../interactive/SmallCalendar.svelte";

  import { NoOp } from "$lib/client/placeholders";
  import { t } from "@sveltia/i18n";

  interface Props {
    showModal: (initial: Date) => Promise<Date>;
  }

  let {
    showModal = $bindable(),
  }: Props = $props();

  let showModalInternal: () => Promise<Date> = $state(Promise.reject);
  let success: (result: Date) => void = $state(NoOp);

  let date = $state(new Date());

  showModal = (initial: Date) => {
    date = new Date(initial); 
    return showModalInternal();
  };

  function dateSelected(selectedDate: Date) {
    // keep the time of day the same:
    selectedDate.setHours(date.getHours(), date.getMinutes(), date.getSeconds(), date.getMilliseconds());
    success(selectedDate);
  }
</script>

<Modal title={t("date.title")} bind:showModal={showModalInternal} bind:success>
  <MonthSelection bind:date />
  <SmallCalendar bind:date onDayClick={dateSelected} />
</Modal>
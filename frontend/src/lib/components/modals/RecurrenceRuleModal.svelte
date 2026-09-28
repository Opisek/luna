<script lang="ts">
  import Modal from "./Modal.svelte";

  import { NoOp } from "$lib/scripts/client/placeholders";
  import { RRule, type Options } from "rrule";
  import RecurrenceInput from "../forms/RecurrenceInput.svelte";

  interface Props {
    showModal: (initial: Partial<Options>) => Promise<Partial<Options>>;
    dtstart: Date;
    allDay: boolean; }

  let success: (result: Partial<Options>) => void = $state(NoOp);
  let failure: (reason?: string | Error) => void = $state(NoOp);

  let {
    showModal = $bindable(),
    dtstart,
    allDay
  }: Props = $props();

  let showModalInternal: () => Promise<Partial<Options>> = $state(Promise.reject);
  let options = $state<Partial<Options>>((new RRule()).origOptions)

  showModal = async (initial) => {
    options = initial;
    return showModalInternal();
  }
</script>

<Modal title={"Recurrence editing"} bind:showModal={showModalInternal} bind:success bind:failure>
  <RecurrenceInput
    dtstart={dtstart} 
    bind:options={options}
    allDay={allDay}
    editable={true}
  />
</Modal>
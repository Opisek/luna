<script lang="ts">
  import Event from "../calendar/Event.svelte";
  import Modal from "./Modal.svelte";

  import { NoOp } from "$lib/client/placeholders";

  interface Props {
    date: Date;
    events: (EventModel | null)[];
    showModal?: () => any;
    hideModal?: () => any;
  }

  let {
    date,
    events,
    showModal = $bindable(),
    hideModal = $bindable(NoOp),
  }: Props = $props();
</script>

<style lang="scss">
  @use "../../styles/dimensions.scss";

  .wrapper {
    height: fit-content;
    display: flex;
    flex-direction: column;
    gap: dimensions.$gapSmaller;
  }
</style>

<Modal title={date.toDateString()} bind:showModal={showModal} bind:hideModal={hideModal}>
  <div class="wrapper">
    {#each events as event, i ((event?.id || 0) + i.toString())}
      <Event
        event={event}
        isFirstDay={true}
        isLastDay={true}
        date={date}
        visible={true}
      />
    {/each}
  </div>
</Modal>
<script lang="ts">
  import Event from "../calendar/Event.svelte";
  import Modal from "./Modal.svelte";

  import { setContext } from "svelte";

  interface Props {
    showModal?: (date: Date, events: (EventModel | null)[]) => Promise<void>;
  }

  let {
    showModal = $bindable(),
  }: Props = $props();

  let date = $state(new Date());
  let events: (EventModel | null)[] = $state([]);

  let showModalInternal: () => Promise<void> = $state(Promise.reject);
  showModal = (_date: Date, _events: (EventModel | null)[]) => {
    date = _date;
    events = _events;
    return showModalInternal();
  };

  let currentlyClickedEvent = $state<EventModel | null>(null);
  let currentlyHoveredEvent = $state<EventModel | null>(null);
  setContext("currentlyClickedEvent", () => currentlyClickedEvent);
  setContext("currentlyHoveredEvent", () => currentlyHoveredEvent);
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

<Modal title={date.toDateString()} bind:showModal={showModalInternal}>
  {#if events.length === 0}
    No events
  {:else}
    <div class="wrapper">
      {#each events as event, i ((event?.id || 0) + i.toString())}
        <Event
          event={event}
          isFirstDay={true}
          date={date}
          visible={true}
          view="day"
        />
      {/each}
    </div>
  {/if}
</Modal>
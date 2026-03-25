<script lang="ts">
  import Popup from "./Popup.svelte";

  import { AsyncNoOp, NoOp } from '$lib/client/placeholders';
  import Button from "../interactive/Button.svelte";
  import { getRepository } from "../../lib/client/data/repository.svelte";

  interface Props {
    showPopup?: () => Promise<void>;
    hidePopup?: () => void;
    addSource: () => Promise<void>;
    addCalendar: () => Promise<void>;
    addEvent: () => Promise<void>;
  }

  let {
    showPopup = $bindable(AsyncNoOp),
    hidePopup = $bindable(NoOp),
    addSource,
    addCalendar,
    addEvent,
  }: Props = $props();

  const repository = getRepository();
  const canAddCalendars = $derived(repository.sources.some(x => x.can_add_calendars));
  const canAddEvents = $derived(repository.calendars.some(x => x.can_add_events));

  let internalShow: () => Promise<void> = $state(AsyncNoOp);
  let internalClose: () => void = $state(NoOp);

  /* Popup */
  showPopup = async () => {
    if (!canAddCalendars && !canAddEvents) return addSource();
    return internalShow();
  }

  hidePopup = () => {
    internalClose();
  }

  /* Buttons */
  function onAddSourceButtonClick() {
    addSource().finally(hidePopup);
  }

  function onAddCalendarButtonClick() {
    addCalendar().finally(hidePopup);
  }

  function onAddEventButtonClick() {
    addEvent().finally(hidePopup);
  }
</script>

<style lang="scss">
  @use "../../styles/animations.scss";
  @use "../../styles/colors.scss";
  @use "../../styles/dimensions.scss";

  div.body {
    overflow: hidden;
    position: relative;
  }

  div.grid {
    display: grid;
    gap: dimensions.$gapSmall;
  }

  div.grid:not(:first-child) {
    position: absolute;
    top: 0;
    left: 0;
  }

  div.grid.month {
    grid-template-columns: repeat(4, 1fr);
    grid-template-rows: repeat(3, 1fr);
  }

  div.grid.year {
    grid-template-columns: repeat(5, 1fr);
    grid-template-rows: repeat(2, 1fr);
  }

  div.topRow {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  button.display {
    all: unset;
    display: flex;
    justify-content: center;
    align-items: center;
    cursor: pointer;
    user-select: none;
    position: relative;
  } 

  button.button {
    all: unset;
    display: flex;
    justify-content: center;
    align-items: center;
    border-radius: dimensions.$borderRadiusSmall;
    //color: colors.$foregroundTertiary;
    //background-color: colors.$backgroundTertiary;
    padding: dimensions.$gapSmall;
    cursor: pointer;
    user-select: none;
    position: relative;
    overflow: hidden;
  }

  button.month {
    width: 2em;
  }

  button.year {
    width: 3em;
  }
</style>

<Popup bind:showPopup={internalShow} bind:hidePopup={internalClose} tooltip={false}>
  <Button onClick={onAddSourceButtonClick}>Add Source</Button>
  {#if canAddCalendars}
    <Button onClick={onAddCalendarButtonClick}>Add Calendar</Button>
  {/if}
  {#if canAddEvents}
    <Button onClick={onAddEventButtonClick}>Add Event</Button>
  {/if}
</Popup>
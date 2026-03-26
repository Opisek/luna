<script lang="ts">
  import Popup from "./Popup.svelte";

  import { AsyncNoOp, NoOp } from '$lib/client/placeholders';
  import Button from "../interactive/Button.svelte";
  import { getRepository } from "../../lib/client/data/repository.svelte";

  interface Props {
    showPopup?: () => Promise<void>;
    hidePopup?: () => void;
    addSource: () => Promise<SourceModel>;
    addCalendar: () => Promise<CalendarModel>;
    addEvent: () => Promise<EventModel>;
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
    if (!canAddCalendars && !canAddEvents) return addSource().catch(NoOp).then(NoOp);
    return internalShow();
  }

  hidePopup = () => {
    internalClose();
  }

  /* Buttons */
  function onAddSourceButtonClick() {
    addSource().catch(NoOp).finally(hidePopup);
  }

  function onAddCalendarButtonClick() {
    addCalendar().catch(NoOp).finally(hidePopup);
  }

  function onAddEventButtonClick() {
    addEvent().catch(NoOp).finally(hidePopup);
  }
</script>

<Popup bind:showPopup={internalShow} bind:hidePopup={internalClose} tooltip={false}>
  <Button onClick={onAddSourceButtonClick}>Add Source</Button>
  {#if canAddCalendars}
    <Button onClick={onAddCalendarButtonClick}>Add Calendar</Button>
  {/if}
  {#if canAddEvents}
    <Button onClick={onAddEventButtonClick}>Add Event</Button>
  {/if}
</Popup>
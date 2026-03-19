<script lang="ts">
  import ColorCircle from "../misc/ColorCircle.svelte";
  import Spinner from "../decoration/Spinner.svelte";
  import Tooltip from "../interactive/Tooltip.svelte";
  import VisibilityToggle from "../interactive/VisibilityToggle.svelte";

  import { GetCalendarColor } from "$lib/common/colors";
  import { NoOp } from "$lib/client/placeholders";
  import { focusIndicator } from "$lib/client/decoration";
  import { getMetadata } from "$lib/client/data/metadata.svelte";
  import { draggable } from "$lib/client/reordering";

  import { getContext } from "svelte";
  import { queueNotification } from "../../lib/client/notifications";
  import { ColorKeys } from "../../types/colors";
  import { getRepository } from "../../lib/client/data/repository.svelte";

  interface Props {
    calendar: CalendarModel;
  }

  let { calendar = $bindable() }: Props = $props();

  const metadata = getMetadata();
  const repository = getRepository();

  let hasErrored = $derived(calendar && metadata.faultyCalendars.has(calendar.id));
  let isLoading = $derived(calendar && metadata.loadingCalendars.get(calendar.id));
  let calendarVisible = $derived(calendar && metadata.hiddenCalendars.has(calendar.id));

  let showModal: ((calendar: CalendarModel) => Promise<CalendarModel>) = getContext("showCalendarModal");
  function showModalInternal() {
    showModal(calendar).then(newCalendar => calendar = newCalendar).catch(NoOp);
  }

  $effect(() => {
    const shouldBeVisible = !metadata.hiddenCalendars.has(calendar.id);
    if (shouldBeVisible == calendarVisible) return;
    calendarVisible = shouldBeVisible;
  });
  function setVisible(visible: boolean) {
    getMetadata().setCalendarVisibility(calendar.id, visible);
  }

  function reorderCalendar(newIndex: number) {
    repository.changeCalendarDisplayOrder(calendar, newIndex).catch((err) => {
      queueNotification(ColorKeys.Danger, err);
    });
  }
</script>

<style lang="scss">
  @use "../../styles/dimensions.scss";

  div.calendarEntry {
    display: flex;
    flex-direction: row;
    gap: dimensions.$gapTiny;
    width: 100%;
    align-items: center;
    justify-content: space-between;
    user-select: none;
    cursor: grab;
  }

  span {
    display: flex;
    flex-direction: row;
    align-items: center;
  }

  span.name {
    gap: dimensions.$gapSmall;
    min-width: 0;
  }

  span.buttons {
    gap: dimensions.$gapTiny;
  }

  button {
    all: unset;
    cursor: pointer;
    display: inline;
    width: max-content;
    position: relative;
    text-wrap: nowrap;
    text-overflow: ellipsis;
    min-width: 0;
    overflow: hidden;
  }
</style>

<div class="calendarEntry" use:draggable={{ ownClass: "calendarEntry", childClasses: [], callback: reorderCalendar}}>
  <span class="name">
    <ColorCircle
      color={GetCalendarColor(calendar)}
      size="small"
    />
    <button onclick={showModalInternal} use:focusIndicator={{ type: "underline" }}>
      {calendar.name}
    </button>
  </span>
  <span class="buttons">
    {#if isLoading}
      <Spinner/>
    {/if}
    <VisibilityToggle bind:visible={calendarVisible} onClick={setVisible}/>
    {#if hasErrored}
      <Tooltip error={true}>An error occurred trying to retrieve events from this calendar.</Tooltip>
    {/if}
  </span>
</div>
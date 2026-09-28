<script lang="ts">
  import ColorCircle from "../misc/ColorCircle.svelte";
  import Spinner from "../decoration/Spinner.svelte";
  import Tooltip from "../interactive/Tooltip.svelte";
  import VisibilityToggle from "../interactive/VisibilityToggle.svelte";

  import { GetCalendarColor } from "$lib/scripts/common/colors";
  import { NoOp } from "$lib/scripts/client/placeholders";
  import { focusIndicator } from "$lib/scripts/client/decoration";
  import { getMetadata } from "$lib/scripts/client/data/metadata.svelte";
  import { draggable } from "$lib/scripts/client/reordering.svelte";

  import { getContext } from "svelte";
  import { queueNotification } from "$lib/scripts/client/notifications";
  import { ColorKeys } from "$lib/types/colors";
  import { getRepository } from "$lib/scripts/client/data/repository.svelte";
  import { t } from "@sveltia/i18n";
  import { extendUniqueElementId, generateUniqueElementId } from "$lib/scripts/common/dom";

  interface Props {
    calendar: CalendarModel;
    idSeed: string;
  }

  let {
    calendar = $bindable(),
    idSeed,
  }: Props = $props();
  let calendarEntryId = $derived(generateUniqueElementId(["calendarentry", calendar.id], idSeed));

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

  async function reorderCalendar(newIndex: number) {
    await repository.changeCalendarDisplayOrder(calendar, newIndex).catch((err) => {
      queueNotification(ColorKeys.Danger, err);
    });
  }
</script>

<style lang="scss">
  @use "$lib/styles/dimensions.scss";

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

<div
  class="calendarEntry"
  use:draggable={{ ownClass: "calendarEntry", childClasses: [], callback: reorderCalendar}}
  id={calendarEntryId}
  aria-label={t("calendar.aria", { values: { name: calendar.name } })}
  aria-level="2"
  role="treeitem"
  aria-selected={document.activeElement?.id === calendarEntryId}
  aria-details={hasErrored ? extendUniqueElementId(["tooltip"], calendarEntryId) : undefined}
>
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
      <Tooltip error={true} describes={calendarEntryId}>{t("calendar.error.events.tooltip")}</Tooltip>
    {/if}
  </span>
</div>
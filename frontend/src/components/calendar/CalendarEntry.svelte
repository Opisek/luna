<script lang="ts">
  import ColorCircle from "../misc/ColorCircle.svelte";
  import Tooltip from "../interactive/Tooltip.svelte";
  import VisibilityToggle from "../interactive/VisibilityToggle.svelte";

  import { GetCalendarColor } from "$lib/common/colors";
  import { faultyCalendars } from "$lib/client/repository";
  import { focusIndicator } from "$lib/client/decoration";
  import { hiddenCalendars, isCalendarVisible, setCalendarVisibility } from "$lib/client/localStorage";

  interface Props {
    calendar: CalendarModel;
  }

  let { calendar = $bindable() }: Props = $props();

  let calendarVisible = $state(calendar ? isCalendarVisible(calendar.id) : false);

  let hasErrored = $state(false);
  faultyCalendars.subscribe((faulty) => {
    hasErrored = faulty.has(calendar.id);
  });

  $effect(() => {
    if (calendar && calendar.id) setCalendarVisibility(calendar.id, calendarVisible);
  });
</script>

<style lang="scss">
  @import "../../styles/dimensions.scss";

  div.entry {
    display: flex;
    flex-direction: row;
    gap: $gapTiny;
    width: 100%;
    align-items: center;
    justify-content: space-between;
  }

  span {
    display: flex;
    flex-direction: row;
    align-items: center;
  }

  span.name {
    gap: $gapSmall;
    min-width: 0;
  }

  span.buttons {
    gap: $gapTiny;
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

<div class="entry">
  <span class="name">
    <ColorCircle
      color={GetCalendarColor(calendar)}
      size="small"
    />
    <button use:focusIndicator={{ type: "underline" }}>
      {calendar.name}
    </button>
  </span>
  <span class="buttons">
    <VisibilityToggle bind:visible={calendarVisible}/>
    {#if hasErrored}
      <Tooltip msg="An error occurred trying to retrieve events from this calendar." error={true}/>
    {/if}
  </span>
</div>
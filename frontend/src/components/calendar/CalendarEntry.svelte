<script lang="ts">
  import { focusIndicator } from "$lib/client/decoration";
  import { faultyCalendars } from "$lib/client/repository";
  import { GetCalendarColor } from "$lib/common/colors";
  import { hiddenCalendars, setCalendarVisibility } from "../../lib/client/localStorage";
  import Tooltip from "../interactive/Tooltip.svelte";
  import VisibilityToggle from "../interactive/VisibilityToggle.svelte";
  import ColorCircle from "../misc/ColorCircle.svelte";

  export let calendar: CalendarModel;

  let hasErrored = false;
  faultyCalendars.subscribe((faulty) => {
    hasErrored = faulty.has(calendar.id);
  });

  $: if (calendar && calendar.id) setCalendarVisibility(calendar.id, calendar.visible);
  hiddenCalendars.subscribe((hidden) => {
    calendar.visible = !hidden.has(calendar.id);
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
    <VisibilityToggle bind:visible={calendar.visible}/>
    {#if hasErrored}
      <Tooltip msg="An error occurred trying to retrieve events from this calendar." error={true}/>
    {/if}
  </span>
</div>
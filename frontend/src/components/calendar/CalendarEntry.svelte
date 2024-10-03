<script lang="ts">
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
    gap: $gapSmall;
    width: 100%;
    align-items: center;
  }


  span {
    flex-grow: 1;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
</style>

<div class="entry">
  <ColorCircle
    color={GetCalendarColor(calendar)}
    size="small"
  />
  <span>
    {calendar.name}
  </span>
  <VisibilityToggle bind:visible={calendar.visible}/>
  {#if hasErrored}
    <Tooltip msg="An error occurred trying to retrieve events from this calendar." error={true}/>
  {/if}
</div>
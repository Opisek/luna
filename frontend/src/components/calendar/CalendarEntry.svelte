<script lang="ts">
  import { faultyCalendars } from "$lib/client/repository";
  import Tooltip from "../interactive/Tooltip.svelte";
  import VisibilityToggle from "../interactive/VisibilityToggle.svelte";

  export let calendar: CalendarModel;

  let hasErrored = false;
  faultyCalendars.subscribe((faulty) => {
    hasErrored = faulty.has(calendar.id);
  });

  let visibleTest = true;
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

  div.color {
    width: 0.75em;
    height: 0.75em;
    border-radius: 50%;
  }

  span {
    flex-grow: 1;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
</style>

<div class="entry">
  <div class="color" style="background-color:{calendar.color}"></div>
  <span>
    {calendar.name}
  </span>
  <VisibilityToggle visible={visibleTest}/>
  {#if hasErrored}
    <Tooltip msg="An error occurred trying to retrieve events from this calendar." error={true}/>
  {/if}
</div>
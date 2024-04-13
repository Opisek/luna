<script lang="ts">
  import CalendarEvent from "./CalendarEvent.svelte";
  import { compareEventsByStartDate } from "../../lib/common/comparators";

  export let day: number;
  export let dayOfWeek: number;
  export let isCurrentMonth: boolean;

  export let events: CalendarEventModel[];
</script>

<style lang="scss">
  @import "../../styles/dimensions.scss";

  div.day {
    display: flex;
    flex-direction: column;
    gap: $gapSmall;
    padding: $gapSmall;
    border-radius: $borderRadiusSmall;
    background-color: #f0f0f0;
  }
  div.otherMonth {
    opacity: .5;
  }
  div.sunday {
    color: red;
  }

  span {
    text-align: center;
    width: 100%;
    display: block;
  }

  div.events {
    display: flex;
    flex-direction: column;
    gap: $gapSmall;
  }
</style>

<div class="day" class:otherMonth={!isCurrentMonth} class:sunday={dayOfWeek === 0}>
  <span>
    {day}
  </span>
  <div class="events">
    {#each events.sort(compareEventsByStartDate) as event}
      <CalendarEvent {event}/>
    {/each}
  </div>
</div>
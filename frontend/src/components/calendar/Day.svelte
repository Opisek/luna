<script lang="ts">
  import CalendarEvent from "./CalendarEvent.svelte";

  export let date: Date;

  export let isCurrentMonth: boolean;

  export let isFirstDay: boolean;
  export let isLastDay: boolean;

  export let events: (CalendarEventModel | null)[];
  
  export let currentlyHoveredEvent: CalendarEventModel | null;
  export let currentlyClickedEvent: CalendarEventModel | null;
  export let clickCallback: (event: CalendarEventModel) => void;
</script>

<style lang="scss">
  @import "../../styles/dimensions.scss";

  div.day {
    display: flex;
    flex-direction: column;
    gap: $gapSmall;
    padding-top: $paddingSmaller;
    border-radius: $borderRadiusSmall;
    background-color: #f0f0f0;
  }
  div.otherMonth {
    opacity: .5;
  }

  span {
    text-align: center;
    width: 100%;
    display: block;
  }
  div.sunday > span {
    color: red;
  }

  div.events {
    display: flex;
    flex-direction: column;
    gap: $gapTiny;
  }
</style>

<div class="day" class:otherMonth={!isCurrentMonth} class:sunday={date.getDay() === 0}>
  <span>
    {date.getDate()}
  </span>
  <div class="events">
    {#each events as event}
      <CalendarEvent
        event={event}
        isFirstDay={isFirstDay}
        isLastDay={isLastDay}
        date={date}
        bind:currentlyHoveredEvent={currentlyHoveredEvent}
        bind:currentlyClickedEvent={currentlyClickedEvent}
        clickCallback={clickCallback}
      />
    {/each}
  </div>
</div>
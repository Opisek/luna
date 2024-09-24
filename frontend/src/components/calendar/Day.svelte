<script lang="ts">
  import CalendarEvent from "./CalendarEvent.svelte";

  export let date: Date;

  export let isCurrentMonth: boolean;

  export let isFirstDay: boolean;
  export let isLastDay: boolean;

  export let events: (EventModel | null)[];
  
  export let currentlyHoveredEvent: EventModel | null;
  export let currentlyClickedEvent: EventModel | null;
  export let clickCallback: (event: EventModel) => void;
</script>

<style lang="scss">
  @import "../../styles/colors.scss";
  @import "../../styles/dimensions.scss";

  div.day {
    min-width: 0;
    overflow: hidden;
  }

  div.background {
    display: flex;
    flex-direction: column;
    gap: $gapSmall;
    margin: calc($gapSmall / 2);
    padding-top: $paddingSmaller;
    border-radius: $borderRadiusSmall;
    background-color: $backgroundSecondary;
    height: calc(100% - $gapSmall);
  }
  div.otherMonth {
    opacity: .5;
  }

  span.date {
    text-align: center;
    width: 100%;
    display: block;
  }
  div.sunday > span.date {
    color: red;
  }

  span.more {
    text-align: center;
    color: $foregroundFaded;
    font-size: $fontSizeSmall;
  }

  div.events {
    display: flex;
    flex-direction: column;
    gap: $gapTiny;
  }
</style>

<div class="day">
  <div class="background" class:otherMonth={!isCurrentMonth} class:sunday={date.getDay() === 0}>
    <span class="date">
      {date.getDate()}
    </span>
    <div class="events">
      {#each events.slice(0,2) as event}
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
      {#if events.length > 2}
        <span class="more">
          and {events.length - 2} more
        </span>
      {/if}
    </div>
  </div>
</div>
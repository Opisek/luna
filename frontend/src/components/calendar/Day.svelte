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

  let containerHeight: number;
  let maxEvents: number = 10;

  // TODO: figure out why sometimes 
  $: ((height: number) => {
    if (height == 0) {
      maxEvents = 0;
      return;
    }
    // TODO: figure out how to extract the proper height (instead of hard-coded 20)
    const slots = Math.floor(height / 22)
    if (events.length > slots) {
      maxEvents = slots - 1;
    } else {
      maxEvents = events.length;
    }
  })(containerHeight);
</script>

<style lang="scss">
  @import "../../styles/colors.scss";
  @import "../../styles/dimensions.scss";

  div.day {
    min-width: 0;
    overflow: hidden;
    height: 100%;
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
    height: 100%;

    // TODO: figure out how to circumvent the css restriction of overflow-y: hidden and overflow-x: visible not being combinable
    //overflow: hidden;
  }
</style>

<div class="day">
  <div class="background" class:otherMonth={!isCurrentMonth} class:sunday={date.getDay() === 0}>
    <span class="date">
      {date.getDate()}
    </span>
    <div class="events" bind:offsetHeight={containerHeight}>
      {#each events.slice(0,maxEvents) as event}
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
      {#if events.length > maxEvents}
        <span class="more">
          and {events.length - maxEvents} more
        </span>
      {/if}
    </div>
  </div>
</div>
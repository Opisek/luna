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

  export let containerHeight: number;
  export let maxEvents: number = 1;
  let actualMaxEvents: number = 1;
  $: actualMaxEvents = maxEvents <= events.length - 1 ? maxEvents - 1 : maxEvents;
</script>

<style lang="scss">
  @import "../../styles/colors.scss";
  @import "../../styles/dimensions.scss";

  div.day {
    min-width: 0;
    overflow: hidden;
    height: 100%;
    position: relative;
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
    margin-right: 1em;
  }

  div.events {
    position: absolute;
    top: calc($gapSmall / 2 + $fontSize + $paddingSmaller + $gapSmall);
    display: flex;
    flex-direction: column;
    gap: $gapTiny;
    height: 100%;
    width: calc(100% + 1em); // +1em needed for long events, otherwise the boundary is visible

    // TODO: figure out how to circumvent the css restriction of overflow-y: hidden and overflow-x: visible not being combinable
    overflow: hidden;
  }
</style>

<div class="day">
  <div class="background" class:otherMonth={!isCurrentMonth} class:sunday={date.getDay() === 0}>
    <span class="date">
      {date.getDate()}
    </span>
  </div>
  {#if isFirstDay}
    <div class="events" bind:offsetHeight={containerHeight}>
      {#each events as event, i}
        <CalendarEvent
          event={event}
          isFirstDay={isFirstDay}
          isLastDay={isLastDay}
          date={date}
          visible={i < actualMaxEvents}
          bind:currentlyHoveredEvent={currentlyHoveredEvent}
          bind:currentlyClickedEvent={currentlyClickedEvent}
          clickCallback={clickCallback}
        />
      {/each}
      {#if events.length > maxEvents}
        <span class="more">
          and {events.length - actualMaxEvents} more
        </span>
      {/if}
    </div>
  {:else}
    <div class="events">
      {#each events as event, i}
        <CalendarEvent
          event={event}
          isFirstDay={isFirstDay}
          isLastDay={isLastDay}
          date={date}
          visible={i < actualMaxEvents}
          bind:currentlyHoveredEvent={currentlyHoveredEvent}
          bind:currentlyClickedEvent={currentlyClickedEvent}
          clickCallback={clickCallback}
        />
      {/each}
      {#if events.length > maxEvents && actualMaxEvents >= 0}
        <span class="more">
          and {events.length - actualMaxEvents} more
        </span>
      {/if}
    </div>
  {/if}
</div>


<!--
TODO: use snippets when svelte 5 is out
{#snippet eventRows}
  {#each events as event, i}
    <CalendarEvent
      event={event}
      isFirstDay={isFirstDay}
      isLastDay={isLastDay}
      date={date}
      visible={i < actualMaxEvents}
      bind:currentlyHoveredEvent={currentlyHoveredEvent}
      bind:currentlyClickedEvent={currentlyClickedEvent}
      clickCallback={clickCallback}
    />
  {/each}
  {#if events.length > maxEvents}
    <span class="more">
      and {events.length - actualMaxEvents} more
    </span>
  {/if}
{/snippet}
-->
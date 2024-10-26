<script lang="ts">
  import EventEntry from "./EventEntry.svelte";
  import IconButton from "../interactive/IconButton.svelte";
  import { PlusIcon } from "lucide-svelte";
  import EventModal from "../modals/EventModal.svelte";

  export let date: Date;

  export let isCurrentMonth: boolean;

  export let isFirstDay: boolean;
  export let isLastDay: boolean;

  export let events: (EventModel | null)[];
  
  export let currentlyHoveredEvent: EventModel | null;
  export let currentlyClickedEvent: EventModel | null;
  export let clickCallback: (event: EventModel) => void;

  let newEvent: EventModel;
  let dummy: () => any;
  let showCreateEventModal: () => any;
  let createEventButtonClick = () => {
    const start = new Date(date);
    start.setHours(12, 0, 0, 0);

    const end = new Date(date);
    end.setHours(13, 0, 0, 0);

    newEvent = {
      id: "",
      calendar: "",
      name: "",
      desc: "",
      color: "",
      date: {
        start: start,
        end: end,
        allDay: false,
      }
    };

    setTimeout(() => {
      showCreateEventModal();
    }, 0);
  };

  export let containerHeight: number;
  export let maxEvents: number = 1;
  let actualMaxEvents: number = 1;
  $: actualMaxEvents = maxEvents <= events.length - 1 ? maxEvents - 1 : maxEvents;
</script>

<style lang="scss">
  @import "../../styles/animations.scss";
  @import "../../styles/colors.scss";
  @import "../../styles/dimensions.scss";
  @import "../../styles/text.scss";

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
    padding: $paddingSmaller;
    border-radius: $borderRadiusSmall;
    background-color: $backgroundSecondary;
    height: calc(100% - $gapSmall);
  }

  div.otherMonth {
    opacity: .5;
  }

  span.top {
    display: grid;
    grid-template-columns: 1fr 1fr 1fr;
    grid-template-areas: "none date add";
  }
  span.date {
    text-align: center;
    width: 100%;
    display: block;
    grid-area: date;
    user-select: none;
  }
  span.sunday {
    color: $foregroundSunday;
  }
  span.add {
    grid-area: add;
    display: flex;
    align-items: center;
    justify-content: right;
    opacity: 0;
    transition: opacity $animationSpeed;
  }
  div.day:hover span.add {
    opacity: 1;
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
    // TODO: z-index so long event names are not truncated
    width: calc(100% + 1em); // +1em needed for long events, otherwise the boundary is visible
    overflow: hidden;
  }
</style>

<div class="day">
  <div class="background" class:otherMonth={!isCurrentMonth}>
    <span class="top">
      <span class="date" class:sunday={date.getDay() === 0}>
        {date.getDate()}
      </span>
      <span class="add">
        <IconButton click={createEventButtonClick} tabindex={-1}>
          <PlusIcon size={13}/>
        </IconButton>
        <EventModal bind:showCreateModal={showCreateEventModal} bind:showModal={dummy} event={newEvent}/>
      </span>
    </span>
  </div>
  {#if isFirstDay}
    <div class="events" bind:offsetHeight={containerHeight}>
      {#each events as event, i}
        <EventEntry
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
  {:else}
    <div class="events">
      {#each events as event, i}
        <EventEntry
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
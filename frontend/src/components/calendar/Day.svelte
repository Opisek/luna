<script lang="ts">
  import { PlusIcon } from "lucide-svelte";

  import EventEntry from "./EventEntry.svelte";
  import EventModal from "../modals/EventModal.svelte";
  import IconButton from "../interactive/IconButton.svelte";

  import { EmptyEvent, NoOp } from '$lib/client/placeholders';

  interface Props {
    date: Date;
    isCurrentMonth: boolean;
    isFirstDay: boolean;
    isLastDay: boolean;
    events: (EventModel | null)[];
    maxEvents?: number;
    containerHeight: number;
    clickCallback: (event: EventModel) => void;
  }

  let {
    date,
    isCurrentMonth,
    isFirstDay,
    isLastDay,
    events,
    maxEvents = 1,
    containerHeight = $bindable(),
    clickCallback
  }: Props = $props();

  let newEvent: EventModel = $state(EmptyEvent);
  let createNewEvent: boolean = $state(false);
  let showCreateEventModal: () => any = $state(NoOp);

  let createEventButtonClick = () => {
    createNewEvent = true;

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

  let actualMaxEvents: number = $derived(maxEvents <= events.length - 1 ? maxEvents - 1 : maxEvents);
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
        {#if createNewEvent}
          <EventModal bind:showCreateModal={showCreateEventModal} bind:event={newEvent}/>
        {/if}
      </span>
    </span>
  </div>
  {#if isFirstDay}
    <div class="events" bind:offsetHeight={containerHeight}>
      {@render eventRows()}
    </div>
  {:else}
    <div class="events">
      {@render eventRows()}
    </div>
  {/if}
</div>

{#snippet eventRows()}
  {#each events as event, i}
    <EventEntry
      event={event}
      isFirstDay={isFirstDay}
      isLastDay={isLastDay}
      date={date}
      visible={i < actualMaxEvents}
      clickCallback={clickCallback}
    />
  {/each}
  {#if events.length > maxEvents && actualMaxEvents >= 0}
    <span class="more">
      and {events.length - actualMaxEvents} more
    </span>
  {/if}
{/snippet}
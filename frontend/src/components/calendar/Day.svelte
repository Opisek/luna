<script lang="ts">
  import { PlusIcon } from "lucide-svelte";
  import { getContext } from "svelte";

  import Event from "./Event.svelte";
  import IconButton from "../interactive/IconButton.svelte";

  import { queueNotification } from "$lib/client/notifications";
  import { NoOp } from "$lib/client/placeholders";

  interface Props {
    date: Date;
    isCurrentMonth: boolean;
    isFirstDay: boolean;
    isLastDay: boolean;
    events: (EventModel | null)[];
    maxEvents?: number;
    containerHeight: number;
    showMore?: (date: Date, events: (EventModel | null)[]) => any;
  }

  let {
    date,
    isCurrentMonth,
    isFirstDay,
    isLastDay,
    events,
    maxEvents = 1,
    containerHeight = $bindable(),
    showMore = NoOp,
  }: Props = $props();

  let showCreateEventModal: ((date: Date) => Promise<EventModel>) = getContext("showNewEventModal");
  let createEventButtonClick = () => {
    showCreateEventModal(date).catch((err) => {
      queueNotification("failure", `Could not create event: ${err.message}`);
    });
  };

  let actualMaxEvents: number = $derived(maxEvents <= events.length - 1 ? maxEvents - 1 : maxEvents);
</script>

<style lang="scss">
  @use "../../styles/animations.scss";
  @use "../../styles/colors.scss";
  @use "../../styles/dimensions.scss";
  @use "../../styles/text.scss";

  div.day {
    min-width: 0;
    overflow: hidden;
    height: 100%;
    position: relative;
  }

  div.background {
    display: flex;
    flex-direction: column;
    gap: dimensions.$gapSmall;
    margin: calc(dimensions.$gapSmall / 2);
    padding: dimensions.$gapSmall;
    border-radius: dimensions.$borderRadiusSmall;
    background-color: colors.$backgroundSecondary;
    height: calc(100% - dimensions.$gapSmall);
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
    color: colors.$foregroundSunday;
  }
  span.add {
    grid-area: add;
    display: flex;
    align-items: center;
    justify-content: right;
    opacity: 0;
    transition: opacity animations.$animationSpeed;
  }
  div.day:hover span.add {
    opacity: 1;
  }

  button.more {
    all: unset;
    text-align: center;
    color: colors.$foregroundDim;
    font-size: text.$fontSizeSmall;
    margin-right: 1em;
    cursor: pointer;
  }

  div.events {
    position: absolute;
    top: calc(text.$fontSize + 2.5 * dimensions.$gapSmall);
    display: flex;
    flex-direction: column;
    gap: dimensions.$gapTiny;
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
      </span>
    </span>
  </div>
  {#if isFirstDay}
    <div class="events" bind:offsetHeight={containerHeight}>
      {@render eventEntries()}
    </div>
  {:else}
    <div class="events">
      {@render eventEntries()}
    </div>
  {/if}
</div>

{#snippet eventEntries()}
  <!-- TODO: forcing EventEntry to be unique for each event and i like that
  fixes a few issues but might be less performant. figure out the right
  compromise -->
  {#each events as event, i ((event?.id || 0) + i.toString())}
    <Event
      event={event}
      isFirstDay={isFirstDay}
      isLastDay={isLastDay}
      date={date}
      visible={i < actualMaxEvents}
    />
  {/each}
  {#if events.length > maxEvents && actualMaxEvents >= 0}
    <button class="more" onclick={() => showMore(date, events)}>
      {#if actualMaxEvents == 0}
       {events.length} events
      {:else}
        and {events.length - actualMaxEvents} more
      {/if}
    </button>
  {/if}
{/snippet}
<script lang="ts">
  import { TextIcon } from "lucide-svelte";

  import { GetEventColor, GetEventHoverColor, GetEventRGB, isDark } from "$lib/common/colors";
  import { passIfEnter } from "$lib/common/inputs";

  import { getContext } from "svelte";
  import type { Writable } from "svelte/store";
  import { NoOp } from "$lib/client/placeholders";

  interface Props {
    visible?: boolean;
    event: EventModel | null;
    isFirstDay: boolean;
    isLastDay: boolean;
    date: Date;
    view: "month" | "week" | "day";
  }

  let {
    visible = true,
    event,
    isFirstDay,
    isLastDay,
    date,
    view
  }: Props = $props();

  let remainingDays = $derived.by(() => {
    // keep in mind start of the week is monday for now
    // Monday: 1, Tuesday: 2, ..., Sunday: 0

    if (!date || !event) return 0;
    if (view === "day") return 1;

    const remainingTime = event.date.end.getTime() - date.getTime();
    const remainingDays = Math.ceil(remainingTime / (1000 * 60 * 60 * 24));

    return remainingDays;
  })

  let remainingDaysThisWeek = $derived.by(() => {
    const myDayIndex = (date.getDay() + 6) % 7;
    const remainingDaysThisWeek = Math.min(remainingDays, 7 - myDayIndex);

    return remainingDaysThisWeek;
  })

  let eventEndsThisWeek = $derived(remainingDays == remainingDaysThisWeek);

  let currentlyHoveredEvent = getContext("currentlyHoveredEvent") as Writable<EventModel | null>;
  let currentlyClickedEvent = getContext("currentlyClickedEvent") as Writable<EventModel | null>;

  let showModal: ((event: EventModel) => Promise<EventModel>) = getContext("showEventModal");

  let element: HTMLDivElement | null = $state(null);

  let isEventStart = $derived(event !== null && event.date.start.getTime() >= date.getTime());
  let isFirstDisplay = $derived(isFirstDay || isEventStart);

  let isBackgroundDark: boolean = $derived(event ? isDark(GetEventRGB(event)) : false);

  function mouseEnter() {
    if (event == null) return;

    $currentlyHoveredEvent = event;
  }
  function mouseLeave() {
    if (event == null) return;

    if ($currentlyHoveredEvent == event)
      $currentlyHoveredEvent = null;
    if ($currentlyClickedEvent == event)
      $currentlyClickedEvent = null;
  }
  function mouseDown() {
    if (event == null) return;

    $currentlyClickedEvent = event;
  }
  function mouseUp() {
    if (event == null) return;

    if ($currentlyClickedEvent == event) {
      $currentlyClickedEvent = null;
      showModal(event).then(newEvent => event = newEvent).catch(NoOp);
      element?.blur();
    }
  }
  function keyPress(e: KeyboardEvent) {
    passIfEnter(e, () => {
      if (event) showModal(event).then(newEvent => event = newEvent).catch(NoOp);
      element?.blur();
    });
  }
</script>

<style lang="scss">
  @use "../../styles/animations.scss";
  @use "../../styles/colors.scss";
  @use "../../styles/dimensions.scss";
  @use "../../styles/text.scss";

  div {
    padding: dimensions.$gapSmaller;
    padding-left: calc(var(--gapBetweenDays) + dimensions.$gapSmaller);
    font-size: text.$fontSizeSmall;
    margin: 0;

    display: flex;
    gap: dimensions.$gapTiny;
    flex-direction: row;
    flex-wrap: nowrap;
    align-items: center;

    user-select: none;
    cursor: pointer;

    white-space: nowrap;
    overflow: visible;

    flex-shrink: 0;

    transition: background-color linear animations.$animationSpeedFast;

    z-index: 10;
  }

  div:focus {
    outline: none;
  }

  div::after {
    content: ".";
    visibility: hidden;
  }
  div.placeholder {
    visibility: hidden;
  }
  div.start {
    border-top-left-radius: dimensions.$borderRadius;
    border-bottom-left-radius: dimensions.$borderRadius;
    margin-left: var(--gapBetweenDays);
    padding-left: dimensions.$gapSmaller;
  }
  div.end {
    border-top-right-radius: dimensions.$borderRadius;
    border-bottom-right-radius: dimensions.$borderRadius;
    margin-right: var(--gapBetweenDays);
  }

  div.hidden {
    display: none;
  }

  div.foregroundBright {
    color: colors.$foregroundBright;
  }
  div.foregroundDark {
    color: colors.$foregroundDark;
  }

  span.name {
    text-overflow: ellipsis;
    overflow: hidden;
    min-width: 0;
    flex-shrink: 1;
  }
  span.time {
    flex-shrink: 0;
    text-align: center;
    font-weight: text.$fontWeightLight;
    font-family: text.$fontFamilyTime;
    font-size: text.$fontSizeSmaller;
  }
  span.icons {
    flex-shrink: 0;
    display: flex;
    align-items: center;
  }
</style>

{#if event && (isFirstDisplay || date.getDay() == 1)}
  <div
    bind:this={element}
    class:start={isEventStart}
    class:end={eventEndsThisWeek}
    class:hover={$currentlyHoveredEvent == event}
    class:active={$currentlyClickedEvent == event}
    class:hidden={!visible}
    class:foregroundBright={isBackgroundDark}
    class:foregroundDark={!isBackgroundDark}
    onmouseenter={mouseEnter}
    onmouseleave={mouseLeave}
    onmousedown={mouseDown}
    onmouseup={mouseUp}
    onfocusin={mouseEnter}
    onfocusout={mouseLeave}
    onkeypress={keyPress}
    role="button"
    tabindex={isFirstDisplay ? 0 : -1}
    style="
      background-color:{$currentlyHoveredEvent == event ? GetEventHoverColor(event) : GetEventColor(event)};
      width: calc({remainingDaysThisWeek * 100}% - {(isEventStart ? 1 : 0) + (eventEndsThisWeek ? 1 : 0)} * var(--gapBetweenDays));
    "
  >
    {#if !event.date.allDay && event.date.start >= date}
      <span class="time">
        {event.date.start.toLocaleTimeString([], {hour: '2-digit', minute: '2-digit'})}
      </span>
    {/if}
    <span class="name">
      {event.name}
    </span>
    {#if (event.desc && event.desc != "")}
      <span class="icons">
        <TextIcon size={12}/>
      </span>
    {/if}
  </div>
{:else}
  <div class="placeholder">
  </div>
{/if}

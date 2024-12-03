<script lang="ts">
  import { TextIcon } from "lucide-svelte";

  import { GetEventColor, GetEventHoverColor, GetEventRGB, isDark } from "$lib/common/colors";
  import { passIfEnter } from "$lib/common/inputs";

  import { getContext } from "svelte";
  import type { Writable } from "svelte/store";

  interface Props {
    visible?: boolean;
    event: EventModel | null;
    isFirstDay: boolean;
    isLastDay: boolean;
    date: Date;
  }

  let {
    visible = true,
    event,
    isFirstDay,
    isLastDay,
    date,
  }: Props = $props();

  let currentlyHoveredEvent = getContext("currentlyHoveredEvent") as Writable<EventModel | null>;
  let currentlyClickedEvent = getContext("currentlyClickedEvent") as Writable<EventModel | null>;

  let showModal: ((event: EventModel) => Promise<EventModel>) = getContext("showEventModal");

  let nextDate: Date = $derived(new Date(date.getFullYear(), date.getMonth(), date.getDate() + 1))
  let element: HTMLDivElement; // TODO: do we really need to make a new element when we just want to bind to something else?

  let isEventStart = $derived(event !== null && event.date.start.getTime() >= date.getTime());
  let isEventEnd = $derived(event !== null && nextDate.getTime() >= event.date.end.getTime());
  let isFirstDisplay = $derived(isFirstDay || isEventStart);
  let isLastDisplay: boolean = $derived(isLastDay || isEventEnd);

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
      showModal(event).then(newEvent => event = newEvent);
      element.blur();
    }
  }
  function keyPress(e: KeyboardEvent) {
    passIfEnter(e, () => {
      if (event) showModal(event).then(newEvent => event = newEvent);
      element.blur();
    });
  }
</script>

<style lang="scss">
  @import "../../styles/animations.scss";
  @import "../../styles/colors.scss";
  @import "../../styles/dimensions.scss";
  @import "../../styles/text.scss";

  div {
    padding: $gapSmaller;
    font-size: $fontSizeSmall;
    margin: 0;

    display: flex;
    gap: $gapTiny;
    flex-direction: row;
    flex-wrap: nowrap;
    align-items: center;

    user-select: none;
    cursor: pointer;

    white-space: nowrap;
    overflow: hidden;

    flex-shrink: 0;

    transition: background-color linear $animationSpeedFast;
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
    border-top-left-radius: $borderRadius;
    border-bottom-left-radius: $borderRadius;
    margin-left: calc($gapSmall / 2);
  }
  div.end {
    border-top-right-radius: $borderRadius;
    border-bottom-right-radius: $borderRadius;
    margin-right: calc($gapSmall / 2 + (1em * $fontSize / $fontSizeSmall));
  }

  div.hidden {
    display: none;
  }

  div.foregroundBright {
    color: $foregroundBright;
  }
  div.foregroundDark {
    color: $foregroundDark;
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
    font-weight: $fontWeightLight;
    font-family: $fontFamilyTime;
    font-size: $fontSizeSmaller;
  }
  span.icons {
    flex-shrink: 0;
    display: flex;
    align-items: center;
  }
</style>

<div
  bind:this={element}
  class:placeholder={!event}
  class:start={isFirstDisplay}
  class:end={isLastDisplay}
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
  style="background-color:{$currentlyHoveredEvent == event ? GetEventHoverColor(event) : GetEventColor(event)}"
>
  {#if event && isFirstDisplay}
    {#if !event.date.allDay}
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
  {/if}
</div>
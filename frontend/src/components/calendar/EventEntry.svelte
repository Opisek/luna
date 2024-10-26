<script lang="ts">
  import { GetEventColor, GetEventHoverColor, GetEventRGB, isDark } from "$lib/common/colors";
  import { TextIcon } from "lucide-svelte";
  import { passIfEnter } from "../../lib/common/inputs";

  export let visible: boolean = true;

  export let event: EventModel | null;
  export let isFirstDay: boolean;
  export let isLastDay: boolean;
  export let date: Date;
  let nextDate: Date;
  $: if (date) {
    nextDate = new Date(date);
    nextDate.setDate(date.getDate() + 1);
  }

  let element: HTMLDivElement;

  export let currentlyHoveredEvent: EventModel | null;
  export let currentlyClickedEvent: EventModel | null;
  export let clickCallback: (event: EventModel) => void;

  let isEventStart, isEventEnd, isFirstDisplay, isLastDisplay: boolean;
  $: isEventStart = event !== null && event.date.start.getTime() >= date.getTime();
  $: isEventEnd = event !== null && nextDate.getTime() >= event.date.end.getTime();
  $: isFirstDisplay = isFirstDay || isEventStart;
  $: isLastDisplay = isLastDay || isEventEnd;

  let isBackgroundDark: boolean;
  $: isBackgroundDark = event ? isDark(GetEventRGB(event)) : false;

  function mouseEnter() {
    if (event == null) return;

    currentlyHoveredEvent = event;
  }
  function mouseLeave() {
    if (event == null) return;

    if (currentlyHoveredEvent == event)
      currentlyHoveredEvent = null;
    if (currentlyClickedEvent == event)
      currentlyClickedEvent = null;
  }
  function mouseDown() {
    if (event == null) return;

    currentlyClickedEvent = event;
  }
  function mouseUp() {
    if (event == null) return;

    if (currentlyClickedEvent == event) {
      currentlyClickedEvent = null;
      clickCallback(event);
      element.blur();
    }
  }
  function keyPress(e: KeyboardEvent) {
    passIfEnter(e, () => {
      if (event) clickCallback(event);
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
    padding: $paddingTiny;
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
  class:hover={currentlyHoveredEvent == event}
  class:active={currentlyClickedEvent == event}
  class:hidden={!visible}
  class:foregroundBright={isBackgroundDark}
  class:foregroundDark={!isBackgroundDark}
  on:mouseenter={mouseEnter}
  on:mouseleave={mouseLeave}
  on:mousedown={mouseDown}
  on:mouseup={mouseUp}
  on:focusin={mouseEnter}
  on:focusout={mouseLeave}
  on:keypress={keyPress}
  role="button"
  tabindex={isFirstDisplay ? 0 : -1}
  style="background-color:{currentlyHoveredEvent == event ? GetEventHoverColor(event) : GetEventColor(event)}"
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
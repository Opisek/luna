<script lang="ts">

  export let event: CalendarEventModel | null;
  export let isFirstDay: boolean;
  export let isLastDay: boolean;
  export let date: Date;

  export let currentlyHoveredEvent: CalendarEventModel | null;
  export let currentlyClickedEvent: CalendarEventModel | null;
  export let clickCallback: (event: CalendarEventModel) => void;

  const previousDate = new Date(date);
  previousDate.setDate(date.getDate() - 1);
  const nextDate = new Date(date);
  nextDate.setDate(date.getDate() + 1);

  const isEventStart = event && previousDate < event.start;
  const isEventEnd = event && nextDate >= event.end;

  const isFirstDisplay = isFirstDay || isEventStart;
  const isLastDisplay = isLastDay || isEventEnd;

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
    }
  
  }
</script>

<style lang="scss">
  @import "../../styles/dimensions.scss";

  div {
    background-color: #cbe6ec;
    padding: $paddingTiny;
    font-size: $fontSizeSmall;
    margin: 0 (-$gap);
  }
  div.hover {
    background-color: #dbecf0;
  }
  div.active {
    transform: scale(1.1);
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
    margin-left: 0;
  }
  div.end {
    border-top-right-radius: $borderRadius;
    border-bottom-right-radius: $borderRadius;
    margin-right: 0;
  }
</style>

<div
  class:placeholder={!event}
  class:start={isFirstDisplay}
  class:end={isLastDisplay}
  class:hover={currentlyHoveredEvent == event}
  class:active={currentlyClickedEvent == event}
  on:mouseenter={mouseEnter}
  on:mouseleave={mouseLeave}
  on:mousedown={mouseDown}
  on:mouseup={mouseUp}
  role="button"
  tabindex="0"
>
  {#if event && isFirstDisplay}
    {event.title}
  {/if}
</div>
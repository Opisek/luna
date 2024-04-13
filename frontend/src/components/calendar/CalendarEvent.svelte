<script lang="ts">
  export let event: CalendarEventModel | null;
  export let isFirstDay: boolean;
  export let isLastDay: boolean;
  export let date: Date;

  const previousDate = new Date(date);
  previousDate.setDate(date.getDate() - 1);
  const nextDate = new Date(date);
  nextDate.setDate(date.getDate() + 1);

  const isEventStart = event && previousDate < event.start;
  const isEventEnd = event && nextDate >= event.end;

  const isFirstDisplay = isFirstDay || isEventStart;
  const isLastDisplay = isLastDay || isEventEnd;
</script>

<style lang="scss">
  @import "../../styles/dimensions.scss";

  div {
    background-color: #cbe6ec;
    padding: $paddingTiny;
    font-size: $fontSizeSmall;
    margin: 0 (-$gap);
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

<div class:placeholder={!event} class:start={isFirstDisplay} class:end={isLastDisplay}>
  {#if event && isFirstDisplay}
    {event.title}
  {/if}
</div>
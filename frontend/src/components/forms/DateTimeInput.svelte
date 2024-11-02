<script lang="ts">
  import { barFocusIndicator } from "../../lib/client/decoration";
  import DateModal from "../modals/DateModal.svelte";
  import TimeModal from "../modals/TimeModal.svelte";
  import Label from "./Label.svelte";

  export let value: Date;
  export let allDay: boolean;
  export let placeholder: string;
  export let name: string;
  export let editable: boolean;

  export let onChange: (value: Date) => void = () => {};

  let dateButton: HTMLButtonElement;
  let timeButton: HTMLButtonElement;

  let showDateModal = () => {};
  let showTimeModal = () => {};

  function dateClick(e: MouseEvent | KeyboardEvent) {
    if (editable) {
      showDateModal();
      if (e.detail !== 0) {
        dateButton.blur();
      }
    }
  }

  function timeClick(e: MouseEvent | KeyboardEvent) {
    if (editable) {
      showTimeModal();
      if (e.detail !== 0) {
        timeButton.blur();
      }
    }
  }
</script>

<style lang="scss">
  @import "../../styles/animations.scss";
  @import "../../styles/colors.scss";
  @import "../../styles/dimensions.scss";
  @import "../../styles/text.scss";

  div.row {
    display: flex;
    flex-direction: row;
    gap: $gapSmall;
    margin: $gapSmall;
  }

  div.editable {
    margin: 0;
  }

  button {
    all: unset;
    border-radius: $borderRadius;
    cursor: text;
    transition: padding $animationSpeedFast linear, border-radius $animationSpeedFast linear;
    padding: $gapSmall;
    margin: -$gapSmall;
    position: relative;
    overflow: hidden;
  }

  div.editable button {
    background: $backgroundSecondary;
    cursor: pointer;
    margin: 0;
  }
</style>

<Label name={name}>{placeholder}</Label>
<div class="row" class:editable={editable}>
  <button
    bind:this={dateButton}
    on:click={dateClick}
    type="button"
    tabindex={editable ? 0 : -1}
    use:barFocusIndicator
  >
    {value.toLocaleDateString()}
  </button>
  {#if !allDay}
    <button
      bind:this={timeButton}
      on:click={timeClick}
      type="button"
      tabindex={editable ? 0 : -1}
      use:barFocusIndicator
    >
      {value.toLocaleTimeString([], {hour: "2-digit", minute: "2-digit"})}
    </button>
  {/if}
</div>

<DateModal bind:date={value} bind:showModal={showDateModal} onChange={onChange}/>
<TimeModal bind:date={value} bind:showModal={showTimeModal} onChange={onChange}/>
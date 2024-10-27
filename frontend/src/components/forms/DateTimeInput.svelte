<script lang="ts">
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
  }

  div.editable button {
    background: $backgroundSecondary;
    cursor: pointer;
    margin: 0;
  }

  div.editable > div.wrapper {
    background-color: $backgroundAccent;
    border-radius: calc($borderRadius + 0.1em);
    transition: padding $animationSpeedFast linear;
  } 
  div.editable > div.wrapper:focus-within {
    padding-left: $borderActiveWidth;
    border-top-left-radius: $borderRadius;
    border-bottom-left-radius: $borderRadius;
  }
  div.editable > div.wrapper:focus-within > button {
    padding-left: calc($gapSmall - $borderActiveWidth);
    border-top-left-radius: 0;
    border-bottom-left-radius: 0;
  }
</style>

<Label name={name}>{placeholder}</Label>
<div class="row" class:editable={editable}>
  <div class="wrapper">
    <button
      bind:this={dateButton}
      on:click={dateClick}
      type="button"
    >
      {value.toLocaleDateString()}
    </button>
  </div>
  {#if !allDay}
    <div class="wrapper">
      <button
        bind:this={timeButton}
        on:click={timeClick}
        type="button"
      >
        {value.toLocaleTimeString([], {hour: "2-digit", minute: "2-digit"})}
      </button>
    </div>
  {/if}
</div>

<DateModal bind:date={value} bind:showModal={showDateModal} onChange={onChange}/>
<TimeModal bind:date={value} bind:showModal={showTimeModal} onChange={onChange}/>
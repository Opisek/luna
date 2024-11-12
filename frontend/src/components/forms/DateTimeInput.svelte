<script lang="ts">
  import DateModal from "../modals/DateModal.svelte";
  import Label from "./Label.svelte";
  import TimeModal from "../modals/TimeModal.svelte";

  import { NoOp } from "$lib/client/placeholders";
  import { focusIndicator } from "$lib/client/decoration";

  interface Props {
    value: Date;
    allDay: boolean;
    placeholder: string;
    name: string;
    editable: boolean;
    onChange?: (value: Date) => void;
  }

  let {
    value = $bindable(),
    allDay,
    placeholder,
    name,
    editable,
    onChange = NoOp
  }: Props = $props();

  let dateButton: HTMLButtonElement;
  let timeButton: HTMLButtonElement;

  let showDateModal = $state(NoOp);
  let showTimeModal = $state(NoOp);

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
    onclick={dateClick}
    type="button"
    tabindex={editable ? 0 : -1}
    use:focusIndicator
  >
    {value.toLocaleDateString()}
  </button>
  {#if !allDay}
    <button
      bind:this={timeButton}
      onclick={timeClick}
      type="button"
      tabindex={editable ? 0 : -1}
      use:focusIndicator
    >
      {value.toLocaleTimeString([], {hour: "2-digit", minute: "2-digit"})}
    </button>
  {/if}
</div>

<DateModal bind:date={value} bind:showModal={showDateModal} onChange={onChange}/>
<TimeModal bind:date={value} bind:showModal={showTimeModal} onChange={onChange}/>
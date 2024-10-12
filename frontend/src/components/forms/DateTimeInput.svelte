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

  let showDateModal = () => {};
  let showTimeModal = () => {};

  function dateClick() {
    if (editable) {
      showDateModal();
    }
  }

  function timeClick() {
    if (editable) {
      showTimeModal();
    }
  }
</script>

<style lang="scss">
  @import "../../styles/colors.scss";
  @import "../../styles/dimensions.scss";

  div {
    display: flex;
    flex-direction: row;
    gap: $gapSmall;
    margin: $gapSmall;
  }

  div.editable {
    margin: 0;
  }

  div > button {
    all: unset;
    border-radius: $borderRadius;
    cursor: text;
  }

  div.editable > button {
    padding: $gapSmall;
    background: $backgroundSecondary;
    cursor: pointer;
  }
</style>

<Label name={name}>{placeholder}</Label>
<div class:editable={editable}>
  <button on:click={dateClick}>{value.toLocaleDateString()}</button>
  {#if !allDay}
    <button on:click={timeClick}>{value.toLocaleTimeString([], {hour: "2-digit", minute: "2-digit"})}</button>
  {/if}
</div>

<DateModal bind:date={value} bind:showModal={showDateModal} onChange={onChange}/>
<TimeModal bind:date={value} bind:showModal={showTimeModal} onChange={onChange}/>
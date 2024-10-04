<script lang="ts">
  import DateModal from "../modals/DateModal.svelte";
  import Label from "./Label.svelte";

  export let value: Date;
  export let placeholder: string;
  export let name: string;
  export let editable: boolean;

  let showModal = () => {};
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

  div > span {
    border-radius: $borderRadius;
    cursor: pointer;
  }

  div.editable > span {
    padding: $gapSmall;
    background: $backgroundSecondary;
  }
</style>

<Label name={name}>{placeholder}</Label>
<div class:editable={editable}>
  <span on:click={showModal}>{value.toLocaleDateString()}</span>
  <span>{value.toLocaleTimeString([], {hour: "2-digit", minute: "2-digit"})}</span>
</div>

<DateModal bind:date={value} bind:showModal={showModal}/>
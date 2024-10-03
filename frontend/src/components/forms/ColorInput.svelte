<script lang="ts">
  import IconButton from "../interactive/IconButton.svelte";
  import ColorCircle from "../misc/ColorCircle.svelte";
  import ColorModal from "../modals/ColorModal.svelte";
  import Label from "./Label.svelte";

  export let color: string;
  export let name: string;
  export let editable: boolean;

  let showModal: () => any;
</script>

<style lang="scss">
  @import "../../styles/dimensions.scss";

  div {
    padding: $gapSmall;
  }
  div.editable {
    padding: $gapSmaller;
  }
</style>

<Label name={name}>Color</Label>
<div
  class:editable={editable}
>
  {#if editable}
    <IconButton click={showModal}>
      <ColorCircle color={color} size="medium"/>
    </IconButton>
    <ColorModal
      bind:showModal={showModal} 
      bind:color={color}
    />
  {:else}
    <ColorCircle color={color} size="medium"/>
  {/if}
</div>

<!-- TODO: svelte snippet once out -->
<script lang="ts">
  import ColorCircle from "../misc/ColorCircle.svelte";
  import ColorModal from "../modals/ColorModal.svelte";
  import IconButton from "../interactive/IconButton.svelte";
  import Label from "./Label.svelte";

  import { NoOp } from "$lib/client/placeholders";

  interface Props {
    color: string;
    name: string;
    editable: boolean;
  }

  let { color = $bindable(), name, editable }: Props = $props();

  let showModal: () => any = $state(NoOp);
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
      {@render circle()}
    </IconButton>
    <ColorModal
      bind:showModal={showModal} 
      bind:color={color}
    />
  {:else}
    {@render circle()}
  {/if}
</div>

{#snippet circle()}
  <ColorCircle color={color} size="medium"/>
{/snippet}
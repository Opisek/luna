<script lang="ts">
  import ColorCircle from "../misc/ColorCircle.svelte";
  import ColorModal from "../modals/ColorModal.svelte";
  import IconButton from "../interactive/IconButton.svelte";
  import Label from "./Label.svelte";
  import { NoOp } from "$lib/scripts/client/placeholders";
  import { t } from "@sveltia/i18n";

  interface Props {
    color: string;
    name: string;
    editable: boolean;
  }

  let { color = $bindable(), name, editable }: Props = $props();

  let showModal: (initial: string | null, anchor?: HTMLElement) => Promise<string> = $state(Promise.reject);
  let anchor: HTMLElement | undefined = $state();

  async function pickColor() {
    await showModal(color, anchor).then((pickedColor) => color = pickedColor).catch(NoOp);
  }
</script>

<style lang="scss">
  @use "$lib/styles/dimensions.scss";

  div {
    padding: dimensions.$gapSmall;
    margin: -(dimensions.$gapSmaller) 0;
  }
  div.editable {
    padding: dimensions.$gapSmaller;
  }
</style>

<Label name={name}>{t("color.display")}</Label>
<div
  class:editable={editable}
>
  {#if editable}
    <IconButton bind:button={anchor} onClick={pickColor} alt={t("color.display")}>
      {@render circle()}
    </IconButton>
    <ColorModal
      bind:showModal={showModal} 
    />
  {:else}
    {@render circle()}
  {/if}
</div>

{#snippet circle()}
  <ColorCircle color={color} size="medium"/>
{/snippet}
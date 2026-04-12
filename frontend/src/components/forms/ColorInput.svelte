<script lang="ts">
  import ColorCircle from "../misc/ColorCircle.svelte";
  import ColorModal from "../modals/ColorModal.svelte";
  import IconButton from "../interactive/IconButton.svelte";
  import Label from "./Label.svelte";
  import { NoOp } from "../../lib/client/placeholders";
  import { t } from "@sveltia/i18n";

  interface Props {
    color: string;
    name: string;
    editable: boolean;
  }

  let { color = $bindable(), name, editable }: Props = $props();

  let showModal: () => Promise<string> = $state(Promise.reject);

  async function pickColor() {
    await showModal().then((pickedColor) => color = pickedColor).catch(NoOp);
  }
</script>

<style lang="scss">
  @use "../../styles/dimensions.scss";

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
    <IconButton onClick={pickColor} alt={t("color.display")}>
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
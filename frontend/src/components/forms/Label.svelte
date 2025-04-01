<script lang="ts">
  import type { Snippet } from "svelte";
  import Tooltip from "../interactive/Tooltip.svelte";

  interface Props {
    name: string;
    info?: string;
    ownPositioning?: boolean;
    children?: Snippet;
  }

  let {
    name,
    info = "",
    ownPositioning = true,
    children
  }: Props = $props();
</script>

<style lang="scss">
  @use "../../styles/colors.scss";
  @use "../../styles/dimensions.scss";
  @use "../../styles/text.scss";

  label {
    color: colors.$foregroundDim;
    font-size: text.$fontSizeSmall;
    cursor: text;
    display: flex;
    flex-direction: row;
    gap: dimensions.$gapSmaller;
    align-items: center;
  }

  .ownPositioning {
    margin-bottom: -(dimensions.$gapMiddle);
    padding-left: calc(dimensions.$gapSmall * (text.$fontSize / text.$fontSizeSmall));
  }
</style>

<label for={name} tabindex="-1" class:ownPositioning={ownPositioning}>
  {@render children?.()}
  {#if info}
    <Tooltip tight tiny>
      {info}
    </Tooltip>
  {/if}
</label>
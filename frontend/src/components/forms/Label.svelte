<script lang="ts">
  import type { Snippet } from "svelte";
  import Tooltip from "../interactive/Tooltip.svelte";
  import { extendUniqueElementId } from "$lib/common/dom";

  interface Props {
    describes: string;
    info?: string;
    ownPositioning?: boolean;
    children?: Snippet;
  }

  let {
    describes,
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
    color: color-mix(in srgb, colors.$foregroundPrimary 50%, transparent);
    font-size: text.$fontSizeSmall;
    cursor: text;
    display: flex;
    flex-direction: row;
    gap: dimensions.$gapSmaller;
    align-items: center;
    padding-left: calc(dimensions.$gapSmall / text.$fontSizeSmallRatio);
  }

  .ownPositioning {
    margin-bottom: -(dimensions.$gapMiddle);
  }
</style>

<label for={describes} id={extendUniqueElementId(["label"], describes)} tabindex="-1" class:ownPositioning={ownPositioning}>
  {@render children?.()}
  {#if info !== undefined && info.trim().length != 0}
    <Tooltip tight tiny describes={describes}>
      {info}
    </Tooltip>
  {/if}
</label>
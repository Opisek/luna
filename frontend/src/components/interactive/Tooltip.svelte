<script lang="ts">
  import { CircleAlert, Info } from "lucide-svelte";

  import TooltipPopup from "../popups/TooltipPopup.svelte";
  import type { Snippet } from "svelte";

  interface Props {
    error?: boolean;
    children?: Snippet;
    tight?: boolean;
    tiny?: boolean;
  }

  let {
    error = false,
    tight = false,
    tiny = false,
    children
  }: Props = $props();
</script>

<style lang="scss">
  @use "../../styles/colors.scss";
  @use "../../styles/dimensions.scss";

  div {
    position: relative;
    color: color-mix(in srgb, colors.$foregroundPrimary 50%, transparent);
    cursor: help;
    display: flex;
    justify-content: center;
    outline: 0;
    padding: dimensions.$gapSmaller;
  }

  div.error {
    color: colors.$backgroundFailure;
  }

  div.tight {
    padding: 0
  }
</style>

<!-- svelte-ignore a11y_no_noninteractive_tabindex -->
<div
  class:error={error}
  class:tight={tight}
  role="tooltip"
  tabindex="0"
>
  {#if error}
    <CircleAlert size={tiny ? 14 : 16}/>
  {:else}
    <Info size={tiny ? 14 : 16}/>
  {/if}

  <TooltipPopup>
    {@render children?.()}
  </TooltipPopup>
</div>
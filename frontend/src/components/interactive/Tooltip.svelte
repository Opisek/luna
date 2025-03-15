<script lang="ts">
  import { CircleAlert, Info } from "lucide-svelte";

  import TooltipPopup from "../popups/TooltipPopup.svelte";

  interface Props {
    msg: string;
    error?: boolean;
  }

  let {
    msg,
    error = false
  }: Props = $props();
</script>

<style lang="scss">
  @use "../../styles/colors.scss";
  @use "../../styles/dimensions.scss";

  div {
    position: relative;
    color: colors.$foregroundDim;
    cursor: help;
    display: flex;
    justify-content: center;
    padding: dimensions.$gapSmaller;
    outline: 0;
  }

  div.error {
    color: colors.$backgroundFailure;
  }
</style>

<!-- svelte-ignore a11y_no_noninteractive_tabindex -->
<div
  class:error={error}
  role="tooltip"
  tabindex="0"
>
  {#if error}
    <CircleAlert size={16}/>
  {:else}
    <Info size={16}/>
  {/if}

  <TooltipPopup>
    {msg}
  </TooltipPopup>
</div>
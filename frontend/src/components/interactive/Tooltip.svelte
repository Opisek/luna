<script lang="ts">
  import { CircleAlert, Info } from "lucide-svelte";

  import TooltipPopup from "../layout/TooltipPopup.svelte";

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
  @import "../../styles/colors.scss";
  @import "../../styles/dimensions.scss";

  div {
    position: relative;
    color: $foregroundFaded;
    cursor: help;
    display: flex;
    justify-content: center;
    padding: $gapSmaller;
    outline: 0;
  }

  div.error {
    color: $backgroundFailure;
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
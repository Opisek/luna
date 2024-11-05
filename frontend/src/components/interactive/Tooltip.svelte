<script lang="ts">
  import { CircleAlert, Info } from "lucide-svelte";
  import Popup from "../layout/TooltipPopup.svelte";
  import TooltipPopup from "../layout/TooltipPopup.svelte";

  export let msg: string;
  export let error: boolean = false;
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
    padding: $paddingTiny;
    outline: 0;
  }
  div.error {
    color: $backgroundFailure;
  }

</style>

<!-- svelte-ignore a11y-no-noninteractive-tabindex -->
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
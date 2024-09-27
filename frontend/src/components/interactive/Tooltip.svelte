<script lang="ts">
  import { CircleAlert, Info } from "lucide-svelte";

  export let msg: string;
  export let error: boolean = false;

  let bottom: boolean = false;
  let center: boolean = false;
  let right: boolean = false;

  let tooltip: HTMLElement;

  function checkPosition() {
    const rect = tooltip.getBoundingClientRect();

    const x = rect.left + (rect.right - rect.left) / 2;
    const y = rect.top + (rect.bottom - rect.top) / 2;

    bottom = y > window.innerHeight / 2;
    right = x > window.innerWidth / 2;
    center = x < window.innerWidth / 3 * 2 && x > window.innerWidth / 3;
  }
</script>

<style lang="scss">
  @import "../../styles/animations.scss";
  @import "../../styles/colors.scss";
  @import "../../styles/decoration.scss";
  @import "../../styles/dimensions.scss";

  div {
    position: relative;
    color: $foregroundFaded;
    cursor: help;
  }
  div.error {
    color: $backgroundFailure;
  }
  span {
    opacity: 0;
    border-radius: $borderRadius;
    padding: $paddingSmaller;
    background-color: $backgroundPrimary;
    color: $foregroundPrimary;
    transition: opacity $animationSpeed $cubic;
    position: absolute;
    z-index: 10;
    pointer-events: none;
    min-width: none;
    width: max-content;
    max-width: 80vw;
    box-shadow: $boxShadow;
    left: 0;
    white-space: nowrap;
  }
  span.over {
    bottom: calc(100% + $gapTiny);
  }
  span.below {
    top: calc(100% + $gapTiny);
  }
  span.left {
    right: 0;
  }
  span.right {
    left: 0;
  }
  span.center {
    left: 50%;
    transform: translateX(-50%);
  }
  div:hover span {
    opacity: 1;
  }
</style>

<div
  bind:this={tooltip}
  on:mouseenter={checkPosition}
  class:error={error}
  role="tooltip"
>
  {#if error}
    <CircleAlert size={16}/>
  {:else}
    <Info size={16}/>
  {/if}
  <span
    class:over={bottom}
    class:below={!bottom}
    class:left={!right && !center}
    class:right={right && !center}
    class:center={center}
  >
    {msg}
  </span>
</div>
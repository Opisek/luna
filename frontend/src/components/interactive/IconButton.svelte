<script lang="ts">
  import type { Snippet } from "svelte";

  import { NoOp } from "$lib/client/placeholders";

  interface Props {
    up?: () => void;
    down?: () => void;
    click?: () => void;
    visible?: boolean;
    style?: string;
    tabindex?: number;
    children?: Snippet;
  }

  let {
    up = NoOp,
    down = NoOp,
    click = NoOp,
    visible = true,
    style = "",
    tabindex = 0,
    children
  }: Props = $props();

  let button: HTMLElement;

  function clickInternal(e: MouseEvent) {
    e.stopPropagation();
    click();
  }

  function leaveInternal(e: MouseEvent) {
    button.blur();
    up();
  }
  function upInternal(e: MouseEvent) {
    button.blur();
    up();
  }
</script>

<style lang="scss">
  @import "../../styles/animations.scss";
  @import "../../styles/colors.scss";
  @import "../../styles/dimensions.scss";

  button {
    all: unset;
    border-radius: 50%;
    display: flex;
    align-items: center;
    padding: $gapSmaller;
    cursor: pointer;
    position: relative;
    transition: all $cubic $animationSpeed;
  }

  button.hidden {
    visibility: hidden;
  }

  div.circle {
    position: absolute;
    background-color: $backgroundSecondary;
    z-index: -1;
    border-radius: 50%;
    left: 50%;
    top: 50%;
    width: 0%;
    height: 0%;
    transition: all $cubic $animationSpeed;
  }

  button:hover div.circle, button:focus div.circle {
    left: 0;
    top: 0;
    width: 100%;
    height: 100%;
  }

  button:active div.circle {
    width: 125%;
    height: 125%;
    left: -12.5%;
    top: -12.5%;
  }
</style>

<button
  bind:this={button}
  onclick={clickInternal}
  onmousedown={down}
  onmouseleave={leaveInternal}
  onmouseup={upInternal}
  class:hidden={!visible}
  type="button"
  style={style}
  tabindex="{tabindex}"
>
  <div class="circle"></div>
  {@render children?.()}
</button>
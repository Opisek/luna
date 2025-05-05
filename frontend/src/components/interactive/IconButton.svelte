<script lang="ts">
  import type { Snippet } from "svelte";

  import { NoOp } from "$lib/client/placeholders";
  import Tooltip from "./Tooltip.svelte";

  interface Props {
    up?: () => void;
    down?: () => void;
    click?: () => void;
    visible?: boolean;
    info?: string;
    style?: string;
    tabindex?: number;
    href?: string;
    children?: Snippet;
  }

  let {
    up = NoOp,
    down = NoOp,
    click = NoOp,
    visible = true,
    info = "",
    style = "",
    tabindex = 0,
    href = "",
    children
  }: Props = $props();

  // svelte-ignore non_reactive_update
  // isLink is set once and never changed
  let button = $state<HTMLElement | null>(null);

  function clickInternal(e: MouseEvent) {
    e.stopPropagation();
    click();
  }

  function leaveInternal() {
    if (!button) return;
    button.blur();
    up();
  }
  function upInternal() {
    if (!button) return;
    button.blur();
    up();
  }
</script>

<style lang="scss">
  @use "../../styles/animations.scss";
  @use "../../styles/colors.scss";
  @use "../../styles/dimensions.scss";

  button, a {
    all: unset;
    border-radius: 50%;
    display: flex;
    align-items: center;
    padding: dimensions.$gapSmaller;
    cursor: pointer;
    position: relative;
    transition: all animations.$cubic animations.$animationSpeed;
  }

  button.hidden, a.hidden {
    visibility: hidden;
  }

  div.circle {
    position: absolute;
    background-color: colors.$backgroundSecondary;
    z-index: -1;
    border-radius: 50%;
    left: 50%;
    top: 50%;
    width: 0%;
    height: 0%;
    transition: all animations.$cubic animations.$animationSpeed;
    pointer-events: none;
  }

  button:hover div.circle, button:focus div.circle, a:hover div.circle, a:focus div.circle {
    left: 0;
    top: 0;
    width: 100%;
    height: 100%;
  }

  button:active div.circle, a:active div.circle {
    width: 125%;
    height: 125%;
    left: -12.5%;
    top: -12.5%;
  }
</style>

{#if info == ""}
  {@render buttonSnippet()}
{:else}
  <Tooltip
    icon={buttonSnippet} 
    inheritColor={true}
    tight={true}
    pointerCursor={true}
  >
    {info}
  </Tooltip>
{/if}

{#snippet buttonSnippet()}
  {#if href !== ""}
    <a
      bind:this={button}
      class:hidden={!visible}
      href={href}
      style={style}
      tabindex="{tabindex}"
    >
      <div class="circle"></div>
      {@render children?.()}
    </a>
  {:else}
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
  {/if}
{/snippet}
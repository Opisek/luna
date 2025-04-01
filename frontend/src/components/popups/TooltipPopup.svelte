<script lang="ts">
  import type { Snippet } from 'svelte';

  import { browser } from "$app/environment";

  import { calculateOptimalPopupPosition } from "$lib/common/calculations";

  interface Props {
    children?: Snippet;
  }

  let { children }: Props = $props();

  let popup: HTMLElement;

  let bottom: boolean = $state(false);
  let center: boolean = $state(false);
  let right: boolean = $state(false);

  function setupListener(popup: HTMLElement) {
    if (!popup || !popup.parentElement) return;

    popup.parentElement.addEventListener("mouseenter", checkPosition);
  }

  function checkPosition() {
    if (!popup || !popup.parentElement || !browser) return;

    const res = calculateOptimalPopupPosition(popup.parentElement, 3);

    bottom = res.bottom;
    right = res.right;
    center = res.center;
  }

  $effect(() => {
    setupListener(popup);
  });
</script>

<style lang="scss">
  @use "../../styles/animations.scss";
  @use "../../styles/colors.scss";
  @use "../../styles/decorations.scss";
  @use "../../styles/dimensions.scss";

  span.popup {
    position: absolute;
    z-index: 50;
    pointer-events: none;
    left: 50%;
    top: 0;
  }
  span.left.below {
    transform: translateY(75%);
  }
  span.left.over {
    transform: translateY(-100%);
  }
  span.right.below {
    transform: translate(-100%, 75%);
  }
  span.right.over {
    transform: translate(-100%, -100%);
  }
  span.center.below {
    transform: translate(-50%, 75%);
  }
  span.center.over {
    transform: translate(-50%, -100%);
  }
  span.contents {
    position: fixed;
    opacity: 0;
    border-radius: dimensions.$borderRadius;
    padding: dimensions.$gapSmall;
    background-color: colors.$backgroundSecondary;
    color: colors.$foregroundSecondary;
    transition: opacity animations.$animationSpeed animations.$cubic;
    z-index: 10;
    pointer-events: none;
    min-width: none;
    width: fit-content;
    max-width: 30em;
    box-shadow: decorations.$boxShadow;
  }
  :global(*:hover) > span.popup > span.contents,
  :global(*:focus-within) > span.popup > span.contents {
    opacity: 1;
  }
</style>

<span
  bind:this={popup}
  class="popup"
>
  <span
    class="contents"
    class:over={bottom}
    class:below={!bottom}
    class:left={!right && !center}
    class:right={right && !center}
    class:center={center}
  >
    {@render children?.()}
  </span>
</span>
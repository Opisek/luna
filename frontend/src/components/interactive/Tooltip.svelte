<script lang="ts">
  import { CircleAlert, Info } from "lucide-svelte";

  import type { Snippet } from "svelte";
  import { calculateOptimalPopupPosition } from "../../lib/common/calculations";

  interface Props {
    error?: boolean;
    children?: Snippet;
    icon?: Snippet;
    tight?: boolean;
    tiny?: boolean;
    inline?: boolean;
    inheritColor?: boolean;
    pointerCursor?: boolean;
  }

  let {
    error = false,
    tight = false,
    tiny = false,
    inline = false,
    inheritColor = false,
    pointerCursor = false,
    children,
    icon,
  }: Props = $props();

  let iconElement = $state<HTMLElement>();
  let popover = $state<HTMLElement>();
  let popoverTop = $state(0);
  let popoverLeft = $state(0);
  let popoverOpen = $state(false);

  function show() {
    if (!popover || !iconElement) return;
    popover.showPopover();
    popoverOpen = true;

    const popoverRect = popover.getBoundingClientRect();
    const iconRect = iconElement.getBoundingClientRect();

    const marginSize = popoverRect.y - popoverTop;

    const optimalPosition = calculateOptimalPopupPosition(iconElement, 3);

    if (!optimalPosition.center) {
      // The popover is not centered vertically with respect to the icon
      popoverTop = iconRect.top + (iconRect.height - popoverRect.height) / 2 - marginSize;
    } else if (optimalPosition.bottom) {
      // The bottom edge of the popover is above the top edge of the icon
      popoverTop = iconRect.top - popoverRect.height - 2 * marginSize;
    } else {
      // The top edge of the popover is below the bottom edge of the icon
      popoverTop = iconRect.bottom;
    }

    if (optimalPosition.center) {
      // The popover is centered horizontally with respect to the icon
      popoverLeft = iconRect.left + (iconRect.width - popoverRect.width) / 2 - marginSize;
    } else if (optimalPosition.right) {
      // The right edge of the popover is to the left of the left edge of the icon
      popoverLeft = iconRect.left - popoverRect.width - 2 * marginSize;
    } else {
      // The left edge of the popover is to the right of the right edge of the icon
      popoverLeft = iconRect.right;
    }
  }

  function hide() {
    if (!popover) return;
    popoverOpen = false;
  }

  function transitionEnd() {
    if (popover && !popoverOpen && popover.matches(":popover-open")) popover.hidePopover();
  }
</script>

<style lang="scss">
  @use "../../styles/animations.scss";
  @use "../../styles/colors.scss";
  @use "../../styles/dimensions.scss";
  @use "../../styles/decorations.scss";
  @use "../../styles/text.scss";

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

  div.inline {
    display: inline-flex;
    vertical-align: bottom;
    padding: 0;
    margin-bottom: 0.25ch;
  }

  div.inheritColor {
    color: inherit;
  }

  div.pointerCursor {
    cursor: pointer;
  }

  .popover {
    border-radius: dimensions.$borderRadius;
    padding: dimensions.$gapSmall;
    box-shadow: decorations.$boxShadow;

    background-color: colors.$backgroundSecondary;
    color: colors.$foregroundSecondary;

    pointer-events: none;

    min-width: none;
    width: fit-content;
    max-width: 30em;

    font-size: text.$fontSize;
    white-space: pre-wrap;

    outline: 0;
    border: 0;
    margin: dimensions.$gapSmaller;
    box-sizing: border-box;

    opacity: 0;
    transition: opacity animations.$animationSpeed;
  }

  :global(html[data-frost="true"]) .popover {
    background-color: color-mix(in srgb, colors.$backgroundSecondary 65%, transparent);
    backdrop-filter: blur(dimensions.$blurLarge);
  }

  .popover.visible {
    opacity: 1;
  }
</style>

<svelte:window
  onresize={hide}
/>

<!-- svelte-ignore a11y_no_noninteractive_tabindex -->
<div
  class:error={error}
  class:tight={tight}
  class:inline={inline}
  class:inheritColor={inheritColor}
  class:pointerCursor={pointerCursor}
  role="tooltip"
  tabindex="0"
  onmouseenter={show}
  onmouseleave={hide}
  onfocus={show}
  onblur={hide}
  bind:this={iconElement}
>
  {#if icon}
    {@render icon?.()}
  {:else if error}
    <CircleAlert size={tiny ? 14 : 16}/>
  {:else}
    <Info size={tiny ? 14 : 16}/>
  {/if}

  <span
    class="popover" 
    class:visible={popoverOpen}
    bind:this={popover}
    popover="manual"
    style="top: {popoverTop}px; left: {popoverLeft}px"
    ontransitionend={transitionEnd}
  >
    {@render children?.()}
  </span>
</div>
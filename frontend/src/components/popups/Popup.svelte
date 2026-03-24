<script lang="ts">
  import type { Snippet } from "svelte";
  import { NoOp } from "../../lib/client/placeholders";

  interface Props {
    tooltip?: boolean;
    delayed?: boolean;
    children?: Snippet;
    showPopup?: () => void;
    hidePopup?: () => void;
  }

  let {
    tooltip = true,
    delayed = false,
    children,
    showPopup = $bindable(),
    hidePopup = $bindable(NoOp),
  }: Props = $props();

  let visible = $state(false);
  let anchorName = Math.floor(Math.random() * 100000000).toString();
  let popover: (HTMLElement | undefined) = $state();

  $effect(() => {
    if (!popover || !popover.parentElement) return;
    Object.assign(popover.parentElement.style, {
      "anchor-name": `--anchor${anchorName}`,
    });
    if (tooltip) popover.parentElement.setAttribute("aria-describedby", `tooltip${anchorName}`);
  })

  let openTimeout = $state<ReturnType<typeof setTimeout>>();
  showPopup = () => {
    clearTimeout(openTimeout);
    openTimeout = setTimeout(() => {
      if (!popover || popover.matches(":popover-open")) return;
      visible = true;
      popover.showPopover();
      if (!tooltip) {
        setTimeout(() => {
          if (!popover) return;
          popover.focus();
        }, 0);
      }
    }, delayed ? 1000 : 0);
  }

  hidePopup = () => {
    clearTimeout(openTimeout);
    if (!popover) return;
    visible = false;
  }

  function transitionEnd() {
    if (!popover || visible || !popover.matches(":popover-open")) return;
    popover.hidePopover();
  }
</script>

<style lang="scss">
  @use "../../styles/animations.scss";
  @use "../../styles/colors.scss";
  @use "../../styles/decorations.scss";
  @use "../../styles/dimensions.scss";
  @use "../../styles/text.scss";
  
  .popup {
    inset: unset;

    border: 0;
    padding: dimensions.$gapSmall;
    border-radius: dimensions.$borderRadius;
    max-width: 30vw;
    box-shadow: decorations.$boxShadow;
    font-size: text.$fontSize;
    background-color: colors.$backgroundSecondary;
    color: colors.$foregroundSecondary;

    position: fixed;
    position-area: top;
    position-try-fallbacks: bottom, right, left;
    position-try-order: most-width;
    container-type: anchored;

    margin: dimensions.$gapSmall;

    opacity: 0;
    transition: opacity animations.$animationSpeed;

    &.visible:popover-open {
      opacity: 1;

      @starting-style {
        & {
          opacity: 0;
        }
      }
    }
  }

  :global(html[data-frost="true"]) .popup {
    background-color: color-mix(in srgb, colors.$backgroundSecondary 50%, transparent) !important;
    backdrop-filter: blur(dimensions.$blurLarge);
  }

  .popup:focus {
    outline: 0;
  }
</style>

<!-- The typecast to "auto" is because the linter does not yet know about "hint" -->
<div
  bind:this={popover}
  class="popup"
  popover={(tooltip ? "hint" : "auto") as "auto"}
  style={`position-anchor: --anchor${anchorName};`}
  id={`tooltip${anchorName}`}
  class:visible={visible}
  tabindex="-1"
  ontransitionend={transitionEnd}
  role={tooltip ? "tooltip" : "dialog"}
>
  {@render children?.()}
</div>
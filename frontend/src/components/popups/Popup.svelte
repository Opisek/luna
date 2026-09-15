<script lang="ts">
  import { type Snippet } from "svelte";
  import { NoOp } from "../../lib/client/placeholders";
  import { extendUniqueElementId, generateUniqueElementId } from "$lib/common/dom";

  interface Props {
    tooltip?: boolean;
    delayed?: boolean;
    anchor?: HTMLElement | undefined;
    matchWidth?: boolean;
    triangle?: boolean;
    dialog?: boolean;
    visible?: boolean;
    describes?: string;
    labelled?: string;
    children?: Snippet;
    showPopup?: () => Promise<void>;
    hidePopup?: () => void;
  }

  let {
    tooltip = true,
    delayed = false,
    anchor = undefined,
    matchWidth = false,
    triangle = true,
    dialog = false,
    visible = $bindable(false),
    describes,
    labelled,
    children,
    showPopup = $bindable(),
    hidePopup = $bindable(NoOp),
  }: Props = $props();
  let uniqueId = $props.id();

  let popover: (HTMLElement | undefined) = $state();
  let anchorElement = $derived(anchor || (!popover ? undefined : popover.parentElement))
  let anchorName = $state<string | undefined>(undefined);
  let tooltipId = $derived(extendUniqueElementId(tooltip ? ["tooltip"] : ["popup"], describes ?? anchorName));

  let promiseResolve: () => void = $state(NoOp);
  let promiseReject: (reason?: any) => void = $state(NoOp);

  $effect(() => {
    if (!visible || !anchorElement) return;
    // @ts-ignore
    const currentAnchor = anchorElement.style["anchor-name"] as string;
    if (currentAnchor.startsWith("--anchor-") && !currentAnchor.includes("undefined")) {
      anchorName = currentAnchor.substring(9);
    } else {
      anchorName = anchorElement.id || uniqueId;
      Object.assign(anchorElement.style, {
        "anchor-name": `--anchor-${anchorName}`,
      });
    }
  })

  $effect(() => {
    if (anchorElement && tooltip) anchorElement.setAttribute("aria-describedby", tooltipId);
  })

  let openTimeout = $state<ReturnType<typeof setTimeout>>();
  showPopup = async () => {
    clearTimeout(openTimeout);
    openTimeout = setTimeout(() => {
      if (!popover || popover.matches(":popover-open")) return;
      visible = true;
      popover.showPopover();
    }, delayed ? 1000 : 0);

    return new Promise<void>((resolve) => {
      promiseResolve = (() => {
        resolve();
      });
    })
  }

  hidePopup = () => {
    clearTimeout(openTimeout);
    if (!popover || !visible) return;
    visible = false;
    promiseResolve();
  }

  function transitionEnd() {
    if (!popover || visible || !popover.matches(":popover-open")) return;
    popover.hidePopover();
  }

  function popoverToggled(event: ToggleEvent) {
    if (event.newState != "closed") return;
    if (!visible) return;
    visible = false;
    promiseResolve();
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

    z-index: 1;

    border: 0;
    padding: var(--padding, #{dimensions.$gapSmall});
    max-width: 30vw;
    max-height: 50vh;
    box-shadow: decorations.$boxShadow;
    font-size: text.$fontSize;
    background-color: colors.$backgroundSecondary;
    color: colors.$foregroundSecondary;

    position: fixed;
    position-area: top;
    position-try-fallbacks: bottom, right, left;
    position-try-order: most-width;

    anchor-name: --popup;
    anchor-scope: --popup;

    --distance: #{dimensions.$gapSmall};
    margin: var(--distance);

    --currentBorderRadius: #{dimensions.$borderRadius};
    border-radius: var(--currentBorderRadius);

    opacity: 0;
    transition: opacity animations.$animationSpeed;
    //transition: opacity animations.$animationSpeed, display animations.$animationSpeed allow-discrete; // blocked by https://bugzilla.mozilla.org/show_bug.cgi?id=1882408

    &.visible:popover-open {
      opacity: 1;

      @starting-style {
        & {
          opacity: 0;
        }
      }
    }
  }

  .popup:not(.triangle) {
    --distance: 0;
  }

  .popup.matchWidth {
    width: anchor-size(width);
    max-width: none;
    overflow-x: hidden;
  }

  .popup.dialog {
    max-width: none;
    max-height: none;
    border: 0;
    border-radius: dimensions.$borderRadius;
    padding: 0;
    background-color: colors.$backgroundPrimary;
    color: colors.$foregroundPrimary;

    position-area: right;
    position-try-fallbacks: left, bottom, top;
    position-try-order: most-height;
  }

  .popup:popover-open {
    display: flex;
    flex-direction: column;
    gap: var(--padding, #{dimensions.$gapSmall});
  }

  :global(html[data-frost="true"]) .popup {
    background-color: color-mix(in srgb, colors.$backgroundSecondary 50%, transparent) !important;
    backdrop-filter: blur(dimensions.$blurLarge);
  }

  .popup:focus {
    outline: 0;
  }

  .tooltip {
    pointer-events: none;
  }

  .popup.triangle::before {
    content: "";
  
    z-index: -1;

    background-color: inherit;

    --size: calc(var(--distance) / 1.41421356 * 2);
    width: var(--size);
    height: var(--size);

    position: fixed;
    left: clamp(
      anchor(--popup left),
      anchor(var(--anchor) center),
      anchor(--popup right),
    );
    top: clamp(
      anchor(--popup top),
      anchor(var(--anchor) center),
      anchor(--popup bottom),
    );

    transform: translate(-50%, -50%) rotate(45deg);
  }

  .popup.triangle::after {
    content: "";
  
    z-index: -1;

    background-color: inherit;

    width: calc(2 * var(--distance));
    height: calc(2 * var(--distance));

    position: fixed;
    left: clamp(
      calc(anchor(--popup left) + var(--distance)),
      anchor(var(--anchor) center),
      calc(anchor(--popup right) - var(--distance)),
    );
    top: clamp(
      calc(anchor(--popup top) + var(--distance)),
      anchor(var(--anchor) center),
      calc(anchor(--popup bottom) - var(--distance)),
    );

    transform: translate(-50%, -50%);
  }
</style>

<!-- The typecast to "auto" is because the linter does not yet know about "hint" -->
<div
  bind:this={popover}
  class="popup"
  popover={tooltip ? "hint" : "auto"}
  style={`--anchor: --anchor-${anchorName}; position-anchor: var(--anchor);`}
  id={tooltipId}
  class:visible
  class:tooltip
  class:matchWidth
  class:triangle
  class:dialog
  tabindex="-1"
  ontransitionend={transitionEnd}
  role={tooltip ? "tooltip" : "dialog"}
  ontoggle={popoverToggled}
  aria-labelledby={labelled}
>
  {@render children?.()}
</div>
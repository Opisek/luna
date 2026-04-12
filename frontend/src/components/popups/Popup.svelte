<script lang="ts">
  import type { Snippet } from "svelte";
  import { NoOp } from "../../lib/client/placeholders";

  interface Props {
    tooltip?: boolean;
    delayed?: boolean;
    anchor?: HTMLElement | undefined;
    children?: Snippet;
    showPopup?: () => Promise<void>;
    hidePopup?: () => void;
  }

  let {
    tooltip = true,
    delayed = false,
    anchor = undefined,
    children,
    showPopup = $bindable(),
    hidePopup = $bindable(NoOp),
  }: Props = $props();

  let visible = $state(false);
  let popover: (HTMLElement | undefined) = $state();
  let anchorElement = $derived(anchor || (!popover ? undefined : popover.parentElement))
  let anchorName = $derived(`${Math.floor(Math.random() * 100000000)}-${anchorElement?.classList.values().toArray().join("-")}`);

  let promiseResolve: () => void = $state(NoOp);
  let promiseReject: (reason?: any) => void = $state(NoOp);

  $effect(() => {
    if (!anchorElement) return;
    // @ts-ignore
    const currentAnchor = anchorElement.style["anchor-name"] as string;
    if (currentAnchor.startsWith("--anchor-") && !currentAnchor.includes("undefined")) {
      anchorName = currentAnchor.substring(9);
    } else {
      Object.assign(anchorElement.style, {
        "anchor-name": `--anchor-${anchorName}`,
      });
      if (tooltip) anchorElement.setAttribute("aria-describedby", `tooltip-${anchorName}`);
    }
  })

  let openTimeout = $state<ReturnType<typeof setTimeout>>();
  showPopup = async () => {
    clearTimeout(openTimeout);
    openTimeout = setTimeout(() => {
      if (!popover || popover.matches(":popover-open")) return;
      visible = true;
      popover.showPopover();
    }, delayed ? 1000 : 0);

    if (!delayed && popover && !visible) {
      return new Promise<void>((resolve, reject) => {
        promiseResolve = (() => {
          resolve();
        });
        promiseReject = ((err) => {
          reject(err);
        });
      })
    }
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
    promiseReject();
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

    margin: dimensions.$gapSmall;

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

  .popup:popover-open {
    display: flex;
    flex-direction: column;
    gap: dimensions.$gapSmall;
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
</style>

<!-- The typecast to "auto" is because the linter does not yet know about "hint" -->
<div
  bind:this={popover}
  class="popup"
  popover={(tooltip ? "hint" : "auto") as "auto"}
  style={`position-anchor: --anchor-${anchorName};`}
  id={`${tooltip ? "tooltip" : "popup"}-${anchorName}`}
  class:visible={visible}
  class:tooltip={tooltip}
  tabindex="-1"
  ontransitionend={transitionEnd}
  role={tooltip ? "tooltip" : "dialog"}
  ontoggle={popoverToggled}
>
  {@render children?.()}
</div>
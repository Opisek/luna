<script lang="ts">
  import { browser } from "$app/environment";
  import { calculateOptimalPopupPosition } from "../../lib/common/calculations";

  let popup: HTMLElement;

  let bottom: boolean = false;
  let center: boolean = false;
  let right: boolean = false;

  $: setupListener(popup);

  function setupListener(popup: HTMLElement) {
    if (!popup || !popup.parentElement) return;

    popup.parentElement.addEventListener("mouseenter", checkPosition);
  }

  function checkPosition() {
    if (!popup || !popup.parentElement || !browser) return;

    const res = calculateOptimalPopupPosition(popup, 3);

    bottom = res.bottom;
    right = res.right;
    center = res.center;
  }
</script>

<style lang="scss">
  @import "../../styles/animations.scss";
  @import "../../styles/colors.scss";
  @import "../../styles/decoration.scss";
  @import "../../styles/dimensions.scss";

  span.popup {
    position: absolute;
    z-index: 10;
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
    border-radius: $borderRadius;
    padding: $paddingSmaller;
    background-color: $backgroundPrimary;
    color: $foregroundPrimary;
    transition: opacity $animationSpeed $cubic;
    z-index: 10;
    pointer-events: none;
    min-width: none;
    width: max-content;
    max-width: 80vw;
    box-shadow: $boxShadow;
    white-space: nowrap;
  }
  :global(*:hover) > span.popup > span.contents {
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
    <slot/>
  </span>
</span>
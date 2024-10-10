<script lang="ts">
  import { browser } from "$app/environment";
  import { calculateOptimalPopupPosition } from "../../lib/common/calculations";

  export let visible: boolean = false;

  let dialog: HTMLDialogElement;

  function clickOutside(event: MouseEvent) {
    if (!dialog) return;

    const clickX = event.clientX;
    const clickY = event.clientY;

    const rect = dialog.getBoundingClientRect();

    const minX = rect.left;
    const maxX = rect.right;
    const minY = rect.top;
    const maxY = rect.bottom;

    if (clickX < minX || clickX > maxX || clickY < minY || clickY > maxY) {
      close();
      event.stopPropagation();
    }
  }

  export const show = () => {
    visible = true;
    checkPosition();
    if (browser) {
      window.addEventListener("click", clickOutside);
    }
    setTimeout(() => dialog.show(), 0);
  }

  export const close = () => {
    visible = false;
    dialog.close();
    if (browser) {
      window.removeEventListener("click", clickOutside);
    }
  }

  let bottom: boolean = false;
  let center: boolean = false;
  let right: boolean = false;

  function checkPosition() {
    if (!dialog || !dialog.parentElement || !browser) return;

    const res = calculateOptimalPopupPosition(dialog, 5);

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

  dialog {
    border: 0;
    padding: $gapSmall $gap $gap $gap;
    border-radius: $borderRadius;
    max-width: 50vw;
    min-width: fit-content;
    box-shadow: $boxShadow;
    position: absolute !important;
    z-index: 10;
  }

  dialog[open] {
		animation: zoom $animationSpeed $cubic forwards;
	}

  div.contents {
    display: flex;
    flex-direction: column;
    gap: $gapSmall;
  }

  dialog.center {
    left: 0;
  }
  dialog.left {
    left: 100%;
  }
  dialog.right {
    left: -100%;
  }
  dialog.below {
    top: -100%;
  }
  dialog.above {
    top: 100%;
  }
</style>

<dialog
  bind:this={dialog}
  class:center={center}
  class:left={!right && !center}
  class:right={right && !center}
  class:below={bottom}
  class:above={!bottom}
>
  <div class="contents">
    <slot/>
  </div>
</dialog>
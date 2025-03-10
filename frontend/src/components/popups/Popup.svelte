<script lang="ts">
  import type { Snippet } from "svelte";

  import { browser } from "$app/environment";

  import { calculateOptimalPopupPosition } from "$lib/common/calculations";
  import { NoOp } from "../../lib/client/placeholders";

  interface Props {
    visible?: boolean;
    children?: Snippet;
    showPopup?: () => void;
    hidePopup?: () => void;
  }

  let {
    visible = $bindable(false),
    children,
    showPopup = $bindable(),
    hidePopup = $bindable(NoOp)
  }: Props = $props();

  let dialog: HTMLDialogElement;

  function clickOutside(event: MouseEvent) {
    if (!dialog || event.detail === 0) return;

    const clickX = event.clientX;
    const clickY = event.clientY;

    const rect = dialog.getBoundingClientRect();

    const minX = rect.left;
    const maxX = rect.right;
    const minY = rect.top;
    const maxY = rect.bottom;

    if (clickX < minX || clickX > maxX || clickY < minY || clickY > maxY) {
      hidePopup();
      event.stopPropagation();
    }
  }

  showPopup = () => {
    visible = true;
    checkPosition();
    dialog.show();

    if (browser) window.addEventListener("click", clickOutside);

    setTimeout(() => {
      dialog.focus();
    }, 0);
  }

  hidePopup = () => {
    visible = false;
    dialog.close();

    if (browser) window.removeEventListener("click", clickOutside);
  }

  let bottom: boolean = $state(false);
  let center: boolean = $state(false);
  let right: boolean = $state(false);

  function checkPosition() {
    if (!dialog || !dialog.parentElement || !browser) return;

    const res = calculateOptimalPopupPosition(dialog.parentElement, 5);

    bottom = res.bottom;
    right = res.right;
    center = res.center;
  }
</script>

<style lang="scss">
  @use "../../styles/animations.scss";
  @use "../../styles/colors.scss";
  @use "../../styles/decorations.scss";
  @use "../../styles/dimensions.scss";

  dialog {
    border: 0;
    padding: dimensions.$gapSmall dimensions.$gap dimensions.$gap dimensions.$gap;
    border-radius: dimensions.$borderRadius;
    max-width: 50vw;
    min-width: fit-content;
    box-shadow: decorations.$boxShadow;
    position: absolute !important;
    z-index: 10;
  }

  dialog[open] {
		animation: zoom animations.$animationSpeed animations.$cubic forwards;
	}

  dialog:focus {
    outline: none;
  }

  div.contents {
    display: flex;
    flex-direction: column;
    gap: dimensions.$gapSmall;
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
    {@render children?.()}
  </div>
</dialog>
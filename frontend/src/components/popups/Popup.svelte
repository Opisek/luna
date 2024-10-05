<script lang="ts">
  import { browser } from "$app/environment";

  let dialog: HTMLDialogElement;

  export const show = () => {
    checkPosition();
    setTimeout(() => dialog.show(), 0);
  }

  let bottom: boolean = false;
  let center: boolean = false;
  let right: boolean = false;

  function checkPosition() {
    if (!dialog || !dialog.parentElement || !browser) return;

    const rect = dialog.parentElement.getBoundingClientRect();

    const x = rect.left + (rect.right - rect.left) / 2;
    const y = rect.top + (rect.bottom - rect.top) / 2;

    bottom = y > window.innerHeight / 2;
    right = x > window.innerWidth / 2;
    center = x < window.innerWidth / 3 * 2 && x > window.innerWidth / 3;
  }
</script>

<style lang="scss">
  @import "../../styles/animations.scss";
  @import "../../styles/colors.scss";
  @import "../../styles/decoration.scss";
  @import "../../styles/dimensions.scss";

  dialog {
    border: 0;
    padding: $gap;
    border-radius: $borderRadius;
    max-width: 50vw;
    min-width: fit-content;
    box-shadow: $boxShadow;
    position: absolute !important;
    z-index: 10;
    left: 0;
    top: 0;
  }

  dialog[open] {
		animation: zoom $animationSpeed $cubic forwards;
	}
</style>

<dialog bind:this={dialog}>
  <div class="contents">
    <slot/>
  </div>
</dialog>
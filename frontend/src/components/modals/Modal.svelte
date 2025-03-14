<script lang="ts">
  import type { Snippet } from "svelte";

  import Button from "../interactive/Button.svelte";
  import CloseButton from "../interactive/CloseButton.svelte";
  import Horizontal from "../layout/Horizontal.svelte";
  import Title from "../layout/Title.svelte";

  import { NoOp } from "$lib/client/placeholders";
  import { redrawNotifications } from "$lib/client/notifications";

  interface Props {
    title: string;
    focusElement?: HTMLElement | null;
    onModalHide?: any;
    onModalSubmit?: any;
    showModal?: () => any;
    hideModal?: () => any;
    resetFocus?: () => any;
    children?: Snippet;
    buttons?: Snippet;
  }

  let {
    title,
    focusElement = null,
    onModalHide = () => {},
    showModal = $bindable(),
    hideModal = $bindable(NoOp),
    resetFocus = $bindable(),
    onModalSubmit = hideModal,
    children,
    buttons
  }: Props = $props();

  let dialog: HTMLDialogElement;

  let visible = $state(false);

  function clickOutside(event: MouseEvent) {
    if (!dialog) return;
    if (event.target === dialog) {
      hideModal();
      event.stopPropagation();
    }
  }

  $effect(() => {
    if (visible) dialog.showModal();
  });

  resetFocus = () => {
    if (focusElement) focusElement.focus();
    else dialog.focus();
  }
  showModal = () => {
    window.addEventListener("click", clickOutside);
    visible = true
    setTimeout(resetFocus, 0);
    setTimeout(redrawNotifications, 0); // hacky way to make sure that notifications are always on the very top. sometimes has a visible blink. should revisit one day.
  }
  hideModal = () => {
    window.removeEventListener("click", clickOutside);
    dialog.close();
    onModalHide();
  }

  function submitInternal(event: Event) {
    event.preventDefault();
    onModalSubmit();
    return false;
  }
</script>

<style lang="scss">
  @use "../../styles/animations.scss";
  @use "../../styles/colors.scss";
  @use "../../styles/dimensions.scss";

  dialog {
    border: 0;
    max-width: 50vw;
    min-width: 30em;
    border-radius: dimensions.$borderRadius;
    padding: 0;
  }
  dialog::backdrop {
    backdrop-filter: blur(dimensions.$blur);
  }

  dialog[open] {
		animation: zoom animations.$animationSpeed animations.$cubic forwards;
	}

  dialog:focus {
    outline: none;
  }
  
  form {
    padding: dimensions.$gap dimensions.$gapLarge dimensions.$gapLarge dimensions.$gapLarge;
    border-radius: dimensions.$borderRadius;
    display: flex;
    width: 100%;
    flex-direction: column;
    flex-wrap: nowrap;
    gap: dimensions.$gap;
  }
</style>

<dialog
  bind:this={dialog}
  onclose={() => (visible = false)}
  class:closed={visible}
>
  {#if visible}
    <form onsubmit={submitInternal}>
      <Horizontal>
        <Title>
          {title}
        </Title>
        <CloseButton onClick={hideModal} />
      </Horizontal>
      {@render children?.()}
      <Horizontal position="right">
        {#if buttons}{@render buttons()}{:else}
          <Button onClick={hideModal}>Close</Button>
        {/if}
      </Horizontal>
    </form>
  {/if}
</dialog>
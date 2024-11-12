<script lang="ts">
  import type { Snippet } from "svelte";

  import Button from "../interactive/Button.svelte";
  import CloseButton from "../interactive/CloseButton.svelte";
  import Horizontal from "../layout/Horizontal.svelte";
  import Title from "../layout/Title.svelte";

  import { NoOp } from "$lib/client/placeholders";

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
  @import "../../styles/animations.scss";
  @import "../../styles/colors.scss";
  @import "../../styles/dimensions.scss";

  dialog {
    border: 0;
    max-width: 50vw;
    min-width: 30em;
    border-radius: $borderRadius;
    padding: 0;
  }
  dialog::backdrop {
    backdrop-filter: blur($blur);
  }

  dialog[open] {
		animation: zoom $animationSpeed $cubic forwards;
	}

  dialog:focus {
    outline: none;
  }
  
  form {
    padding: $gap $gapLarge $gapLarge $gapLarge;
    border-radius: $borderRadius;
    display: flex;
    width: 100%;
    flex-direction: column;
    flex-wrap: nowrap;
    gap: $gap;
  }
</style>

<dialog
  bind:this={dialog}
  onclose={() => (visible = false)}
  class:closed={visible}
>
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
        <Button onClick={hideModal} color="accent">Close</Button>
      {/if}
    </Horizontal>
	</form>
</dialog>
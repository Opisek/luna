<script lang="ts">
  import Button from "../interactive/Button.svelte";
  import CloseButton from "../interactive/CloseButton.svelte";
  import Horizontal from "../layout/Horizontal.svelte";
  import Title from "../layout/Title.svelte";

  let dialog: HTMLDialogElement;

  export let title: string;
  export let focusElement: HTMLElement | null = null;

  let visible = false;

  function clickOutside(event: MouseEvent) {
    if (!dialog) return;
    if (event.target === dialog) {
      hideModal();
      event.stopPropagation();
    }
  }

  $: if (dialog && visible) dialog.showModal();

  export const showModal = () => {
    window.addEventListener("click", clickOutside);
    visible = true
    setTimeout(resetFocus, 0);
  }
  export const hideModal = () => {
    window.removeEventListener("click", clickOutside);
    dialog.close();
    onModalHide();
  }
  export let onModalHide = () => {};

  export let onModalSubmit = hideModal;

  function submitInternal(event: Event) {
    event.preventDefault();
    onModalSubmit();
    return false;
  }

  export let resetFocus = () => {
    if (focusElement) focusElement.focus();
    else dialog.focus();
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
  on:close={() => (visible = false)}
  class:closed={visible}
>
	<form on:submit={submitInternal}>
    <Horizontal>
      <Title>
        {title}
      </Title>
      <CloseButton onClick={hideModal} />
    </Horizontal>
		<slot />
    <Horizontal position="right">
      <slot name="buttons">
        <Button onClick={hideModal} color="accent">Close</Button>
      </slot>
    </Horizontal>
	</form>
</dialog>
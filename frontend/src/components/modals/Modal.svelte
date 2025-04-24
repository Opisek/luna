<script lang="ts">
  import type { Snippet } from "svelte";

  import Button from "../interactive/Button.svelte";
  import CloseButton from "../interactive/CloseButton.svelte";
  import Horizontal from "../layout/Horizontal.svelte";
  import Title from "../layout/Title.svelte";

  import { NoOp } from "$lib/client/placeholders";
  import { redrawNotifications } from "$lib/client/notifications";
  import { isChildOfModal } from "../../lib/common/misc";

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
    topButtons?: Snippet;
  }

  let {
    title,
    focusElement = null,
    onModalHide = NoOp,
    showModal = $bindable(),
    hideModal = $bindable(NoOp),
    resetFocus = $bindable(),
    onModalSubmit = hideModal,
    children,
    buttons,
    topButtons,
  }: Props = $props();

  let dialog: HTMLDialogElement;

  let visible = $state(false);

  let ignoreClickOutside = $state(false);
  function mouseDown(event: MouseEvent) {
    ignoreClickOutside = isChildOfModal(event.target as HTMLElement) && event.target !== dialog;
  }

  function clickOutside(event: MouseEvent) {
    if (!dialog || ignoreClickOutside) return;
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
    window.addEventListener("mousedown", mouseDown);
    window.addEventListener("click", clickOutside);
    visible = true
    setTimeout(resetFocus, 0);
    setTimeout(redrawNotifications, 0); // hacky way to make sure that notifications are always on the very top. sometimes has a visible blink. should revisit one day.
  }
  hideModal = () => {
    window.removeEventListener("mousedown", mouseDown);
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
    border-radius: dimensions.$borderRadius;
    padding: 0;
    background-color: colors.$backgroundPrimary;
    color: colors.$foregroundPrimary;
  }
  :global(html[data-frost="true"]) dialog {
    background-color: color-mix(in srgb, colors.$backgroundPrimary 75%, transparent);
    backdrop-filter: blur(dimensions.$blurLarge);
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
    padding: dimensions.$gapLarge dimensions.$gapLarger dimensions.$gapLarger dimensions.$gapLarger;
    border-radius: dimensions.$borderRadius;
    display: flex;
    flex-direction: column;
    flex-wrap: nowrap;
    gap: dimensions.$gapMiddle;
    box-sizing: content-box;
    min-width: 30em;
    width: fit-content;
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
        {#if topButtons}
          <Horizontal position="right">
            {@render topButtons()}
            <CloseButton onClick={hideModal} />
          </Horizontal>
        {:else}
          <CloseButton onClick={hideModal} />
        {/if}
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
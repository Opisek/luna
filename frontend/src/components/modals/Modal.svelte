<script lang="ts" generics="T">
  import type { Snippet } from "svelte";

  import CloseButton from "../interactive/CloseButton.svelte";
  import Horizontal from "../layout/Horizontal.svelte";
  import Title from "../layout/Title.svelte";

  import { AsyncNoOp, NoOp } from "$lib/client/placeholders";
  import { redrawNotifications } from "$lib/client/notifications";
  import Popup from "../popups/Popup.svelte";

  interface Props {
    title: string;
    onModalHide?: any;
    onModalSubmit?: any;
    showModal: (anchor?: HTMLElement) => Promise<T>;
    success?: (result: T) => void;
    failure?: (reason?: string | Error) => void;
    deanchor?: () => void;
    children?: Snippet;
    buttons?: Snippet;
    topButtons?: Snippet;
  }

  let {
    title,
    onModalHide = NoOp,
    showModal = $bindable(),
    success = $bindable(),
    failure = $bindable(),
    deanchor = $bindable(),
    onModalSubmit = NoOp,
    children,
    buttons,
    topButtons,
  }: Props = $props();

  let dialog: HTMLDialogElement | undefined = $state();
  let showPopup = $state(AsyncNoOp);
  let hidePopup = $state(NoOp);

  const dialogId = $state(`dialog-${Math.floor(Math.random() * 100000000)}`);

  let visible = $state(false);
  let anchor: HTMLElement | undefined = $state()
  let popup = $derived(anchor != undefined);

  let promiseResolve: (result: T) => void = $state(NoOp);
  let promiseReject: (reason?: string | Error) => void = $state(NoOp);

  showModal = (anchorElement?: HTMLElement) => {
    anchor = anchorElement;
    visible = true
    setTimeout(() => {
      if (dialog) dialog.showModal();
      else showPopup().finally(failure);
    }, 10);
    setTimeout(redrawNotifications, 0); // hacky way to make sure that notifications are always on the very top. sometimes has a visible blink. should revisit one day.
    return new Promise<T>((resolve, reject) => {
      promiseResolve = ((result) => {
        promiseResolve = NoOp;
        promiseReject = NoOp;
        resolve(result);
      });
      promiseReject = ((err) => {
        promiseResolve = NoOp;
        promiseReject = NoOp;
        reject(err);
      });
    })
  }
  success = (result) => {
    promiseResolve(result);
    hideModal();
  }
  failure = (error) => {
    promiseReject(error);
    hideModal();
  }
  deanchor = () => {
    if (!visible || anchor == undefined) return;
    anchor = undefined;
    setTimeout(() => dialog?.showModal(), 10);
  }

  function hideModal() {
    if (dialog) dialog.close();
    else hidePopup();
  }

  function modalHideInternal() {
    visible = false;
    promiseReject();
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
    min-width: 30rem;
    width: fit-content;
  }

  form.popup {
    padding: dimensions.$gapLarge;
    min-width: 0;
  }
</style>

{#if popup}
  <Popup
    tooltip={false}
    anchor={anchor}
    dialog={true}
    bind:visible
    bind:showPopup
    bind:hidePopup
  >
    {@render content()}
  </Popup>
{:else}
  <dialog
    bind:this={dialog}
    closedby="any"
    onclose={modalHideInternal}
    tabindex="-1"
    id={dialogId}
    aria-describedby={`title-${dialogId}`}
    aria-modal="true"
  >
    {@render content()}
  </dialog>
{/if}

{#snippet content()}
  {#if visible}
    <form onsubmit={submitInternal} class:popup>
      {#if !popup}
        <Horizontal>
          <Title id={`title-${dialogId}`}>
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
      {/if}
      {@render children?.()}
      {#if buttons}
        <Horizontal position="right">
          {@render buttons()}
        </Horizontal>
      {/if}
    </form>
  {/if}
{/snippet}
<script lang="ts" generics="T">
  import type { Snippet } from "svelte";

  import CloseButton from "../interactive/CloseButton.svelte";
  import Horizontal from "../layout/Horizontal.svelte";
  import Title from "../layout/Title.svelte";

  import { NoOp } from "$lib/client/placeholders";
  import { redrawNotifications } from "$lib/client/notifications";

  interface Props {
    title: string;
    onModalHide?: any;
    onModalSubmit?: any;
    showModal: () => Promise<T>;
    success?: (result: T) => void;
    failure?: (reason?: string | Error) => void;
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
    onModalSubmit = NoOp,
    children,
    buttons,
    topButtons,
  }: Props = $props();

  let dialog: HTMLDialogElement;

  let visible = $state(false);

  let promiseResolve: (result: T) => void = $state(NoOp);
  let promiseReject: (reason?: string | Error) => void = $state(NoOp);

  showModal = () => {
    visible = true
    dialog.showModal();
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

  function hideModal() {
    dialog.close();
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
    min-width: 30em;
    width: fit-content;
  }
</style>

<dialog
  bind:this={dialog}
  closedby="any"
  onclose={modalHideInternal}
  class:closed={visible}
  tabindex="-1"
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
          <!--<Button onClick={hideModal}>Close</Button>-->
        {/if}
      </Horizontal>
    </form>
  {/if}
</dialog>
<script lang="ts">
  import Loader from "../decoration/Loader.svelte";
  import Button from "../interactive/Button.svelte";
  import Modal from "./Modal.svelte";

  import { queueNotification } from "$lib/client/notifications";
  import { NoOp } from "../../lib/client/placeholders";
  import { ColorKeys } from "../../types/colors";

  interface Props {
    confirmCallback: () => Promise<void>;
    cancelCallback?: () => void;
    showModal: () => any;
    hideModal?: () => any;
    children?: import('svelte').Snippet;
  }

  let {
    showModal = $bindable(),
    hideModal = $bindable(NoOp),
    confirmCallback,
    cancelCallback = () => {},
    children
  }: Props = $props();

  let awaitingConfirm = $state(false);
  function confirm() {
    awaitingConfirm = true;
    confirmCallback().catch(err => {
      queueNotification(ColorKeys.Danger, err)
    }).finally(() => {
      awaitingConfirm = false;
      hideModal()
    });
  }

  function cancel() {
    hideModal()
    cancelCallback()
  }
</script>

<Modal title="Confirmation" bind:showModal={showModal} bind:hideModal={hideModal}>
  {@render children?.()}
  {#snippet buttons()}
      <Button onClick={confirm} color={ColorKeys.Success} type="submit">
        {#if awaitingConfirm}
          <Loader/>
        {:else}
          Confirm
        {/if}
      </Button>
      <Button onClick={cancel} color={ColorKeys.Danger}>Cancel</Button>
  {/snippet}
</Modal>
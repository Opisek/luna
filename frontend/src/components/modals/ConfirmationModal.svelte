<script lang="ts">
  import Loader from "../decoration/Loader.svelte";
  import Button from "../interactive/Button.svelte";
  import Modal from "./Modal.svelte";

  import { queueNotification } from "$lib/client/notifications";
  import { NoOp } from "../../lib/client/placeholders";

  interface Props {
    confirmCallback: () => Promise<string>;
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
  async function confirm() {
    // TOOD: error message if returned value is not empty string
    awaitingConfirm = true;
    const res = await confirmCallback()
    awaitingConfirm = false;
    hideModal()

    if (res !== "") {
      queueNotification("failure", res)
    }
  }

  function cancel() {
    hideModal()
    cancelCallback()
  }
</script>

<Modal title="Confirmation" bind:showModal={showModal} bind:hideModal={hideModal}>
  {@render children?.()}
  {#snippet buttons()}
      <Button onClick={confirm} color="success">
        {#if awaitingConfirm}
          <Loader/>
        {:else}
          Confirm
        {/if}
      </Button>
      <Button onClick={cancel} color="failure">Cancel</Button>
  {/snippet}
</Modal>
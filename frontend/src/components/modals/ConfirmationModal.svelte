<script lang="ts">
  import Loader from "../decoration/Loader.svelte";
import Button from "../interactive/Button.svelte";
  import Modal from "./Modal.svelte";

  export let showModal: () => boolean;
  export let confirmCallback: () => Promise<string>;
  export let cancelCallback: () => void = () => {};

  let hideModal: () => void;

  let awaitingConfirm = false;
  async function confirm() {
    // TOOD: error message if returned value is not empty string
    awaitingConfirm = true;
    await confirmCallback()
    awaitingConfirm = false;
    hideModal()
  }

  function cancel() {
    hideModal()
    cancelCallback()
  }
</script>

<Modal title="Confirmation" bind:showModal={showModal} bind:hideModal={hideModal}>
  <slot/>
  <svelte:fragment slot="buttons">
    <Button onClick={confirm} color="success">
      {#if awaitingConfirm}
        <Loader/>
      {:else}
        Confirm
      {/if}
    </Button>
    <Button onClick={cancel} color="failure">Cancel</Button>
  </svelte:fragment>
</Modal>
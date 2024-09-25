<script lang="ts">
  import Button from "../interactive/Button.svelte";
  import Modal from "./Modal.svelte";

  export let showModal: () => boolean;
  export let confirmCallback: () => void;
  export let cancelCallback: () => void = () => {};

  let hideModal: () => void;

  function confirm() {
    hideModal()
    confirmCallback()
  }

  function cancel() {
    hideModal()
    cancelCallback()
  }
</script>

<Modal title="Confirmation" bind:showModal={showModal} bind:hideModal={hideModal}>
  <slot/>
  <svelte:fragment slot="buttons">
    <Button onClick={confirm} color="success">Confirm</Button>
    <Button onClick={cancel} color="failure">Cancel</Button>
  </svelte:fragment>
</Modal>
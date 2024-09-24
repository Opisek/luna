<script lang="ts">
  import Button from "../interactive/Button.svelte";
  import Horizontal from "../layout/Horizontal.svelte";
  import Vertical from "../layout/Vertical.svelte";
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
  <Vertical>
    <span>
      <slot/>
    </span>
    <Horizontal position="right">
      <Button onClick={confirm} color="success">Confirm</Button>
      <Button onClick={cancel} color="failure">Cancel</Button>
    </Horizontal>
  </Vertical>
</Modal>
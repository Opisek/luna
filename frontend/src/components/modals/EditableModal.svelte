<script lang="ts">
  import Loader from "../decoration/Loader.svelte";
  import Button from "../interactive/Button.svelte";
  import ConfirmationModal from "./ConfirmationModal.svelte";
  import Modal from "./Modal.svelte";

  export const showModal = () => {
    editMode = isNew;
    showModalInternal();
  };
  let showModalInternal: () => boolean;
  let hideModal: () => boolean;

  export let title: string;
  export let deleteConfirmation: string;
  export let isNew: boolean;

  export let editMode: boolean = false;

  function startEditMode() {
    editMode = true;
  }

  function cancelEdit() {
    editMode = false;
  }

  let awaitingEdit = false;
  async function saveEdit() {
    // TOOD: error message if returned value is not empty string
    awaitingEdit = true;
    await onEdit();
    awaitingEdit = false;

    editMode = false;
    hideModal();
  }

  let showDeleteModal: () => boolean;
  const confirmDelete = async () => {
    const returnValue = await onDelete();
    hideModal();
    return returnValue;
  }

  export let onEdit: () => Promise<string>;
  export let onDelete: () => Promise<string>;
</script>

<Modal title={title} bind:showModal={showModalInternal} bind:hideModal={hideModal} onModalHide={cancelEdit}>
  <slot/>
  <svelte:fragment slot="buttons">
    {#if editMode}
      <Button onClick={saveEdit} color="success">
        {#if awaitingEdit}
          <Loader/>
        {:else}
          Save
        {/if}
      </Button>
      <Button onClick={cancelEdit} color="failure">Cancel</Button>
    {:else}
      <Button onClick={startEditMode} color="accent">Edit</Button>
      <Button onClick={showDeleteModal} color="failure">Delete</Button>
    {/if}
  </svelte:fragment>
</Modal>

<ConfirmationModal
  bind:showModal={showDeleteModal}
  confirmCallback={confirmDelete}
>
  {deleteConfirmation}
  <br>
  This action is irreversible.
</ConfirmationModal>
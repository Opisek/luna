<script lang="ts">
  import { queueNotification } from "$lib/client/notifications";
  import Loader from "../decoration/Loader.svelte";
  import Button from "../interactive/Button.svelte";
  import ConfirmationModal from "./ConfirmationModal.svelte";
  import Modal from "./Modal.svelte";

  export let title: string;
  export let deleteConfirmation: string;

  export let editMode: boolean = false;
  let creating = false;

  export let submittable = true;

  export const showCreateModal = () => {
    creating = true;
    editMode = true;
    showModalInternal();
  };
  export const showModal = () => {
    creating = false;
    editMode = false;
    showModalInternal();
  };
  let showModalInternal: () => boolean;
  let hideModal: () => boolean;


  let resetFocus: () => any;

  function startEditMode() {
    resetFocus();
    editMode = true;
  }

  function cancelEdit() {
    resetFocus();
    editMode = false;
    if (creating) {
      hideModal();
    }
  }

  let awaitingEdit = false;
  async function saveEdit() {
    // TOOD: error message if returned value is not empty string
    awaitingEdit = true;
    const res = await onEdit();
    awaitingEdit = false;

    if (res === "") {
      editMode = false;
      hideModal();
    } else {
      queueNotification("failure", res)
    }
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

<Modal title={title} bind:showModal={showModalInternal} bind:hideModal={hideModal} onModalHide={() => {editMode = false}} bind:resetFocus>
  <slot/>
  <svelte:fragment slot="buttons">
    {#if editMode}
      <Button onClick={saveEdit} color="success" enabled={submittable}>
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
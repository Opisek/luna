<script lang="ts">
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

  function saveEdit() {
    editMode = false;
    onEdit();
  }

  let showDeleteModal: () => boolean;
  const confirmDelete = () => {
    hideModal();
    onDelete();
  }

  export let onEdit: () => void;
  export let onDelete: () => void;
</script>

<Modal title={title} bind:showModal={showModalInternal} bind:hideModal={hideModal} onModalHide={cancelEdit}>
  <slot/>
  <svelte:fragment slot="buttons">
    {#if editMode}
      <Button onClick={saveEdit} color="success">Save</Button>
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
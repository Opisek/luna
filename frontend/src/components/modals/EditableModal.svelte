<script lang="ts">
  import type { Snippet } from "svelte";

  import Loader from "../decoration/Loader.svelte";
  import Button from "../interactive/Button.svelte";
  import ConfirmationModal from "./ConfirmationModal.svelte";

  import { NoOp } from "../../lib/client/placeholders";
  import { queueNotification } from "$lib/client/notifications";
 
  interface Props {
    title: string;
    deleteConfirmation: string;
    editMode?: boolean;
    submittable?: boolean;
    onEdit: () => Promise<string>;
    onDelete: () => Promise<string>;
    showCreateModal?: () => any;
    showModal?: () => any;
    hideModal?: () => any;
    children?: Snippet;
  }

  let {
    title,
    deleteConfirmation,
    editMode = $bindable(false),
    submittable = true,
    onEdit,
    onDelete,
    showCreateModal = $bindable(),
    showModal = $bindable(),
    hideModal = $bindable(NoOp),
    children,
  }: Props = $props(); import Modal from "./Modal.svelte";

  let creating = false;

  let showModalInternal: () => any = $state(NoOp);
  let showDeleteModal: () => any = $state(NoOp);
  let resetFocus: () => any = $state(NoOp);

  showCreateModal = () => {
    creating = true;
    editMode = true;
    showModalInternal();
  };
  
  showModal = () => {
    creating = false;
    editMode = false;
    showModalInternal();
  };

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

  let awaitingEdit = $state(false);
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

  const confirmDelete = async () => {
    const returnValue = await onDelete();
    hideModal();
    return returnValue;
  }

</script>

<Modal title={title} bind:showModal={showModalInternal} bind:hideModal={hideModal} onModalHide={() => {editMode = false}} bind:resetFocus>
  {@render children?.()}
  {#snippet buttons()}
  
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
    
  {/snippet}
</Modal>

<ConfirmationModal
  bind:showModal={showDeleteModal}
  confirmCallback={confirmDelete}
>
  {deleteConfirmation}
  <br>
  This action is irreversible.
</ConfirmationModal>
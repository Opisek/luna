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
    editable?: boolean;
    deletable?: boolean;
    submittable?: boolean;
    onEdit: () => Promise<void>;
    onDelete: () => Promise<void>;
    showCreateModal?: () => any;
    showModal?: () => any;
    hideModal?: () => any;
    children?: Snippet;
    extraButtonsLeft?: Snippet;
    extraButtonsRight?: Snippet;
  }

  let {
    title,
    deleteConfirmation,
    editMode = $bindable(false),
    editable = true,
    deletable = true,
    submittable = true,
    onEdit,
    onDelete,
    showCreateModal = $bindable(),
    showModal = $bindable(),
    hideModal = $bindable(NoOp),
    children,
    extraButtonsLeft,
    extraButtonsRight,
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
  function saveEdit() {
    awaitingEdit = true;
    onEdit().then(() => {
      editMode = false;
      queueNotification("success", "Saved successfully")
      hideModal()
    }).catch((err) => {
      queueNotification("failure", err)
    }).finally(() => {
      awaitingEdit = false;
    });
  }

  const confirmDelete = async () => {
    await onDelete().then(() => {
      queueNotification("success", "Deleted successfully")
      hideModal();
    }).catch((err) => {
      queueNotification("failure", err)
    });
  }
</script>

<Modal title={title} bind:showModal={showModalInternal} bind:hideModal={hideModal} onModalHide={() => {editMode = false}} bind:resetFocus>
  {@render children?.()}
  {#snippet buttons()}
    {@render extraButtonsLeft?.()}
    {#if editMode}
      <Button onClick={saveEdit} color="success" enabled={submittable} type="submit">
        {#if awaitingEdit}
          <Loader/>
        {:else}
          Save
        {/if}
      </Button>
      <Button onClick={cancelEdit} color="failure">Cancel</Button>
    {:else}
      {#if editable}
        <Button onClick={startEditMode} color="accent">Edit</Button>
      {/if}
      {#if deletable}
        <Button onClick={showDeleteModal} color="failure">Delete</Button>
      {/if}
      {#if !editable && !deletable}
        <Button onClick={hideModal}>Close</Button>
      {/if}
    {/if}
    {@render extraButtonsRight?.()}
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
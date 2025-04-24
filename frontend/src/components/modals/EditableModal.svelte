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
    onCancel?: () => any;
    showCreateModal?: () => any;
    showModal?: () => any;
    hideModal?: () => any;
    children?: Snippet;
    extraButtonsTop?: Snippet;
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
    onCancel,
    showCreateModal = $bindable(),
    showModal = $bindable(),
    hideModal = $bindable(NoOp),
    children,
    extraButtonsTop,
    extraButtonsLeft,
    extraButtonsRight,
  }: Props = $props(); import Modal from "./Modal.svelte";
  import { ColorKeys } from "../../types/colors";
  import { D } from "svelte-simples";

  let creating = false;

  let showModalInternal: () => any = $state(NoOp);
  let hideModalInternal: () => any = $state(NoOp);
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
  hideModal = () => {
    onCancel?.();
    hideModalInternal();
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
      queueNotification(ColorKeys.Success, "Saved successfully")
      hideModal()
    }).catch((err) => {
      queueNotification(ColorKeys.Danger, err)
    }).finally(() => {
      awaitingEdit = false;
    });
  }

  const confirmDelete = async () => {
    await onDelete().then(() => {
      queueNotification(ColorKeys.Success, "Deleted successfully")
      hideModal();
    }).catch((err) => {
      queueNotification(ColorKeys.Danger, err)
    });
  }
</script>

<Modal
  title={title}
  bind:showModal={showModalInternal}
  bind:hideModal={hideModalInternal}
  onModalHide={() => {editMode = false}}
  bind:resetFocus
  topButtons={extraButtonsTop}
>
  {@render children?.()}
  {#snippet buttons()}
    {@render extraButtonsLeft?.()}
    {#if editMode}
      <Button onClick={saveEdit} color={ColorKeys.Success} enabled={submittable} type="submit">
        {#if awaitingEdit}
          <Loader/>
        {:else}
          Save
        {/if}
      </Button>
      <Button onClick={cancelEdit} color={ColorKeys.Danger}>Cancel</Button>
    {:else}
      {#if editable}
        <Button onClick={startEditMode} color={ColorKeys.Accent}>Edit</Button>
      {/if}
      {#if deletable}
        <Button onClick={showDeleteModal} color={ColorKeys.Danger}>Delete</Button>
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
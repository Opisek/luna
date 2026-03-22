<script lang="ts">
  import type { Snippet } from "svelte";

  import ConfirmationModal from "./ConfirmationModal.svelte";
  import Modal from "./Modal.svelte";

  import { NoOp } from "../../lib/client/placeholders";
  import { queueNotification } from "$lib/client/notifications";

  import { ColorKeys } from "../../types/colors";
  import IconButton from "../interactive/IconButton.svelte";
  import { Check, Pencil, Trash, Trash2, X } from "lucide-svelte";
 
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
    showEditModal?: () => any;
    showModal?: () => any;
    hideModal?: () => any;
    onModalHide?: () => any;
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
    showEditModal = $bindable(),
    showModal = $bindable(),
    hideModal = $bindable(NoOp),
    onModalHide = $bindable(NoOp),
    children,
    extraButtonsTop,
    extraButtonsLeft,
    extraButtonsRight,
  }: Props = $props();

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
  showEditModal = () => {
    creating = false;
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

  async function saveEdit() {
    return onEdit().then(() => {
      editMode = false;
      queueNotification(ColorKeys.Success, "Saved successfully")
      hideModal()
    }).catch((err) => {
      queueNotification(ColorKeys.Danger, err)
    });
  }

  const confirmDelete = async () => {
    return onDelete().then(() => {
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
  onModalHide={() => {editMode = false; onModalHide();}}
  bind:resetFocus
  topButtons={extraButtonsTop}
>
{@render children?.()}
{#snippet buttons()}
    {@render extraButtonsLeft?.()}
    {#if editMode}
      <IconButton onClick={saveEdit} color={ColorKeys.Success} enabled={submittable} type="submit" alt="Save" canRenderAsButton={true}>
        <Check/>
      </IconButton>
      <IconButton onClick={cancelEdit} color={ColorKeys.Danger} alt="Cancel" canRenderAsButton={true}><X/></IconButton>
    {:else}
      {#if editable}
        <IconButton onClick={startEditMode} alt="Edit" canRenderAsButton={true}>
          <Pencil/>
        </IconButton>
      {/if}
      {#if deletable}
        <IconButton onClick={showDeleteModal} color={ColorKeys.Danger} alt="Delete" canRenderAsButton={true}>
          <Trash2/>
        </IconButton>
      {/if}
      <!--{#if !editable && !deletable}
        <Button onClick={hideModal}>Close</Button>
      {/if}-->
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
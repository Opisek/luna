<script lang="ts" generics="T extends { id: string }">
  import type { Snippet } from "svelte";

  import ConfirmationModal from "./ConfirmationModal.svelte";
  import Modal from "./Modal.svelte";

  import { AsyncNoOp, NoOp } from "../../lib/client/placeholders";
  import { queueNotification } from "$lib/client/notifications";

  import { ColorKeys } from "../../types/colors";
  import IconButton from "../interactive/IconButton.svelte";
  import { Check, Pencil, Trash2, X } from "lucide-svelte";
 
  interface Props {
    title: string;
    deleteConfirmation: string;
    editMode?: boolean;
    editable?: boolean;
    deletable?: boolean;
    submittable?: boolean;

    onEdit: () => Promise<T>;
    onDelete: () => Promise<T>;

    showModal?: (initial?: T, edit?: boolean) => Promise<T>;
    success?: (result: T) => void;
    failure?: (reason?: string | Error) => void;

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
    showModal = $bindable(),
    success = $bindable(NoOp),
    failure = $bindable(NoOp),
    onModalHide = $bindable(NoOp),
    children,
    extraButtonsTop,
    extraButtonsLeft,
    extraButtonsRight,
  }: Props = $props();

  let creating = false;

  let showModalInternal: () => Promise<T> = $state(Promise.reject);

  let showDeleteModal: () => Promise<void> = $state(AsyncNoOp);

  showModal = (initial?: T, edit?: boolean) => {
    creating = !initial || initial.id === "";
    editMode = edit || creating;
    return showModalInternal();
  };

  function startEditMode() {
    editMode = true;
  }

  function cancelEdit() {
    editMode = false;
    if (creating) failure();
  }

  async function saveEdit() {
    return onEdit().then((result) => {
      editMode = false;
      queueNotification(ColorKeys.Success, "Saved successfully")
      success(result);
    }).catch((err) => {
      queueNotification(ColorKeys.Danger, err)
    });
  }

  const confirmDelete = async () => {
    return onDelete().then((result) => {
      queueNotification(ColorKeys.Success, "Deleted successfully")
      success(result);
    }).catch((err) => {
      queueNotification(ColorKeys.Danger, err)
    });
  }
</script>

<Modal
  title={title}
  bind:showModal={showModalInternal}
  bind:success
  bind:failure
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
        <IconButton onClick={() => showDeleteModal().then(confirmDelete).catch(NoOp)} color={ColorKeys.Danger} alt="Delete" canRenderAsButton={true}>
          <Trash2/>
        </IconButton>
      {/if}
    {/if}
    {@render extraButtonsRight?.()}
  {/snippet}
</Modal>

<ConfirmationModal bind:showModal={showDeleteModal}>
  {deleteConfirmation}
  <br>
  This action is irreversible.
</ConfirmationModal>
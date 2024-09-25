<script lang="ts">
  import Button from "../interactive/Button.svelte";
  import ConfirmationModal from "./ConfirmationModal.svelte";
  import Modal from "./Modal.svelte";
  import TextInput from "../forms/TextInput.svelte";

  export let source: SourceModel;
  export let showModal: () => boolean;
  let hideModal: () => boolean;

  let title: string;
  $: title = (source && source.id) ? (editMode ? "Edit source" : "Source") : "Create source";

  let editMode: boolean;
  $: editMode = !(source && source.id);

  let sourceCopy: SourceModel;
  function startEditMode() {
    editMode = true;
    sourceCopy = JSON.parse(JSON.stringify(event));
  }

  function cancelEdit() {
    editMode = false;
    source = sourceCopy;
  }

  function saveEdit() {
    editMode = false;
    console.log("Save");
  }

  let showDeleteModal: () => boolean;
  const deleteEvent = () => {
    hideModal();
    console.log("Delete");
  }
</script>

<Modal title={title} bind:showModal={showModal} bind:hideModal={hideModal}>
  {#if event}
    <TextInput bind:value={source.name} name="name" placeholder="Name" editable={editMode} />
  {/if}
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
  confirmCallback={deleteEvent}
>
{#if source}
  Do you really want to delete source "{source.name}"?
  <br>
  This action is irreversible.
{/if}
</ConfirmationModal>
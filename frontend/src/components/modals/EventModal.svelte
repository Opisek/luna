<script lang="ts">
  import DateTimeInput from "../forms/DateTimeInput.svelte";
  import TextInput from "../forms/TextInput.svelte";
  import Button from "../interactive/Button.svelte";
  import Horizontal from "../layout/Horizontal.svelte";
  import Vertical from "../layout/Vertical.svelte";
  import ConfirmationModal from "./ConfirmationModal.svelte";
  import Modal from "./Modal.svelte";

  export let event: EventModel;
  export let showModal: () => boolean;
  let hideModal: () => boolean;

  let title: string;
  $: title = (event && event.id) ? (editMode ? "Edit event" : "Event") : "Create event";

  let editMode: boolean;
  $: editMode = !(event && event.id);

  let eventCopy: EventModel;
  function startEditMode() {
    editMode = true;
    eventCopy = JSON.parse(JSON.stringify(event));
  }

  function cancelEdit() {
    editMode = false;
    event = eventCopy;
    event.date.start = new Date(event.date.start);
    event.date.end = new Date(event.date.start);
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
    <Vertical>
      <TextInput bind:value={event.name} name="name" placeholder="Name" editable={editMode} />
      <TextInput bind:value={event.desc} name="desc" placeholder="Description" multiline={true} editable={editMode} />
      <DateTimeInput bind:value={event.date.start} name="date_start" placeholder="Start" />
      <DateTimeInput bind:value={event.date.end} name="date_end" placeholder="End" />
      <Horizontal position="right">
        {#if editMode}
          <Button onClick={saveEdit} color="success">Save</Button>
          <Button onClick={cancelEdit} color="failure">Cancel</Button>
        {:else}
          <Button onClick={startEditMode} color="accent">Edit</Button>
          <Button onClick={showDeleteModal} color="failure">Delete</Button>
        {/if}
      </Horizontal>
    </Vertical>
  {/if}
</Modal>

<ConfirmationModal
  bind:showModal={showDeleteModal}
  confirmCallback={deleteEvent}
>
Do you really want to delete event "{event.name}"?
<br>
This action is irrevensible.
</ConfirmationModal>
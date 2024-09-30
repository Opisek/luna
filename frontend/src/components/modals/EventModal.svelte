<script lang="ts">
  import { deleteEvent, editEvent } from "$lib/client/repository";
  import DateTimeInput from "../forms/DateTimeInput.svelte";
  import TextInput from "../forms/TextInput.svelte";
  import EditableModal from "./EditableModal.svelte";

  export let event: EventModel;

  export let showCreateModal: () => boolean;
  export let showModal: () => boolean;

  let title: string;
  $: title = (event && event.id) ? (editMode ? "Edit event" : "Event") : "Create event";

  let editMode: boolean;

  const onDelete = async () => {
    const res = await deleteEvent(event.id);
    if (res === "") return "";
    else return `Could not delete event: ${res}`;
  };
  const onEdit = async () => {
    const res = await editEvent(event);
    if (res === "") return "";
    else return `Could not edit event: ${res}`;
  };
</script>

<EditableModal
  title={title}
  deleteConfirmation={`Are you sure you want to delete event "${event ? event.name : ""}"?`}
  bind:editMode={editMode}
  bind:showCreateModal={showCreateModal}
  bind:showModal={showModal}
  onDelete={onDelete}
  onEdit={onEdit}
>
  {#if event}
    <TextInput bind:value={event.name} name="name" placeholder="Name" editable={editMode} />
    <TextInput bind:value={event.desc} name="desc" placeholder="Description" multiline={true} editable={editMode} />
    <DateTimeInput bind:value={event.date.start} name="date_start" placeholder="Start" />
    <DateTimeInput bind:value={event.date.end} name="date_end" placeholder="End" />
  {/if}
</EditableModal>
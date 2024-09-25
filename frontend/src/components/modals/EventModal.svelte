<script lang="ts">
  import DateTimeInput from "../forms/DateTimeInput.svelte";
  import TextInput from "../forms/TextInput.svelte";
  import EditableModal from "./EditableModal.svelte";

  export let event: EventModel;

  export let showModal: () => boolean;

  let title: string;
  $: title = (event && event.id) ? (editMode ? "Edit event" : "Event") : "Create event";

  let editMode: boolean;

  const onDelete = () => {};
  const onEdit = () => {};
</script>

<EditableModal
  title={title}
  deleteConfirmation={`Are you sure you want to delete event "${event ? event.name : ""}"?`}
  isNew={!(event && event.id)}
  bind:editMode={editMode}
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
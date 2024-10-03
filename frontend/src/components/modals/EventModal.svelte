<script lang="ts">
  import { createEvent, deleteEvent, editEvent, getCalendars } from "$lib/client/repository";
  import ColorInput from "../forms/ColorInput.svelte";
  import DateTimeInput from "../forms/DateTimeInput.svelte";
  import SelectInput from "../forms/SelectInput.svelte";
  import TextInput from "../forms/TextInput.svelte";
  import EditableModal from "./EditableModal.svelte";

  export let event: EventModel;

  let currentCalendars: CalendarModel[] = [];

  export const showCreateModal = () => {
    currentCalendars = getCalendars();
    showCreateModalInternal();
  }
  export const showModal = () => {
    showModalInternal();
  };

  let showCreateModalInternal: () => boolean;
  let showModalInternal: () => boolean;

  let title: string;
  $: title = (event && event.id) ? (editMode ? "Edit event" : "Event") : "Create event";

  let editMode: boolean;

  const onDelete = async () => {
    const res = await deleteEvent(event.id);
    if (res === "") return "";
    else return `Could not delete event: ${res}`;
  };
  const onEdit = async () => {
    if (event.id === "") {
      const res = await createEvent(event);
      if (res === "") return "";
      else return `Could not edit event: ${res}`;
    } else {
      const res = await editEvent(event);
      if (res === "") return "";
      else return `Could not create event: ${res}`;
    }
  };
</script>

<EditableModal
  title={title}
  deleteConfirmation={`Are you sure you want to delete event "${event ? event.name : ""}"?`}
  bind:editMode={editMode}
  bind:showCreateModal={showCreateModalInternal}
  bind:showModal={showModalInternal}
  onDelete={onDelete}
  onEdit={onEdit}
>
  {#if event}
    <TextInput bind:value={event.name} name="name" placeholder="Name" editable={editMode} />
    {#if (event.id === "")}
      <SelectInput bind:value={event.calendar} name="calendar" placeholder="Calendar" options={currentCalendars.map(x => ({ value: x.id, name: x.name }))} editable={editMode} />
    {/if}
    <ColorInput bind:color={event.color} name="color" editable={editMode} />
    <TextInput bind:value={event.desc} name="desc" placeholder="Description" multiline={true} editable={editMode} />
    <DateTimeInput bind:value={event.date.start} name="date_start" placeholder="Start" />
    <DateTimeInput bind:value={event.date.end} name="date_end" placeholder="End" />
  {/if}
</EditableModal>
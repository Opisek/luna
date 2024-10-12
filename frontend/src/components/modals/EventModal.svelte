<script lang="ts">
  import { createEvent, deleteEvent, editEvent, getCalendars } from "$lib/client/repository";
  import CheckboxInput from "../forms/CheckboxInput.svelte";
  import ColorInput from "../forms/ColorInput.svelte";
  import DateTimeInput from "../forms/DateTimeInput.svelte";
  import SelectInput from "../forms/SelectInput.svelte";
  import TextInput from "../forms/TextInput.svelte";
  import EditableModal from "./EditableModal.svelte";

  export let event: EventModel;
  let eventCopy: EventModel;
  let lastStartDate: Date;

  let currentCalendars: CalendarModel[] = [];

  export const showCreateModal = () => {
    eventCopy = event;
    lastStartDate = eventCopy.date.start;
    currentCalendars = getCalendars();
    setTimeout(showCreateModalInternal, 0);
  }
  export const showModal = () => {
    eventCopy = { ...event };
    editMode = false;
    setTimeout(showModalInternal, 0);
  };

  let showCreateModalInternal: () => boolean;
  let showModalInternal: () => boolean;

  let title: string;
  $: title = (eventCopy && eventCopy.id) ? (editMode ? "Edit event" : "Event") : "Create event";

  let editMode: boolean;

  const onDelete = async () => {
    const res = await deleteEvent(eventCopy.id);
    if (res === "") return "";
    else return `Could not delete event: ${res}`;
  };
  const onEdit = async () => {
    if (eventCopy.id === "") {
      const res = await createEvent(eventCopy);
      if (res === "") return "";
      else return `Could not edit event: ${res}`;
    } else {
      const res = await editEvent(eventCopy);
      if (res === "") return "";
      else return `Could not create event: ${res}`;
    }
  };

  $: if (eventCopy && eventCopy.date && eventCopy.date.start && eventCopy.date.end && eventCopy.date.start.getTime() > eventCopy.date.end.getTime()) {
    if (eventCopy.date.start === lastStartDate) {
      eventCopy.date.start = new Date(eventCopy.date.end);
    } else {
      eventCopy.date.end = new Date(eventCopy.date.start);
      lastStartDate = eventCopy.date.start;
    }
  }
</script>

<EditableModal
  title={title}
  deleteConfirmation={`Are you sure you want to delete event "${eventCopy ? eventCopy.name : ""}"?`}
  bind:editMode={editMode}
  bind:showCreateModal={showCreateModalInternal}
  bind:showModal={showModalInternal}
  onDelete={onDelete}
  onEdit={onEdit}
>
  {#if eventCopy}
    <TextInput bind:value={eventCopy.name} name="name" placeholder="Name" editable={editMode} />
    {#if (eventCopy.id === "")}
      <SelectInput bind:value={eventCopy.calendar} name="calendar" placeholder="Calendar" options={currentCalendars.map(x => ({ value: x.id, name: x.name }))} editable={editMode} />
    {/if}
    {#if editMode}
      <ColorInput bind:color={eventCopy.color} name="color" editable={editMode} />
    {/if}
    <TextInput bind:value={eventCopy.desc} name="desc" placeholder="Description" multiline={true} editable={editMode} />
    {#if editMode}
        <CheckboxInput bind:value={eventCopy.date.allDay} name="all_day" description="All Day"/>
    {/if}
    <DateTimeInput bind:value={eventCopy.date.start} name="date_start" placeholder="Start" editable={editMode} allDay={eventCopy.date.allDay}/>
    <DateTimeInput bind:value={eventCopy.date.end} name="date_end" placeholder="End" editable={editMode} allDay={eventCopy.date.allDay}/>
  {/if}
</EditableModal>
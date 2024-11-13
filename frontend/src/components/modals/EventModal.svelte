<script lang="ts">
  import CheckboxInput from "../forms/CheckboxInput.svelte";
  import ColorInput from "../forms/ColorInput.svelte";
  import DateTimeInput from "../forms/DateTimeInput.svelte";
  import EditableModal from "./EditableModal.svelte";
  import SelectInput from "../forms/SelectInput.svelte";
  import TextInput from "../forms/TextInput.svelte";

  import { EmptyEvent } from "$lib/client/placeholders";
  import { createEvent, deleteEvent, editEvent, getCalendars } from "$lib/client/repository";

  interface Props {
    event: EventModel;
    showCreateModal?: () => any;
    showModal?: () => any;
  }

  let {
    event = $bindable(EmptyEvent),
    showCreateModal = $bindable(),
    showModal = $bindable(),
  }: Props = $props();

  let eventCopy: EventModel = $state(EmptyEvent);
  let currentCalendars: CalendarModel[] = $state([]);

  showCreateModal = () => {
    editMode = false;
    eventCopy = event;
    currentCalendars = getCalendars();
    setTimeout(showCreateModalInternal, 0);
  }

  showModal = () => {
    editMode = false;
    eventCopy = {
      id: event.id,
      calendar: event.calendar,
      name: event.name,
      desc: event.desc,
      color: event.color,
      date: {
        start: new Date(event.date.start),
        end: new Date(event.date.end),
        allDay: event.date.allDay,
      }
    }
    if (eventCopy.date.allDay) {
      eventCopy.date.end.setDate(eventCopy.date.end.getDate() - 1);
    }
    setTimeout(showModalInternal, 0);
  };

  let showCreateModalInternal: () => boolean = $state(() => false);
  let showModalInternal: () => boolean = $state(() => false);

  let editMode: boolean = $state(false);
  let title: string = $derived((eventCopy && eventCopy.id) ? (editMode ? "Edit event" : "Event") : "Create event");

  const onDelete = async () => {
    const res = await deleteEvent(eventCopy.id);
    if (res === "") return "";
    else return `Could not delete event: ${res}`;
  };
  const onEdit = async () => {
    if (eventCopy.date.allDay) {
      eventCopy.date.end.setDate(eventCopy.date.end.getDate() + 1);
    }
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

  const changeEnd = (value: Date) => {
    if (value.getTime() < eventCopy.date.start.getTime()) {
      const previousStart = eventCopy.date.start;
      eventCopy.date.start = new Date(value);

      if (Math.abs(previousStart.getTime() - value.getTime()) >= 24 * 60 * 60 * 1000) {
        eventCopy.date.start.setHours(previousStart.getHours(), previousStart.getMinutes(), previousStart.getSeconds(), previousStart.getMilliseconds());
      }
    }
  }

  const changeStart = (value: Date) => {
    if (value.getTime() > eventCopy.date.end.getTime()) {
      const previousEnd = eventCopy.date.end;
      eventCopy.date.end = new Date(value);

      if (Math.abs(previousEnd.getTime() - value.getTime()) >= 24 * 60 * 60 * 1000) {
        eventCopy.date.end.setHours(previousEnd.getHours(), previousEnd.getMinutes(), previousEnd.getSeconds(), previousEnd.getMilliseconds());
      }
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
    <DateTimeInput bind:value={eventCopy.date.start} name="date_start" placeholder="Start" editable={editMode} allDay={eventCopy.date.allDay} onChange={changeStart}/>
    <DateTimeInput bind:value={eventCopy.date.end} name="date_end" placeholder="End" editable={editMode} allDay={eventCopy.date.allDay} onChange={changeEnd}/>
  {/if}
</EditableModal>
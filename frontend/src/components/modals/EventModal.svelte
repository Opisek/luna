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
    showCreateModal?: (date: Date) => Promise<EventModel>;
    showModal?: (event: EventModel) => Promise<EventModel>;
  }

  let {
    showCreateModal = $bindable(),
    showModal = $bindable(),
  }: Props = $props();

  let event: EventModel = $state(EmptyEvent);
  let currentCalendars: CalendarModel[] = $state([]);

  let saveEvent = (_: EventModel | PromiseLike<EventModel>) => {};
  let cancelEvent = (_?: any) => {};

  showCreateModal = (date: Date) => {
    cancelEvent();
    
    editMode = false;

    const start = new Date(date);
    start.setHours(12, 0, 0, 0);

    const end = new Date(date);
    end.setHours(13, 0, 0, 0);

    event = {
      id: "",
      calendar: "",
      name: "",
      desc: "",
      color: "",
      date: {
        start: start,
        end: end,
        allDay: false,
      }
    };

    currentCalendars = getCalendars();
    setTimeout(showCreateModalInternal, 0);

    return new Promise((resolve, reject) => {
      saveEvent = resolve;
      cancelEvent = reject;
    });
  }

  showModal = (original: EventModel): Promise<EventModel> => {
    cancelEvent();

    editMode = false;
    event = {
      id: original.id,
      calendar: original.calendar,
      name: original.name,
      desc: original.desc,
      color: original.color,
      date: {
        start: new Date(original.date.start),
        end: new Date(original.date.end),
        allDay: original.date.allDay,
      }
    }
    if (event.date.allDay) {
      event.date.end.setDate(event.date.end.getDate() - 1);
    }
    setTimeout(showModalInternal, 0);

    return new Promise((resolve, reject) => {
      saveEvent = resolve;
      cancelEvent = reject;
    });
  };

  let showCreateModalInternal: () => boolean = $state(() => false);
  let showModalInternal: () => boolean = $state(() => false);

  let editMode: boolean = $state(false);
  let title: string = $derived((event && event.id) ? (editMode ? "Edit event" : "Event") : "Create event");

  const onDelete = async () => {
    await deleteEvent(event.id).catch(err => {
      throw new Error(`Could not delete event: ${err.message}`);
    });
    cancelEvent();
  };
  const onEdit = async () => {
    if (event.date.allDay) {
      event.date.end.setDate(event.date.end.getDate() + 1);
    }
    if (event.id === "") {
      await createEvent(event).catch(err => {
        cancelEvent();
        throw new Error(`Could not create event: ${err.message}`);
      });
      saveEvent(event);
    } else {
      await editEvent(event).catch(err => {
        cancelEvent();
        throw new Error(`Could not edit event: ${err.message}`);
      });
      saveEvent(event);
    }
  };

  const changeEnd = (value: Date) => {
    if (value.getTime() < event.date.start.getTime()) {
      const previousStart = event.date.start;
      event.date.start = new Date(value);

      if (Math.abs(previousStart.getTime() - value.getTime()) >= 24 * 60 * 60 * 1000) {
        event.date.start.setHours(previousStart.getHours(), previousStart.getMinutes(), previousStart.getSeconds(), previousStart.getMilliseconds());
      }
    }
  }

  const changeStart = (value: Date) => {
    if (value.getTime() > event.date.end.getTime()) {
      const previousEnd = event.date.end;
      event.date.end = new Date(value);

      if (Math.abs(previousEnd.getTime() - value.getTime()) >= 24 * 60 * 60 * 1000) {
        event.date.end.setHours(previousEnd.getHours(), previousEnd.getMinutes(), previousEnd.getSeconds(), previousEnd.getMilliseconds());
      }
    }
  }
  
</script>

<EditableModal
  title={title}
  deleteConfirmation={`Are you sure you want to delete event "${event ? event.name : ""}"?`}
  bind:editMode={editMode}
  bind:showCreateModal={showCreateModalInternal}
  bind:showModal={showModalInternal}
  onDelete={onDelete}
  onEdit={onEdit}
  submittable={event.calendar !== "" && event.name !== "" && event.date.start.getTime() < event.date.end.getTime()}
>
  {#if event != EmptyEvent}
    <TextInput bind:value={event.name} name="name" placeholder="Name" editable={editMode} />
    {#if (event.id === "")}
      <SelectInput bind:value={event.calendar} name="calendar" placeholder="Calendar" options={currentCalendars.map(x => ({ value: x.id, name: x.name }))} editable={editMode} />
    {/if}
    {#if editMode}
      <ColorInput bind:color={event.color} name="color" editable={editMode} />
    {/if}
    <TextInput bind:value={event.desc} name="desc" placeholder="Description" multiline={true} editable={editMode} />
    {#if editMode}
        <CheckboxInput bind:value={event.date.allDay} name="all_day" description="All Day"/>
    {/if}
    <DateTimeInput bind:value={event.date.start} name="date_start" placeholder="Start" editable={editMode} allDay={event.date.allDay} onChange={changeStart}/>
    <DateTimeInput bind:value={event.date.end} name="date_end" placeholder="End" editable={editMode} allDay={event.date.allDay} onChange={changeEnd}/>
  {/if}
</EditableModal>
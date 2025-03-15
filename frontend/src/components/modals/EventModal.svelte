<script lang="ts">
  import CheckboxInput from "../forms/CheckboxInput.svelte";
  import ColorInput from "../forms/ColorInput.svelte";
  import DateTimeInput from "../forms/DateTimeInput.svelte";
  import EditableModal from "./EditableModal.svelte";
  import SelectInput from "../forms/SelectInput.svelte";
  import TextInput from "../forms/TextInput.svelte";

  import { EmptyEvent } from "$lib/client/placeholders";
  import { getRepository } from "$lib/client/repository";
  import { deepCopy, deepEquality } from "$lib/common/misc";
  import { isSameDay } from "$lib/common/date";

  interface Props {
    showCreateModal?: (date: Date) => Promise<EventModel>;
    showModal?: (event: EventModel) => Promise<EventModel>;
  }

  let {
    showCreateModal = $bindable(),
    showModal = $bindable(),
  }: Props = $props();

  let event: EventModel = $state(EmptyEvent);
  let originalEvent: EventModel = $state(EmptyEvent);
  let eventSourceType = $state("");

  let saveEvent = (_: EventModel | PromiseLike<EventModel>) => {};
  let cancelEvent = (_?: any) => {};

  showCreateModal = async (date: Date) => {
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
        recurrence: false,
      }
    };

    setTimeout(showCreateModalInternal, 0);

    return new Promise((resolve, reject) => {
      saveEvent = resolve;
      cancelEvent = reject;
    });
  }

  showModal = async (original: EventModel): Promise<EventModel> => {
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
        recurrence: await deepCopy(original.date.recurrence),
      }
    }
    if (event.date.allDay && event.date.end.getTime() !== event.date.start.getTime() && event.date.end.getHours() === 0 && event.date.end.getMinutes() === 0 && event.date.end.getSeconds() === 0 && event.date.end.getMilliseconds() === 0) {
      event.date.end.setDate(event.date.end.getDate() - 1);
    }

    originalEvent = await deepCopy(original);
    const calendar = await getRepository().getCalendar(original.calendar).catch(err => {
      throw new Error(`Could not get calendar: ${err.message}`);
    });
    if (calendar) {
      const source = await getRepository().getSourceDetails(calendar.source).catch(err => {
        throw new Error(`Could not get source details: ${err.message}`);
      });
      if (source) {
        eventSourceType = source.type;
      }
    } else {
      eventSourceType = "";
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
  let editable: boolean = $derived(event && !(
    eventSourceType === "ical" || // iCal files are treated as read-only
    event.date.recurrence != false // for now we won't allow editing recurring events
  ))
  let showEndDate: boolean = $derived(editMode || (event && (!event.date.allDay || !isSameDay(event.date.start, event.date.end))));

  let selectableCalendars = $derived.by(() => {
    let calendars = getRepository().calendars.getArray();

    if (editMode) {
      const sources = getRepository().sources;

      calendars = calendars.filter(x => {
        if (x.id === event.calendar) return true;
        const source = sources.find("id", x.source);
        console.log(source);
        return source && source.type !== "ical";
      });
    } else {
      calendars = calendars.filter(x => x.id === event.calendar);
    }

    return calendars.map(x => ({ value: x.id, name: x.name }))
  });


  const onDelete = async () => {
    await getRepository().deleteEvent(event.id).catch(err => {
      throw new Error(`Could not delete event ${event.name}: ${err.message}`);
    });
    cancelEvent();
  };
  const onEdit = async () => {
    if (event.date.allDay) {
      event.date.end.setDate(event.date.end.getDate() + 1);
    }
    if (event.id === "") {
      await getRepository().createEvent(event).catch(err => {
        cancelEvent();
        throw new Error(`Could not create event ${event.name}: ${err.message}`);
      });
      saveEvent(event);
    } else if (event.calendar == originalEvent.calendar) {
      const changes = {
        name: event.name != originalEvent.name,
        desc: event.desc != originalEvent.desc,
        color: event.color != originalEvent.color,
        date: !deepEquality(event.date, originalEvent.date)
      };
      await getRepository().editEvent(event, changes).catch(err => {
        cancelEvent();
        throw new Error(`Could not edit event ${event.name}: ${err.message}`);
      });
      saveEvent(event);
    } else {
      await getRepository().moveEvent(event).catch(err => {
        cancelEvent();
        throw new Error(`Could not move event ${event.name}: ${err.message}`);
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
  editable={editable}
  submittable={event.calendar !== "" && event.name !== "" && (event.date.start.getTime() < event.date.end.getTime() || (event.date.start.getTime() <= event.date.end.getTime() && event.date.allDay))}
>
  {#if event != EmptyEvent}
    <TextInput bind:value={event.name} name="name" placeholder="Name" editable={editMode} />
    <SelectInput bind:value={event.calendar} name="calendar" placeholder="Calendar" options={selectableCalendars} editable={editMode} />
    {#if editMode}
      <ColorInput bind:color={event.color} name="color" editable={editMode} />
    {/if}
    {#if editMode || event.desc}
      <TextInput bind:value={event.desc} name="desc" placeholder="Description" multiline={true} editable={editMode} />
    {/if}
    {#if editMode}
        <CheckboxInput bind:value={event.date.allDay} name="all_day" description="All Day"/>
    {/if}
    <DateTimeInput bind:value={event.date.start} name="date_start" placeholder={showEndDate ? "Start" : "Date"} editable={editMode} allDay={event.date.allDay} onChange={changeStart}/>
    {#if showEndDate}
      <DateTimeInput bind:value={event.date.end} name="date_end" placeholder="End" editable={editMode} allDay={event.date.allDay} onChange={changeEnd}/>
    {/if}
  {/if}
</EditableModal>
<script lang="ts">
  import Button from "../interactive/Button.svelte";
  import ColorInput from "../forms/ColorInput.svelte";
  import DateTimeInput from "../forms/DateTimeInput.svelte";
  import EditableModal from "./EditableModal.svelte";
  import SelectInput from "../forms/SelectInput.svelte";
  import TextInput from "../forms/TextInput.svelte";

  import { EmptyEvent, NoChangesEvent, NoOp } from "$lib/client/placeholders";
  import { deepCopy, deepEquality } from "$lib/common/misc";
  import { getRepository } from "$lib/client/repository";
  import { isSameDay } from "$lib/common/date";
  import { queueNotification } from "$lib/client/notifications";
  import ToggleInput from "../forms/ToggleInput.svelte";
  import { getSettings } from "$lib/client/settings.svelte";
  import { UserSettingKeys } from "../../types/settings";

  interface Props {
    showCreateModal?: (date: Date) => Promise<EventModel>;
    showModal?: (event: EventModel) => Promise<EventModel>;
  }

  let {
    showCreateModal = $bindable(),
    showModal = $bindable(),
  }: Props = $props();

  const settings = getSettings();

  let event: EventModel = $state(EmptyEvent);
  let originalEvent: EventModel = $state(EmptyEvent);

  let calendars: CalendarModel[] = $state([]);
  let sources: SourceModel[] = $state([]);

  let eventSourceType = $derived.by(() => {
    const calendar = calendars.find(x => x.id === event.calendar);
    if (!calendar) return "unknown";

    const source = sources.find(x => x.id === calendar.source);
    if (!source) return "unknown";

    return source.type;
  });

  let promiseResolve: (value: EventModel | PromiseLike<EventModel>) => void = $state(NoOp);
  let promiseReject: (reason?: any) => void = $state(NoOp);

  showCreateModal = async (date: Date) => {
    promiseReject();

    calendars = getRepository().calendars.getArray();
    sources = getRepository().sources.getArray();
    
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
      },
      overridden: false
    };

    setTimeout(showCreateModalInternal, 0);

    return new Promise((resolve, reject) => {
      promiseResolve = resolve;
      promiseReject = reject;
    });
  }

  showModal = async (original: EventModel): Promise<EventModel> => {
    promiseReject();

    calendars = getRepository().calendars.getArray();
    sources = getRepository().sources.getArray();

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
      },
      overridden: original.overridden
    }
    if (event.date.allDay && event.date.end.getTime() !== event.date.start.getTime() && event.date.end.getHours() === 0 && event.date.end.getMinutes() === 0 && event.date.end.getSeconds() === 0 && event.date.end.getMilliseconds() === 0) {
      event.date.end.setDate(event.date.end.getDate() - 1);
    }

    originalEvent = await deepCopy(original);

    setTimeout(showModalInternal, 0);

    return new Promise((resolve, reject) => {
      promiseResolve = resolve;
      promiseReject = reject;
    });
  };

  let showCreateModalInternal: () => boolean = $state(() => false);
  let showModalInternal: () => boolean = $state(() => false);

  let editMode: boolean = $state(false);
  let title: string = $derived((event && event.id) ? (editMode ? "Edit event" : "Event") : "Create event");
  let showEndDate: boolean = $derived(editMode || (event && (!event.date.allDay || !isSameDay(event.date.start, event.date.end))));

  let selectableCalendars = $derived.by(() => {
    let selectable;

    if (editMode) {
      selectable = calendars.filter(x => {
        if (x.id === event.calendar) return true;
        const source = sources.find(y => y.id === x.source);
        return source && source.type !== "ical";
      });
    } else {
      selectable = calendars.filter(x => x.id === event.calendar);
    }

    return selectable.map(x => ({ value: x.id, name: x.name }))
  });


  const onDelete = async () => {
    await getRepository().deleteEvent(event.id).catch(err => {
      throw new Error(`Could not delete event ${event.name}: ${err.message}`);
    });
    promiseReject();
  };
  const onEdit = async () => {
    if (event.date.allDay) {
      event.date.end.setDate(event.date.end.getDate() + 1);
    }
    if (event.id === "") {
      await getRepository().createEvent(event).catch(err => {
        promiseReject();
        throw new Error(`Could not create event ${event.name}: ${err.message}`);
      });
      promiseResolve(event);
    } else if (event.calendar == originalEvent.calendar) {
      const changes = {
        name: event.name != originalEvent.name,
        desc: event.desc != originalEvent.desc,
        color: event.color != originalEvent.color,
        date: !deepEquality(event.date, originalEvent.date)
      };
      await getRepository().editEvent(event, changes, eventSourceType === "ical").catch(err => {
        promiseReject();
        throw new Error(`Could not edit event ${event.name}: ${err.message}`);
      });
      promiseResolve(event);
    } else {
      await getRepository().moveEvent(event).catch(err => {
        promiseReject();
        throw new Error(`Could not move event ${event.name}: ${err.message}`);
      });
      promiseResolve(event);
    }
  };
  const resetOverrides = async () => {
    event.overridden = false;
    getRepository().editEvent(event, NoChangesEvent, true).catch(err => {
      event.overridden = true;
      queueNotification("failure", `Could not reset event ${event.name}: ${err.message}`);
      return;
    }).then(async () => {
      getRepository().getEvent(event.id, true).catch(err => {
        event.overridden = true;
        queueNotification("failure", `Could not reset event ${event.name}: ${err.message}`);
        return;
      }).then((fetched) => {
        event = fetched as EventModel;
      });
    });
  }

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
  onCancel={promiseReject}
  deletable={event && eventSourceType !== "ical" && !event.date.recurrence}
  editable={event && !event.date.recurrence}
  submittable={event.calendar !== "" && event.name !== "" && (event.date.start.getTime() < event.date.end.getTime() || (event.date.start.getTime() <= event.date.end.getTime() && event.date.allDay))}
>
  {#if event != EmptyEvent}
    <TextInput bind:value={event.name} name="name" placeholder="Name" editable={editMode} />
    <SelectInput bind:value={event.calendar} name="calendar" placeholder="Calendar" options={selectableCalendars} editable={editMode && eventSourceType !== "ical"} />
    {#if editMode}
      <ColorInput bind:color={event.color} name="color" editable={editMode} />
    {/if}
    {#if editMode || event.desc}
      <TextInput bind:value={event.desc} name="desc" placeholder="Description" multiline={true} editable={editMode} />
    {/if}
    {#if editMode}
        <ToggleInput bind:value={event.date.allDay} name="all_day" description="All Day"/>
    {/if}
    <DateTimeInput bind:value={event.date.start} name="date_start" placeholder={showEndDate ? "Start" : "Date"} editable={editMode} allDay={event.date.allDay} onChange={changeStart}/>
    {#if showEndDate}
      <DateTimeInput bind:value={event.date.end} name="date_end" placeholder="End" editable={editMode} allDay={event.date.allDay} onChange={changeEnd}/>
    {/if}
    {#if settings.userSettings[UserSettingKeys.DebugMode]}
      <TextInput bind:value={event.id} name="id" placeholder="ID" editable={false} />
    {/if}
  {/if}
  {#snippet extraButtonsLeft()}
    {#if event != EmptyEvent && !editMode && event.overridden}
      <Button color="accent" onClick={resetOverrides}>Reset</Button>
    {/if}
  {/snippet}
</EditableModal>
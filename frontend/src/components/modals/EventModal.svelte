<script lang="ts">
  import Button from "../interactive/Button.svelte";
  import ColorInput from "../forms/ColorInput.svelte";
  import DateTimeInput from "../forms/DateTimeInput.svelte";
  import EditableModal from "./EditableModal.svelte";
  import SelectInput from "../forms/SelectInput.svelte";
  import TextInput from "../forms/TextInput.svelte";

  import { EmptyEvent, NoChangesEvent, NoOp } from "$lib/client/placeholders";
  import { deepCopy, deepEquality } from "$lib/common/misc";
  import { getRepository } from "$lib/client/data/repository.svelte";
  import { isSameDay } from "$lib/common/date";
  import { queueNotification } from "$lib/client/notifications";
  import ToggleInput from "../forms/ToggleInput.svelte";
  import { getSettings } from "$lib/client/data/settings.svelte";
  import { UserSettingKeys } from "../../types/settings";
  import { ColorKeys } from "../../types/colors";
  import Horizontal from "../layout/Horizontal.svelte";
  import EventCopyModal from "./EventCopyModal.svelte";
  import IconButton from "../interactive/IconButton.svelte";
  import { Copy } from "lucide-svelte";

  //import { RRule } from "rrule";

  interface Props {
    showModal?: (initial?: EventModel, date?: Date) => Promise<EventModel>;
  }

  let {
    showModal = $bindable(),
  }: Props = $props();

  const settings = getSettings();
  const repository = getRepository();

  let showModalInternal: (initial?: EventModel, edit?: boolean) => Promise<EventModel> = $state(Promise.reject);
  let showCopyModal: (event: EventModel) => Promise<boolean> = $state(async () => false);
  let editMode: boolean = $state(false);

  let event: EventModel = $state(EmptyEvent);
  let originalEvent: EventModel = $state(EmptyEvent);
  let eventRecurrenceBoolean = $state(false);
  //let eventRecurrenceObject: RRule | null = $state(null);

  let eventSourceType = $derived.by(() => {
    const calendar = repository.calendars.find(x => x.id === event.calendar);
    if (!calendar) return "unknown";

    const source = repository.sources.find(x => x.id === calendar.source);
    if (!source) return "unknown";

    return source.type;
  });

  showModal = async (initial?: EventModel, date?: Date): Promise<EventModel> => {
    if (!initial) {
      const start = new Date(date || new Date());
      start.setHours(12, 0, 0, 0);

      const end = new Date(date || new Date());
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
        overridden: false,
        can_edit: true,
        can_delete: true,
      };

      eventRecurrenceBoolean = false;
      //eventRecurrenceObject = null;
    } else {
      event = {
        id: initial.id,
        calendar: initial.calendar,
        name: initial.name,
        desc: initial.desc,
        color: initial.color,
        date: {
          start: new Date(initial.date.start),
          end: new Date(initial.date.end),
          allDay: initial.date.allDay,
          recurrence: await deepCopy(initial.date.recurrence),
        },
        overridden: initial.overridden,
        can_edit: initial.can_edit,
        can_delete: initial.can_delete,
      }
      if (event.date.allDay && event.date.end.getTime() !== event.date.start.getTime() && event.date.end.getHours() === 0 && event.date.end.getMinutes() === 0 && event.date.end.getSeconds() === 0 && event.date.end.getMilliseconds() === 0) {
        event.date.end.setDate(event.date.end.getDate() - 1);
      }

      originalEvent = await deepCopy(initial);

      eventRecurrenceBoolean = event.date.recurrence != false;
      //eventRecurrenceObject = eventRecurrenceBoolean ? RRule.fromString(event.date.recurrence) : null;
    }

    return showModalInternal(event);
  };

  let title: string = $derived((event && event.id) ? (editMode ? "Edit event" : "Event") : "Create event");
  let showEndDate: boolean = $derived(editMode || (event && (!event.date.allDay || !isSameDay(event.date.start, event.date.end))));

  let selectableCalendars = $derived(
    repository.calendars
      .filter(calendar => calendar.id === event.calendar || (editMode && calendar.can_add_events))
      .map(calendar => ({ value: calendar.id, name: calendar.name }))
  );

  const onDelete = async () => {
    return await getRepository().deleteEvent(event.id).then(() => event).catch(err => {
      throw new Error(`Could not delete event ${event.name}: ${err.message}`);
    });
  };
  const onEdit = async () => {
    if (event.date.allDay) {
      event.date.end.setDate(event.date.end.getDate() + 1);
    }
    if (event.id === "") {
      return await getRepository().createEvent(event).then(() => event).catch(err => {
        throw new Error(`Could not create event ${event.name}: ${err.message}`);
      });
    } else if (event.calendar == originalEvent.calendar) {
      const changes = {
        name: event.name != originalEvent.name,
        desc: event.desc != originalEvent.desc,
        color: event.color != originalEvent.color,
        date: !deepEquality(event.date, originalEvent.date)
      };
      return await getRepository().editEvent(event, changes, eventSourceType === "ical").then(() => event).catch(err => {
        throw new Error(`Could not edit event ${event.name}: ${err.message}`);
      });
    } else {
      return await getRepository().moveEvent(event).then(() => event).catch(err => {
        throw new Error(`Could not move event ${event.name}: ${err.message}`);
      });
    }
  };
  const resetOverrides = async () => {
    event.overridden = false;
    getRepository().editEvent(event, NoChangesEvent, true).catch(err => {
      event.overridden = true;
      queueNotification(ColorKeys.Danger, `Could not reset event ${event.name}: ${err.message}`);
      return;
    }).then(async () => {
      getRepository().getEvent(event.id, true).catch(err => {
        event.overridden = true;
        queueNotification(ColorKeys.Danger, `Could not reset event ${event.name}: ${err.message}`);
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
  
  const copyEvent = async () => {
    showCopyModal(originalEvent).then(() => {

    }).catch(() => {

    });
  }
</script>

<EditableModal
  title={title}
  deleteConfirmation={`Are you sure you want to delete event "${event ? event.name : ""}"?`}
  bind:editMode={editMode}
  bind:showModal={showModalInternal}
  onDelete={onDelete}
  onEdit={onEdit}
  deletable={event?.can_edit}
  editable={event?.can_delete}
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
    <Horizontal position="left">
      <DateTimeInput bind:value={event.date.start} name="date_start" placeholder={showEndDate ? "Start" : "Date"} editable={editMode} allDay={event.date.allDay} onChange={changeStart} wrap={true}/>
      {#if showEndDate}
        <DateTimeInput bind:value={event.date.end} name="date_end" placeholder="End" editable={editMode} allDay={event.date.allDay} onChange={changeEnd} wrap={true}/>
      {/if}
    </Horizontal>
    {#if event.id && settings.userSettings[UserSettingKeys.DebugMode]}
      <TextInput value={event.id} name="id" placeholder="Event ID" editable={false} />
    {/if}
    {#if editMode}
      <ToggleInput bind:value={eventRecurrenceBoolean} name="repeats" description="Repeats"/>
    {/if}
    {#if eventRecurrenceBoolean}
      {#if editMode}
        <SelectInput bind:value={event.date.recurrence} name="recurrence_freq" placeholder="Frequency" showLabel={true} options={[
          { value: "daily", name: "Daily" },
          { value: "weekly", name: "Weekly" },
          { value: "monthly", name: "Monthly" },
          { value: "yearly", name: "Yearly" },
        ]} editable={editMode} />
      {:else}
        Repeats xyz times or something
      {/if}
    {/if}
  {/if}
  {#snippet extraButtonsLeft()}
    {#if !editMode}
      {#if event != EmptyEvent && event.overridden}
        <Button color={ColorKeys.Accent} onClick={resetOverrides}>Reset</Button>
      {/if}
      <IconButton onClick={copyEvent} alt="Copy" canRenderAsButton={true}>
        <Copy/>
      </IconButton>
    {/if}
  {/snippet}
</EditableModal>

<EventCopyModal bind:copy={showCopyModal}/>
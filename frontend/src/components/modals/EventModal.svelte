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
  import { Frequency, RRule, type Options } from "rrule";
  import { parseTimestampList, serializeTimestampList } from "../../lib/common/ical";
  import { SvelteSet } from "svelte/reactivity";
  import AffectedRecurrencesModal from "./AffectedRecurrencesModal.svelte";
  import { t } from "@sveltia/i18n";
  import RecurrenceInput from "../forms/RecurrenceInput.svelte";
  import Link from "../forms/Link.svelte";
  import RecurrenceRuleModal from "./RecurrenceRuleModal.svelte";
  import Title from "../layout/Title.svelte";

  interface Props {
    showModal?: (initial?: EventModel, date?: Date, anchor?: HTMLElement) => Promise<EventModel>;
  }

  let {
    showModal = $bindable(),
  }: Props = $props();

  const settings = getSettings();
  const repository = getRepository();

  let showModalInternal: (initial?: EventModel, edit?: boolean, anchor?: HTMLElement) => Promise<EventModel> = $state(Promise.reject);
  let showCopyModal: (event: EventModel) => Promise<EventModel> = $state(Promise.reject);
  let showRecurrenceRuleModal: (initial: Partial<Options>) => Promise<Partial<Options>> = $state(Promise.reject);
  let selectAffectedRecurrences: (edit: boolean) => Promise<"this" | "thisandfuture" | "all"> = $state(Promise.reject);
  let editMode: boolean = $state(false);

  let event: EventModel = $state(EmptyEvent);
  let originalEvent: EventModel = $state(EmptyEvent);

  let eventRepeats = $state(false);
  let eventRecurrenceRrule = $state(new RRule());
  let eventRecurrenceRruleOptions = $state<Partial<Options>>({});
  let eventRecurrenceRdate = $state(new SvelteSet<Date>());
  let eventRecurrenceExdate = $state(new SvelteSet<Date>());
  $effect(() => {
    eventRecurrenceRruleOptions.dtstart = event.date.start;
  });

  let eventSourceType = $derived.by(() => {
    const calendar = repository.calendars.find(x => x.id === event.calendar);
    if (!calendar) return "unknown";

    const source = repository.sources.find(x => x.id === calendar.source);
    if (!source) return "unknown";

    return source.type;
  });

  showModal = async (initial?: EventModel, date?: Date, anchor?: HTMLElement): Promise<EventModel> => {
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
          recurrence: undefined,
        },
        overridden: false,
        can_edit: true,
        can_delete: true,
      };

      eventRepeats = false;
      //eventRecurrenceObject = null;

      eventRecurrenceRrule = new RRule({ dtstart: start, freq: Frequency.YEARLY });
      eventRecurrenceRruleOptions = eventRecurrenceRrule.origOptions;
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

      eventRepeats = event.date.recurrence != undefined;
      if (event.date.recurrence) {
        if (event.date.recurrence.RRULE) {
          eventRecurrenceRrule = RRule.fromString(event.date.recurrence.RRULE);
          eventRecurrenceRrule.origOptions.dtstart = event.date.start;
          eventRecurrenceRrule = new RRule(eventRecurrenceRrule.origOptions);
          eventRecurrenceRruleOptions = eventRecurrenceRrule.origOptions;
        }

        if (event.date.recurrence.RDATE)
          eventRecurrenceRdate = new SvelteSet(parseTimestampList(event.date.recurrence.RDATE));

        if (event.date.recurrence.EXDATE)
          eventRecurrenceExdate = new SvelteSet(parseTimestampList(event.date.recurrence.EXDATE));
      }
    }

    return showModalInternal(event, false, anchor);
  };

  let title: string = $derived((event && event.id) ? (editMode ? t("event.title.edit") : t("event.title.view")) : t("event.title.create"));
  let showEndDate: boolean = $derived(editMode || (event && (!event.date.allDay || !isSameDay(event.date.start, event.date.end))));

  let selectableCalendars = $derived(
    repository.calendars
      .filter(calendar => calendar.id === event.calendar || (editMode && calendar.can_add_events))
      .map(calendar => ({ value: calendar.id, name: calendar.name }))
  );

  const onDelete = async () => {
    const affect = originalEvent.date.recurrence != undefined ? await selectAffectedRecurrences(false).catch(() => { throw new Error("Cancelled"); }) : "this";

    return await getRepository().deleteEvent(event.id, affect).then(() => event).catch(err => {
      throw new Error(t("event.error.delete", { values: { name: event.name, msg: err.message } }));
    });
  };
  const onEdit = async () => {
    const affect = (event.id !== "" && originalEvent.date.recurrence != undefined) ? await selectAffectedRecurrences(true).catch(() => { throw new Error("Cancelled"); }) : "this";

    if (eventRepeats) {
      eventRecurrenceRruleOptions.dtstart = event.date.start;

      event.date.recurrence = {
        RRULE: `RRULE:${RRule.optionsToString(new RRule(eventRecurrenceRruleOptions).options).split("RRULE:")[1]}`,
        RDATE: serializeTimestampList("RDATE", event.date.allDay, "UTC", [...eventRecurrenceRdate.values()]),
        EXDATE: serializeTimestampList("EXDATE", event.date.allDay, "UTC", [...eventRecurrenceExdate.values()]),
      };
      if (event.date.recurrence.RDATE === undefined) delete event.date.recurrence.RDATE;
      if (event.date.recurrence.EXDATE === undefined) delete event.date.recurrence.EXDATE;
    } else {
      event.date.recurrence = undefined;
    }

    if (event.date.allDay) {
      event.date.end.setDate(event.date.end.getDate() + 1);
    }
    if (event.id === "") {
      return await getRepository().createEvent(event).then(() => event).catch(err => {
        throw new Error(t("event.error.create", { values: { name: event.name, msg: err.message } }));
      });
    } else if (event.calendar == originalEvent.calendar) {
      const changes = {
        name: event.name != originalEvent.name,
        desc: event.desc != originalEvent.desc,
        color: event.color != originalEvent.color,
        date: !deepEquality(event.date, originalEvent.date) // TODO: this needs to be fixed!
      };
      return await getRepository().editEvent(event, changes, eventSourceType === "ical", affect).then(() => event).catch(err => {
        throw new Error(t("event.error.edit", { values: { name: event.name, msg: err.message } }));
      });
    } else {
      return await getRepository().moveEvent(event).then(() => event).catch(err => {
        throw new Error(t("event.error.move", { values: { name: event.name, msg: err.message } }));
      });
    }
    // TODO: after editing/deleting recurring events with affect != "this", we have to refetch from this calendar, because the frontend does not try to make assumptions about what other events are affected
  };
  const resetOverrides = async () => {
    event.overridden = false;
    getRepository().editEvent(event, NoChangesEvent, true).catch(err => {
      event.overridden = true;
      queueNotification(ColorKeys.Danger, t("event.error.reset", { values: { name: event.name, msg: err.message } }));
      return;
    }).then(async () => {
      getRepository().getEvent(event.id, true).catch(err => {
        event.overridden = true;
        queueNotification(ColorKeys.Danger, t("event.error.reset", { values: { name: event.name, msg: err.message } }));
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
    await showCopyModal(originalEvent).then((newEvent) => {
      console.log(originalEvent);
      console.log(newEvent);
    }).catch(() => {
      console.log("aborted copy")
    });
  }
</script>

<EditableModal
  title={title}
  deleteConfirmation={t("event.confirm.delete", { values: { name: event.name } })}
  bind:editMode={editMode}
  bind:showModal={showModalInternal}
  onDelete={onDelete}
  onEdit={onEdit}
  deletable={event?.can_delete}
  editable={event?.can_edit}
  submittable={event.calendar !== "" && event.name !== "" && (event.date.start.getTime() < event.date.end.getTime() || (event.date.start.getTime() <= event.date.end.getTime() && event.date.allDay))}
>
  {#if event != EmptyEvent}
    {#if editMode}
      <TextInput bind:value={event.name} name="name" placeholder={t("form.name")} editable={true} />
    {:else}
      <Title>{event.name}</Title>
    {/if}
    <SelectInput bind:value={event.calendar} name="calendar" placeholder={t("calendar.display")} options={selectableCalendars} editable={editMode && eventSourceType !== "ical"} />
    {#if editMode}
      <ColorInput bind:color={event.color} name="color" editable={editMode} />
    {/if}
    {#if editMode || event.desc}
      <TextInput bind:value={event.desc} name="desc" placeholder={t("form.desc")} multiline={true} editable={editMode} />
    {/if}
    {#if editMode}
      <ToggleInput bind:value={event.date.allDay} name="all_day" description={t("date.allDay")}/>
    {/if}
    <Horizontal position="left">
      <DateTimeInput bind:value={event.date.start} name="date_start" placeholder={showEndDate ? t("date.start") : t("date.date")} editable={editMode} allDay={event.date.allDay} onChange={changeStart} wrap={true}/>
      {#if showEndDate}
        <DateTimeInput bind:value={event.date.end} name="date_end" placeholder={t("date.end")} editable={editMode} allDay={event.date.allDay} onChange={changeEnd} wrap={true}/>
      {/if}
    </Horizontal>
    {#if editMode}
      <ToggleInput bind:value={eventRepeats} name="repeats" description={t("recurrence.repeats")}/>
    {/if}
    {#if eventRepeats}
      {#if editMode}
        <RecurrenceInput
          bind:options={eventRecurrenceRruleOptions} 
          dtstart={event.date.start}
          allDay={event.date.allDay}
          editable={editMode}
          simple={true}
        />
        <Horizontal position="right">
          <Link onClick={async () => await showRecurrenceRuleModal(eventRecurrenceRruleOptions).then(x => eventRecurrenceRruleOptions = x).catch(NoOp)}>Advanced recurrence editing</Link>
        </Horizontal>
        <RecurrenceRuleModal
          dtstart={event.date.start}
          allDay={event.date.allDay}
          bind:showModal={showRecurrenceRuleModal}
        />
      {/if}
    {/if}
    {#if event.id && settings.userSettings[UserSettingKeys.DebugMode]}
      <TextInput value={event.id} name="id" placeholder={t("event.id")} editable={false} />
    {/if}
  {/if}
  {#snippet extraButtonsLeft()}
    {#if !editMode}
      {#if event != EmptyEvent && event.overridden}
        <Button color={ColorKeys.Accent} onClick={resetOverrides}>{t("button.reset")}</Button>
      {/if}
      <IconButton onClick={copyEvent} alt={t("button.copy")} canRenderAsButton={true}>
        <Copy/>
      </IconButton>
    {/if}
  {/snippet}
</EditableModal>

<EventCopyModal bind:copy={showCopyModal}/>
<AffectedRecurrencesModal bind:showModal={selectAffectedRecurrences}/>
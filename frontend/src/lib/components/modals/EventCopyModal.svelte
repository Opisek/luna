<script lang="ts">
  import Modal from "./Modal.svelte";

  import { EmptyEvent, NoOp } from "$lib/scripts/client/placeholders";
  import Paragraph from "../forms/Paragraph.svelte";
  import { getRepository } from "$lib/scripts/client/data/repository.svelte";
  import SmallCalendar from "../interactive/SmallCalendar.svelte";
  import MonthSelection from "../interactive/MonthSelection.svelte";
  import { SvelteMap } from "svelte/reactivity";
  import { ColorKeys } from "$lib/types/colors";
  import SelectInput from "../forms/SelectInput.svelte";
  import { deepCopy } from "$lib/scripts/common/misc";
  import IconButton from "../interactive/IconButton.svelte";
  import { Check, X } from "lucide-svelte";
  import { parseTimestampList, serializeTimestampList } from "$lib/scripts/common/ical";
  import { t } from "@sveltia/i18n";
  import ToggleInput from "../forms/ToggleInput.svelte";
  
  interface Props {
    copy?: (event: EventModel) => Promise<EventModel>;
  }

  let {
    copy = $bindable(),
  }: Props = $props();

  let showModalInternal: () => Promise<EventModel> = $state(Promise.reject);
  let success: (result: EventModel) => void = $state(NoOp);
  let failure: () => void = $state(NoOp);

  const repository = getRepository();

  let date = $state(new Date());
  let marked: Map<string, Date> = $state(new SvelteMap());
  let copyAsRecurrence = $state(false);
  
  let event = $state(EmptyEvent)
  let original = $state(EmptyEvent)

  let selectableCalendars = $derived(
    repository.calendars
      .filter(calendar => calendar.id === event.calendar || calendar.can_add_events)
      .map(calendar => ({ value: calendar.id, name: calendar.name }))
  );

  copy = async (eventToCopy: EventModel) => {
    event = await deepCopy(eventToCopy);
    original = eventToCopy;

    date = new Date(event.date.start);
    if (event.date.allDay) {
      date.setUTCMilliseconds(0);
      date.setUTCSeconds(0);
      date.setUTCMinutes(0);
      date.setUTCHours(0);
      date.setUTCDate(event.date.start.getDate());
      date.setUTCMonth(event.date.start.getMonth());
      date.setUTCFullYear(event.date.start.getFullYear());
    }
    marked = new SvelteMap();
    copyAsRecurrence = false;

    return showModalInternal();
  }

  function markDay(day: Date) {
    marked.set(day.toISOString().substring(0, 10), day);
  }

  function unmarkDay(day: Date) {
    marked.delete(day.toISOString().substring(0, 10));
  }

  function isMarked(day: Date): boolean {
    return marked.has(day.toISOString().substring(0, 10));
  }

  function daySelected(day: Date) {
    if (!original.date.allDay) {
      day.setHours(date.getHours());
      day.setMinutes(date.getMinutes());
      day.setSeconds(date.getSeconds());
    }
    if (isMarked(day)) unmarkDay(day);
    else markDay(day);
  }

  async function save() {
    if (copyAsRecurrence) {
      if (event.date.recurrence && event.date.recurrence.RDATE) {
        parseTimestampList(event.date.recurrence.RDATE).forEach(markDay);
      } else {
        markDay(date);
      }

      if (!event.date.recurrence) event.date.recurrence = {};
      event.date.recurrence.RDATE = serializeTimestampList("RDATE", event.date.allDay, "UTC", [...marked.values()]);
      return getRepository().editEvent(event, { date: true }, false).then(() => success(event)).catch(err => {
        throw new Error(t("event.error.copy", { values: { name: event.name, msg: err.message } }));
      });
    } else {
      const newDays = [...marked.values()]
      newDays.forEach(newDay => {
        newDay.setHours(date.getHours());
        newDay.setMinutes(date.getMinutes());
        newDay.setSeconds(date.getSeconds());
      })

      const delta = original.date.end.getTime() - original.date.start.getTime();

      const promises = newDays.map(async newDate => new Promise<EventModel>(async (res, rej) => {
        const newEvent = await deepCopy(original);
        newEvent.date.start = newDate;
        newEvent.date.end = new Date(newEvent.date.start.getTime() + delta);
        newEvent.calendar = event.calendar;
        return getRepository().createEvent(newEvent).then(() => res(event)).catch(err => {
          rej(t("event.error.copy", { values: { name: event.name, msg: err.message } }));
        });
      }))
      await Promise.all(promises).then(x => success(event)).catch(err => { throw new Error(err) });
    }
  }
</script>

<style lang="scss">
  @use "$lib/styles/colors.scss";
  @use "$lib/styles/dimensions.scss";
  @use "$lib/styles/text.scss";
</style>

<Modal
  title={t("event.title.copy")}
  bind:showModal={showModalInternal}
  bind:success
  bind:failure
>
  <SelectInput bind:value={event.calendar} name="calendar" placeholder={t("calendar.display")} options={selectableCalendars} />

  <Paragraph>
    {t("event.copy.instruction", { values: { count: marked.size } })}
  </Paragraph>

  <MonthSelection bind:date />
  <SmallCalendar bind:date bind:marked onDayClick={daySelected} />

  <ToggleInput bind:value={copyAsRecurrence} name="copy-rdate" description={t("event.copy.rdate.toggle")} info={t("event.copy.rdate.details", { values: { count: marked.size } })}/>

  {#snippet buttons()}
    <IconButton onClick={save} color={ColorKeys.Success} enabled={marked.size != 0} type="submit" alt={t("button.save")} canRenderAsButton={true}><Check/></IconButton>
    <IconButton onClick={failure} color={ColorKeys.Danger} alt={t("button.cancel")} canRenderAsButton={true}><X/></IconButton>
  {/snippet}
</Modal>
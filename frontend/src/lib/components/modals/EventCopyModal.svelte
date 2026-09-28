<script lang="ts">
  import Modal from "./Modal.svelte";

  import { EmptyEvent, NoOp } from "$lib/scripts/client/placeholders";
  import { getSettings } from "$lib/scripts/client/data/settings.svelte";
  import Paragraph from "../forms/Paragraph.svelte";
  import { getRepository } from "$lib/scripts/client/data/repository.svelte";
  import SmallCalendar from "../interactive/SmallCalendar.svelte";
  import MonthSelection from "../interactive/MonthSelection.svelte";
  import { SvelteMap, SvelteSet } from "svelte/reactivity";
  import Button from "../interactive/Button.svelte";
  import { ColorKeys } from "$lib/types/colors";
  import SelectInput from "../forms/SelectInput.svelte";
  import { deepCopy } from "$lib/scripts/common/misc";
  import IconButton from "../interactive/IconButton.svelte";
  import { Check, Save, X } from "lucide-svelte";
  import { parseTimestampList, serializeTimestampList } from "$lib/scripts/common/ical";
  import { t } from "@sveltia/i18n";
  
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
    marked = new SvelteMap();
    if (event.date.recurrence && event.date.recurrence.RDATE) {
      parseTimestampList(event.date.recurrence.RDATE).forEach(markDay);
    } else {
      markDay(date);
    }

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
    day.setHours(date.getHours());
    day.setMinutes(date.getMinutes());
    day.setSeconds(date.getSeconds());
    if (isMarked(day)) unmarkDay(day);
    else markDay(day);
  }

  async function save() {
    if (!event.date.recurrence) event.date.recurrence = {};
    event.date.recurrence.RDATE = serializeTimestampList("RDATE", event.date.allDay, "UTC", [...marked.values()]);
    return getRepository().editEvent(event, { date: true }, false).then(() => success(event)).catch(err => {
      throw new Error(t("event.error.copy", { values: { name: event.name, msg: err.message } }));
    });
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
    {t("event.copy.instruction")}
  </Paragraph>

  <MonthSelection bind:date />
  <SmallCalendar bind:date bind:marked onDayClick={daySelected} />

  {#snippet buttons()}
    <IconButton onClick={save} color={ColorKeys.Success} enabled={marked.size != 0} type="submit" alt={t("button.save")} canRenderAsButton={true}><Check/></IconButton>
    <IconButton onClick={failure} color={ColorKeys.Danger} alt={t("button.cancel")} canRenderAsButton={true}><X/></IconButton>
  {/snippet}
</Modal>
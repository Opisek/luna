<script lang="ts">
  import Modal from "./Modal.svelte";

  import { EmptyEvent, NoOp } from "$lib/client/placeholders";
  import { getSettings } from "$lib/client/data/settings.svelte";
  import Paragraph from "../forms/Paragraph.svelte";
  import { getRepository } from "../../lib/client/data/repository.svelte";
  import SmallCalendar from "../interactive/SmallCalendar.svelte";
  import MonthSelection from "../interactive/MonthSelection.svelte";
  import { SvelteSet } from "svelte/reactivity";
  import Button from "../interactive/Button.svelte";
  import { ColorKeys } from "../../types/colors";
  import SelectInput from "../forms/SelectInput.svelte";
  import { deepCopy } from "../../lib/common/misc";
  import IconButton from "../interactive/IconButton.svelte";
  import { Check, Save, X } from "lucide-svelte";
  
  interface Props {
    copy?: (event: EventModel) => Promise<boolean>;
  }

  let {
    copy = $bindable(),
  }: Props = $props();

  const settings = getSettings();
  const repository = getRepository();

  let date = $state(new Date());
  let marked: Set<string> = $state(new Set());
  
  let event = $state(EmptyEvent)
  let original = $state(EmptyEvent)

  let selectableCalendars = $derived(
    repository.calendars
      .filter(calendar => calendar.id === event.calendar || calendar.can_add_events)
      .map(calendar => ({ value: calendar.id, name: calendar.name }))
  );

  let showModalInternal = $state(NoOp);
  let hideModalInternal = $state(NoOp);

  let promiseResolve: (copied: boolean | PromiseLike<boolean>) => void = $state(NoOp);
  let promiseReject: (reason?: any) => void = $state(NoOp);

  copy = async (eventToCopy: EventModel) => {
    promiseReject();

    event = await deepCopy(eventToCopy);
    original = eventToCopy;
    date = new Date(event.date.start);
    //marked = new SvelteSet([date.toISOString().substring(0, 10)]);
    marked = new SvelteSet();

    showModalInternal();

    return new Promise((resolve, reject) => {
      promiseResolve = ((res) => {
        resolve(res);
      });
      promiseReject = ((err) => {
        reject(err);
      });
    })
  }

  function daySelected(day: Date) {
    const isoDay = day.toISOString().substring(0, 10);
    if (marked.has(isoDay)) marked.delete(isoDay);
    else marked.add(isoDay);
  }

  async function save() {
  }
</script>

<style lang="scss">
  @use "../../styles/colors.scss";
  @use "../../styles/dimensions.scss";
  @use "../../styles/text.scss";
</style>

<Modal
  title={"Copy Event"}
  bind:showModal={showModalInternal}
  bind:hideModal={hideModalInternal}
  onModalHide={() => {
    promiseReject();
  }}
>
  <SelectInput bind:value={event.calendar} name="calendar" placeholder="Calendar" options={selectableCalendars} />

  <Paragraph>
    Select on which days the event should take place.
  </Paragraph>

  <MonthSelection bind:date />
  <SmallCalendar bind:date bind:marked onDayClick={daySelected} />

  {#snippet buttons()}
    <IconButton onClick={NoOp} color={ColorKeys.Success} enabled={marked.size != 0} type="submit" alt="Save" canRenderAsButton={true}><Check/></IconButton>
    <IconButton onClick={hideModalInternal} color={ColorKeys.Danger} alt="Cancel" canRenderAsButton={true}><X/></IconButton>
  {/snippet}
</Modal>
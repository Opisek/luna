<script lang="ts">
  import EditableModal from "./EditableModal.svelte";
  import TextInput from "../forms/TextInput.svelte";

  import { EmptyCalendar, NoOp } from "$lib/client/placeholders";
  import { getRepository } from "$lib/client/repository";
  import { deepCopy } from "$lib/common/misc";
  import SelectInput from "../forms/SelectInput.svelte";
  import ColorInput from "../forms/ColorInput.svelte";

  interface Props {
    showCreateModal?: () => any;
    showModal?: (calendar: CalendarModel) => Promise<CalendarModel>;
  }

  let {
    showCreateModal = $bindable(),
    showModal = $bindable(),
  }: Props = $props();

  let calendar: CalendarModel = $state(EmptyCalendar);
  let originalCalendar: CalendarModel = $state(EmptyCalendar);

  let saveCalendar = (_: CalendarModel | PromiseLike<CalendarModel>) => {};
  let cancelCalendar = (_?: any) => {};

  showCreateModal = async () => {
    cancelCalendar();

    calendar = {
      id: "",
      source: "",
      name: "",
      desc: "",
      color: ""
    };

    showCreateModalInternal();
  }
  showModal = async (original: CalendarModel): Promise<CalendarModel> => {
    cancelCalendar();

    originalCalendar = await deepCopy(original);
    calendar = original;

    showModalInternal();

    return new Promise((resolve, reject) => {
      saveCalendar = resolve;
      cancelCalendar = reject;
    })
  };

  let showCreateModalInternal: () => any = $state(NoOp);
  let showModalInternal: () => any = $state(NoOp);

  let editMode: boolean = $state(false);
  let title: string = $derived(calendar.id ? (editMode ? "Edit calendar" : "Calendar") : "Add calendar");

  let selectableSources = $derived(
    getRepository().sources.getArray()
      .filter(x => editMode ? x.type !== "ical" || x.id === calendar.source : x.id === calendar.source)
      .map(x => ({ value: x.id, name: x.name }))
  );

  const onDelete = async () => {
    await getRepository().deleteCalendar(calendar.id).catch(err => {
      throw new Error(`Could not delete calendar ${calendar.name}: ${err.message}`);
    });
    cancelCalendar();
  };
  const onEdit = async () => {
    if (calendar.id === "") {
      await getRepository().createCalendar(calendar).catch(err => {
        cancelCalendar();
        throw new Error(`Could not create calendar ${calendar.name}: ${err.message}`);
      });
      saveCalendar(calendar);
    } else if (calendar.source === originalCalendar.source) {
      const changes = {
        name: calendar.name != originalCalendar.name,
        desc: calendar.desc != originalCalendar.desc,
        color: calendar.color != originalCalendar.color
      }
      await getRepository().editCalendar(calendar, changes).catch(err => {
        cancelCalendar();
        throw new Error(`Could not edit calendar ${calendar.name}: ${err.message}`);
      });
      saveCalendar(calendar);
    } else {
      await getRepository().moveCalendar(calendar).catch(err => {
        cancelCalendar();
        throw new Error(`Could not move calendar ${calendar.name}: ${err.message}`);
      });
      saveCalendar(calendar);
    }
  };

  let canSubmit: boolean = $derived(calendar && calendar.name !== "" && calendar.source !== "");
</script>

<EditableModal
  title={title}
  deleteConfirmation={`Are you sure you want to delete calendar "${calendar ? calendar.name : ""}"?`}
  bind:editMode={editMode}
  bind:showCreateModal={showCreateModalInternal}
  bind:showModal={showModalInternal}
  onDelete={onDelete}
  onEdit={onEdit}
  editable={false}
  submittable={canSubmit}
>
  {#if calendar != EmptyCalendar}
    <TextInput bind:value={calendar.name} name="name" placeholder="Name" editable={editMode} />
    <SelectInput bind:value={calendar.source} name="source" placeholder="Source" options={selectableSources} editable={editMode} />
    {#if editMode}
      <ColorInput bind:color={calendar.color} name="color" editable={editMode} />
    {/if}
    <TextInput bind:value={calendar.desc} name="desc" placeholder="Description" multiline={true} editable={editMode} />
  {/if}
</EditableModal>
<script lang="ts">
  import EditableModal from "./EditableModal.svelte";
  import TextInput from "../forms/TextInput.svelte";

  import { EmptyCalendar, NoChangesCalendar, NoOp } from "$lib/client/placeholders";
  import { getRepository } from "$lib/client/data/repository.svelte";
  import { deepCopy } from "$lib/common/misc";
  import SelectInput from "../forms/SelectInput.svelte";
  import ColorInput from "../forms/ColorInput.svelte";
  import { queueNotification } from "../../lib/client/notifications";
  import { getSettings } from "$lib/client/data/settings.svelte";
  import { UserSettingKeys } from "../../types/settings";
  import { ColorKeys } from "../../types/colors";
  import IconButton from "../interactive/IconButton.svelte";
  import { ArchiveRestore } from "lucide-svelte";

  interface Props {
    showCreateModal?: () => any;
    showModal?: (initial?: CalendarModel) => Promise<CalendarModel>;
  }

  let {
    showCreateModal = $bindable(),
    showModal = $bindable(),
  }: Props = $props();

  const settings = getSettings();
  const repository = getRepository();

  let showModalInternal: (initial?: CalendarModel, edit?: boolean) => Promise<CalendarModel> = $state(Promise.reject);

  let calendar: CalendarModel = $state(EmptyCalendar);
  let originalCalendar: CalendarModel = $state(EmptyCalendar);
  let editMode: boolean = $state(false);

  showModal = async (initial?: CalendarModel): Promise<CalendarModel> => {
    if (!initial) {
      calendar = {
        id: "",
        source: "",
        name: "",
        desc: "",
        color: "",
        overridden: false,
        can_edit: true,
        can_delete: true,
        can_add_events: true,
      };
    } else {
      calendar = await deepCopy(initial);
      originalCalendar = await deepCopy(initial);
    }

    return showModalInternal(calendar);
  };

  let title: string = $derived(calendar.id ? (editMode ? "Edit calendar" : "Calendar") : "Add calendar");

  let selectableSources = $derived(
    repository.sources
      .filter(source => source.id === calendar.source || (editMode && source.can_add_calendars))
      .map(source => ({ value: source.id, name: source.name }))
  );

  const onDelete = async () => {
    return await getRepository().deleteCalendar(calendar.id).then(() => calendar).catch(err => {
      throw new Error(`Could not delete calendar ${calendar.name}: ${err.message}`);
    });
  };
  const onEdit = async () => {
    if (calendar.id === "") {
      return await getRepository().createCalendar(calendar).then(() => calendar).catch(err => {
        throw new Error(`Could not create calendar ${calendar.name}: ${err.message}`);
      });
    } else if (calendar.source === originalCalendar.source) {
      const changes = {
        name: calendar.name != originalCalendar.name,
        desc: calendar.desc != originalCalendar.desc,
        color: calendar.color != originalCalendar.color
      }
      return await getRepository().editCalendar(calendar, changes, !calendar.can_edit).then(() => calendar).catch(err => {
        throw new Error(`Could not edit calendar ${calendar.name}: ${err.message}`);
      });
    } else {
      return await getRepository().moveCalendar(calendar).then(() => calendar).catch(err => {
        throw new Error(`Could not move calendar ${calendar.name}: ${err.message}`);
      });
    }
  };
  const resetOverrides = async () => {
    calendar.overridden = false;
    getRepository().editCalendar(calendar, NoChangesCalendar, true).catch(err => {
      calendar.overridden = true;
      queueNotification(ColorKeys.Danger, `Could not reset calendar ${calendar.name}: ${err.message}`);
      return;
    }).then(async () => {
      getRepository().getCalendar(calendar.id, true).catch(err => {
        calendar.overridden = true;
        queueNotification(ColorKeys.Danger, `Could not reset event ${calendar.name}: ${err.message}`);
        return;
      }).then((fetched) => {
        calendar = fetched as CalendarModel;
      });
    });
  }

  let canSubmit: boolean = $derived(calendar && calendar.name !== "" && calendar.source !== "");
</script>

<EditableModal
  title={title}
  deleteConfirmation={`Are you sure you want to delete calendar "${calendar ? calendar.name : ""}"?`}
  bind:editMode={editMode}
  bind:showModal={showModalInternal}
  onDelete={onDelete}
  onEdit={onEdit}
  deletable={calendar?.can_delete}
  submittable={canSubmit}
>
  {#if calendar != EmptyCalendar}
    <TextInput bind:value={calendar.name} name="name" placeholder="Name" editable={editMode} />
    <SelectInput bind:value={calendar.source} name="source" placeholder="Source" options={selectableSources} editable={calendar.id === ""} />
    {#if editMode}
      <ColorInput bind:color={calendar.color} name="color" editable={editMode} />
    {/if}
    {#if editMode || calendar.desc}
      <TextInput bind:value={calendar.desc} name="desc" placeholder="Description" multiline={true} editable={editMode} />
    {/if}
    {#if calendar.id && settings.userSettings[UserSettingKeys.DebugMode]}
      <TextInput value={calendar.id} name="id" placeholder="Calendar ID" editable={false} />
    {/if}
  {/if}
  {#snippet extraButtonsLeft()}
    {#if calendar != EmptyCalendar && !editMode && calendar.overridden}
      <IconButton onClick={resetOverrides} alt="Reset">
        <ArchiveRestore/>
      </IconButton>
    {/if}
  {/snippet}
</EditableModal>
<script lang="ts">
  import Button from "../interactive/Button.svelte";
  import EditableModal from "./EditableModal.svelte";
  import TextInput from "../forms/TextInput.svelte";

  import { EmptyCalendar, NoChangesCalendar, NoOp } from "$lib/client/placeholders";
  import { getRepository } from "$lib/client/repository";
  import { deepCopy } from "$lib/common/misc";
  import SelectInput from "../forms/SelectInput.svelte";
  import ColorInput from "../forms/ColorInput.svelte";
  import { queueNotification } from "../../lib/client/notifications";
  import { getSettings } from "$lib/client/settings.svelte";
  import { UserSettingKeys } from "../../types/settings";

  interface Props {
    showCreateModal?: () => any;
    showModal?: (calendar: CalendarModel) => Promise<CalendarModel>;
  }

  let {
    showCreateModal = $bindable(),
    showModal = $bindable(),
  }: Props = $props();

  const settings = getSettings();

  let calendar: CalendarModel = $state(EmptyCalendar);
  let originalCalendar: CalendarModel = $state(EmptyCalendar);

  let promiseResolve: (value: CalendarModel | PromiseLike<CalendarModel>) => void = $state(NoOp);
  let promiseReject: (reason?: any) => void = $state(NoOp);

  showCreateModal = async () => {
    promiseReject();

    calendar = {
      id: "",
      source: "",
      name: "",
      desc: "",
      color: "",
      overridden: false,
    };

    showCreateModalInternal();
  }
  showModal = async (original: CalendarModel): Promise<CalendarModel> => {
    promiseReject();

    editMode = false;
    calendar = await deepCopy(original);
    originalCalendar = await deepCopy(original);

    setTimeout(showModalInternal(), 0);

    return new Promise((resolve, reject) => {
      promiseResolve = resolve;
      promiseReject = reject;
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
    promiseReject();
  };
  const onEdit = async () => {
    if (calendar.id === "") {
      await getRepository().createCalendar(calendar).catch(err => {
        promiseReject();
        throw new Error(`Could not create calendar ${calendar.name}: ${err.message}`);
      });
      promiseResolve(calendar);
    } else if (calendar.source === originalCalendar.source) {
      const changes = {
        name: calendar.name != originalCalendar.name,
        desc: calendar.desc != originalCalendar.desc,
        color: calendar.color != originalCalendar.color
      }
      await getRepository().editCalendar(calendar, changes, true).catch(err => {
        promiseReject();
        throw new Error(`Could not edit calendar ${calendar.name}: ${err.message}`);
      });
      promiseResolve(calendar);
    } else {
      await getRepository().moveCalendar(calendar).catch(err => {
        promiseReject();
        throw new Error(`Could not move calendar ${calendar.name}: ${err.message}`);
      });
      promiseResolve(calendar);
    }
  };
  const resetOverrides = async () => {
    calendar.overridden = false;
    getRepository().editCalendar(calendar, NoChangesCalendar, true).catch(err => {
      calendar.overridden = true;
      queueNotification("failure", `Could not reset calendar ${calendar.name}: ${err.message}`);
      return;
    }).then(async () => {
      getRepository().getCalendar(calendar.id, true).catch(err => {
        calendar.overridden = true;
        queueNotification("failure", `Could not reset event ${calendar.name}: ${err.message}`);
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
  bind:showCreateModal={showCreateModalInternal}
  bind:showModal={showModalInternal}
  onDelete={onDelete}
  onEdit={onEdit}
  onCancel={promiseReject}
  deletable={false}
  submittable={canSubmit}
>
  {#if calendar != EmptyCalendar}
    <TextInput bind:value={calendar.name} name="name" placeholder="Name" editable={editMode} />
    <SelectInput bind:value={calendar.source} name="source" placeholder="Source" options={selectableSources} editable={false} />
    {#if editMode}
      <ColorInput bind:color={calendar.color} name="color" editable={editMode} />
    {/if}
    {#if editMode || calendar.desc}
      <TextInput bind:value={calendar.desc} name="desc" placeholder="Description" multiline={true} editable={editMode} />
    {/if}
    {#if settings.userSettings[UserSettingKeys.DebugMode]}
      <TextInput bind:value={calendar.id} name="id" placeholder="ID" editable={false} />
    {/if}
  {/if}
  {#snippet extraButtonsLeft()}
    {#if calendar != EmptyCalendar && !editMode && calendar.overridden}
      <Button color="accent" onClick={resetOverrides}>Reset</Button>
    {/if}
  {/snippet}
</EditableModal>
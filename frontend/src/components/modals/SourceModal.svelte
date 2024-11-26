<script lang="ts">
  import EditableModal from "./EditableModal.svelte";
  import SelectButtons from "../forms/SelectButtons.svelte";
  import TextInput from "../forms/TextInput.svelte";

  import { EmptySource, NoOp } from "$lib/client/placeholders";
  import { createSource, deleteSource, editSource } from "$lib/client/repository";
  import { isValidUrl, valid } from "$lib/client/validation";
  import { queueNotification } from "$lib/client/notifications";

  interface Props {
    showCreateModal?: () => any;
    showModal?: (source: SourceModel) => Promise<SourceModel>;
  }

  let {
    showCreateModal = $bindable(),
    showModal = $bindable(),
  }: Props = $props();

  let sourceDetailed: SourceModel = $state(EmptySource);

  let saveSource = (_: SourceModel | PromiseLike<SourceModel>) => {};
  let cancelSource = (_?: any) => {};

  showCreateModal = () => {
    cancelSource();

    sourceDetailed = {
      id: "",
      name: "",
      type: "caldav",
      settings: {},
      auth_type: "none",
      auth: {},
      collapsed: false
    };

    showCreateModalInternal();
  }
  showModal = async (source: SourceModel): Promise<SourceModel> => {
    cancelSource();

    const res = await fetch(`/api/sources/${source.id}`);
    if (res.ok) {
      sourceDetailed = await res.json();
    } else {
      queueNotification("failure", `Failed to fetch source details: ${res.statusText}`);
      return Promise.reject();
    }

    showModalInternal();
    return new Promise((resolve, reject) => {
      saveSource = resolve;
      cancelSource = reject;
    })
  };

  let showCreateModalInternal: () => any = $state(NoOp);
  let showModalInternal: () => any = $state(NoOp);

  let editMode: boolean = $state(false);
  let title: string = $derived(sourceDetailed.id ? (editMode ? "Edit source" : "Source") : "Add source");

  const onDelete = async () => {
    await deleteSource(sourceDetailed.id).catch(err => {
      throw new Error(`Could not delete source: ${err.message}`);
    });
    cancelSource();
  };
  const onEdit = async () => {
    if (sourceDetailed.id === "") {
      await createSource(sourceDetailed).catch(err => {
        cancelSource();
        throw new Error(`Could not create source: ${err.message}`);
      });
      saveSource(sourceDetailed);
    } else {
      await editSource(sourceDetailed).catch(err => {
        cancelSource();
        throw new Error(`Could not edit source: ${err.message}`);
      });
      saveSource(sourceDetailed);
    }
  };

  let caldavLinkValidity: Validity = $state(valid);
  let icalLinkValidity: Validity = $state(valid);

  let canSubmit: boolean = $derived(sourceDetailed && sourceDetailed.name !== "" && sourceDetailed.type !== "" && (
    (sourceDetailed.type === "caldav" && caldavLinkValidity?.valid) ||
    (sourceDetailed.type === "ical" && icalLinkValidity?.valid)
  ));
</script>

<EditableModal
  title={title}
  deleteConfirmation={`Are you sure you want to delete source "${sourceDetailed ? sourceDetailed.name : ""}"?`}
  bind:editMode={editMode}
  bind:showCreateModal={showCreateModalInternal}
  bind:showModal={showModalInternal}
  onDelete={onDelete}
  onEdit={onEdit}
  submittable={canSubmit}
>
  {#if sourceDetailed}
    <TextInput bind:value={sourceDetailed.name} name="name" placeholder="Name" editable={editMode} />
    <!--
    <SelectInput bind:value={sourceType} name="type" placeholder={"Type"} editable={editMode} options={[
      {
        value: "caldav",
        name: "CalDav"
      },
      {
        value: "ical",
        name: "iCal"
      }
    ]}></SelectInput>
    -->
    <SelectButtons bind:value={sourceDetailed.type} name="type" placeholder={"Type"} editable={editMode} options={[
      {
        value: "caldav",
        name: "CalDav"
      },
      {
        value: "ical",
        name: "iCal"
      }
    ]}/>
    {#if sourceDetailed.type === "caldav"}
      <TextInput bind:value={sourceDetailed.settings.url} name="caldav_url" placeholder="CalDav URL" editable={editMode} validation={isValidUrl} bind:validity={caldavLinkValidity} />
    {/if}
    {#if sourceDetailed.type === "ical"}
      <TextInput bind:value={sourceDetailed.settings.url} name="ical_url" placeholder="iCal URL" editable={editMode} validation={isValidUrl} bind:validity={icalLinkValidity} />
    {/if}
    
    <SelectButtons bind:value={sourceDetailed.auth_type} name="auth_type" placeholder={"Authentication Type"} editable={editMode} options={[
      {
        value: "none",
        name: "None",
      },
      {
        value: "basic",
        name: "Password",
      },
      {
        value: "bearer",
        name: "Token",
      },
    ]}/>
    {#if sourceDetailed.auth_type === "basic"}
      <TextInput bind:value={sourceDetailed.auth.username} name="auth_username" placeholder="Username" editable={editMode} />
      <TextInput bind:value={sourceDetailed.auth.password} name="auth_password" placeholder="Password" editable={editMode} password={true} />
    {/if}
    {#if sourceDetailed.auth_type === "bearer"}
      <TextInput bind:value={sourceDetailed.auth.token} name="auth_token" placeholder="Token" editable={editMode} password={true} />
    {/if}
  {/if}
</EditableModal>
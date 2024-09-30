<script lang="ts">
  import EditableModal from "./EditableModal.svelte";
  import TextInput from "../forms/TextInput.svelte";
  import SelectButtons from "../forms/SelectButtons.svelte";
  import { queueNotification } from "$lib/client/notifications";
  import { createSource, deleteSource, editSource } from "$lib/client/repository";

  export let source: SourceModel;
  let sourceDetailed: SourceModel;

  export const showCreateModal = () => {
    sourceDetailed = source;
    showCreateModalInternal();
  }
  export const showModal = async () => {
    const res = await fetch(`/api/sources/${source.id}`);
    if (res.ok) {
      sourceDetailed = await res.json();
    } else {
      queueNotification("failure", `Failed to fetch source details: ${res.statusText}`);
      return
    }

    showModalInternal();
  };
  let showCreateModalInternal: () => boolean;
  let showModalInternal: () => boolean;

  let title: string;
  $: title = (sourceDetailed && sourceDetailed.id) ? (editMode ? "Edit source" : "Source") : "Add source";

  let editMode: boolean;

  const onDelete = async () => {
    const res = await deleteSource(sourceDetailed.id);
    if (res === "") return "";
    else return `Could not delete source: ${res}`;
  };
  const onEdit = async () => {
    let res: string;
    if (sourceDetailed.id === "") {
      res = await createSource(sourceDetailed);
    } else {
      res = await editSource(sourceDetailed);
    }
    if (res === "") return "";
    else return `Could not edit source: ${res}`;
  };
</script>

<EditableModal
  title={title}
  deleteConfirmation={`Are you sure you want to delete source "${sourceDetailed ? sourceDetailed.name : ""}"?`}
  bind:editMode={editMode}
  bind:showCreateModal={showCreateModalInternal}
  bind:showModal={showModalInternal}
  onDelete={onDelete}
  onEdit={onEdit}
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
      <TextInput bind:value={sourceDetailed.settings.url} name="caldav_url" placeholder="CalDav URL" editable={editMode} />
    {/if}
    {#if sourceDetailed.type === "ical"}
      <TextInput bind:value={sourceDetailed.settings.url} name="ical_url" placeholder="iCal URL" editable={editMode} />
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
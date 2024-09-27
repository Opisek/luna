<script lang="ts">
  import EditableModal from "./EditableModal.svelte";
  import TextInput from "../forms/TextInput.svelte";
  import SelectInput from "../forms/SelectInput.svelte";
  import SelectButtons from "../forms/SelectButtons.svelte";

  export let source: SourceModel;

  let sourceType: string = "caldav";

  let caldavUrl: string;
  let icalUrl: string;

  let authType: string = "none";
  let authUsername: string;
  let authPassword: string;
  let authToken: string;

  export let showModal: () => boolean;

  let title: string;
  $: title = (source && source.id) ? (editMode ? "Edit source" : "Source") : "Create source";

  let editMode: boolean;

  const onDelete = async () => {
    return Promise.resolve("");
  };
  const onEdit = () => {
    return Promise.resolve("");
  };
</script>

<EditableModal
  title={title}
  deleteConfirmation={`Are you sure you want to delete source "${source ? source.name : ""}"?`}
  isNew={!(source && source.id)}
  bind:editMode={editMode}
  bind:showModal={showModal}
  onDelete={onDelete}
  onEdit={onEdit}
>
  {#if source}
    <TextInput bind:value={source.name} name="name" placeholder="Name" editable={editMode} />
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
    <SelectButtons bind:value={sourceType} name="type" placeholder={"Type"} editable={editMode} options={[
      {
        value: "caldav",
        name: "CalDav"
      },
      {
        value: "ical",
        name: "iCal"
      }
    ]}/>
    {#if sourceType === "caldav"}
      <TextInput bind:value={caldavUrl} name="caldav_url" placeholder="CalDav URL" editable={editMode} />
    {/if}
    {#if sourceType === "ical"}
      <TextInput bind:value={icalUrl} name="ical_url" placeholder="iCal URL" editable={editMode} />
    {/if}
    
    <SelectButtons bind:value={authType} name="auth_type" placeholder={"Authentication Type"} editable={editMode} options={[
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
    {#if authType === "basic"}
      <TextInput bind:value={authUsername} name="auth_username" placeholder="Username" editable={editMode} />
      <TextInput bind:value={authPassword} name="auth_password" placeholder="Password" editable={editMode} password={true} />
    {/if}
    {#if authType === "bearer"}
      <TextInput bind:value={authToken} name="auth_token" placeholder="Token" editable={editMode} password={true} />
    {/if}
  {/if}
</EditableModal>
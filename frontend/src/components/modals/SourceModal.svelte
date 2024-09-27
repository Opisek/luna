<script lang="ts">
  import EditableModal from "./EditableModal.svelte";
  import TextInput from "../forms/TextInput.svelte";
  import SelectInput from "../forms/SelectInput.svelte";
  import SelectButtons from "../forms/SelectButtons.svelte";

  export let source: SourceModel;
  let sourceType: string;
  let authType: string = "none";

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
  {/if}
</EditableModal>
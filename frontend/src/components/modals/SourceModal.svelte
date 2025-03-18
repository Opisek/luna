<script lang="ts">
  import EditableModal from "./EditableModal.svelte";
  import SelectButtons from "../forms/SelectButtons.svelte";
  import TextInput from "../forms/TextInput.svelte";

  import { EmptySource, NoOp } from "$lib/client/placeholders";
  import { getRepository } from "$lib/client/repository";
  import { deepCopy, deepEquality } from "$lib/common/misc";
  import { isValidFile, isValidIcalFile, isValidPath, isValidUrl, valid } from "$lib/client/validation";
  import { queueNotification } from "$lib/client/notifications";
  import FileUpload from "../forms/FileUpload.svelte";
  import { fetchResponse } from "../../lib/client/net";

  interface Props {
    showCreateModal?: () => any;
    showModal?: (source: SourceModel) => Promise<SourceModel>;
  }

  let {
    showCreateModal = $bindable(),
    showModal = $bindable(),
  }: Props = $props();

  let sourceDetailed: SourceModel = $state(EmptySource);
  let originalSource: SourceModel;

  let saveSource = (_: SourceModel | PromiseLike<SourceModel>) => {};
  let cancelSource = (_?: any) => {};

  showCreateModal = () => {
    cancelSource();

    sourceDetailed = {
      id: "",
      name: "",
      type: "caldav",
      settings: {
        location: "remote",
        file: null,
        fileId: "",
      },
      auth_type: "none",
      auth: {},
    };

    showCreateModalInternal();
  }
  showModal = async (source: SourceModel): Promise<SourceModel> => {
    cancelSource();

    // TODO: this should be a call to repository with force refresh = true
    sourceDetailed = await getRepository().getSourceDetails(source.id).catch(err => {
      queueNotification("failure", `Could not get source details: ${err.message}`);
      return Promise.reject();
    });

    // so that when we edit a caldav source into an ical source, the location selection will default to some value (remote):
    if (sourceDetailed.type !== "ical") sourceDetailed.settings.location = "remote";

    if (sourceDetailed.type === "ical" && sourceDetailed.settings.location === "database" && sourceDetailed.settings.file !== null) {
      const fileId = sourceDetailed.settings.file;
      sourceDetailed.settings.fileId = fileId;

      const res = await fetchResponse(`/api/files/${fileId}`, { method: "HEAD" }).catch(err => {
        queueNotification("failure", `Could not get file: ${err.message}`);
        sourceDetailed.settings.file = null;
      });

      if (res) {
        let filename = `${fileId}.ics`;

        const header = res.headers.get("Content-Disposition")
        if (header) {
          const remoteFilename = header
            .split(";")
            .map(x => x.trim())
            .filter(x => x.startsWith("filename="))
            .map(x => x.split("=")[1]);
          
          if (remoteFilename.length > 0) filename = remoteFilename[0];
        }

        // https://stackoverflow.com/questions/52078853/is-it-possible-to-update-filelist
        const list = new DataTransfer();
        const file = new File([], filename);
        list.items.add(file);
        sourceDetailed.settings.file = list.files;
      }
    } else {
      sourceDetailed.settings.file = null;
      sourceDetailed.settings.fileId = "";
    }

    originalSource = await deepCopy(sourceDetailed);
    if (sourceDetailed.settings.file !== null) originalSource.settings.file = sourceDetailed.settings.file;

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
    await getRepository().deleteSource(sourceDetailed.id).catch(err => {
      throw new Error(`Could not delete source ${sourceDetailed.name}: ${err.message}`);
    });
    cancelSource();
  };
  const onEdit = async () => {
    if (sourceDetailed.id === "") {
      await getRepository().createSource(sourceDetailed).catch(err => {
        cancelSource();
        throw new Error(`Could not create source ${sourceDetailed.name}: ${err.message}`);
      });
      saveSource(sourceDetailed);
    } else {
      if (originalSource.settings.file instanceof String && sourceDetailed.settings.file instanceof FileList && sourceDetailed.settings.file.length === 1 && sourceDetailed.settings.file[0].name === originalSource.settings.file) {
        sourceDetailed.settings.file = sourceDetailed.settings.file[0];
      }
      const changes = {
        name: sourceDetailed.name != originalSource.name,
        type: sourceDetailed.type != originalSource.type || !deepEquality(sourceDetailed.settings, originalSource.settings),
        settings: !deepEquality(sourceDetailed.settings, originalSource.settings),
        auth: sourceDetailed.auth_type != originalSource.auth_type || !deepEquality(sourceDetailed.auth, originalSource.auth)
      }
      await getRepository().editSource(sourceDetailed, changes).catch(err => {
        cancelSource();
        throw new Error(`Could not edit source ${sourceDetailed.name}: ${err.message}`);
      });
      saveSource(sourceDetailed);
    }
  };

  let caldavLinkValidity: Validity = $state(valid);
  let icalLinkValidity: Validity = $state(valid);
  let icalFileValidity: Validity = $state(valid);
  let icalPathValidity: Validity = $state(valid);

  let canSubmit: boolean = $derived(sourceDetailed && sourceDetailed.name !== "" && sourceDetailed.type !== "" && (
    (sourceDetailed.type === "caldav" && caldavLinkValidity?.valid) ||
    (sourceDetailed.type === "ical" && (
      (sourceDetailed.settings.location === "remote"   && icalLinkValidity?.valid) ||
      (sourceDetailed.settings.location === "database" && icalFileValidity?.valid) ||
      (sourceDetailed.settings.location === "local"    && icalPathValidity?.valid)
    ))
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


    {#if sourceDetailed.type === "ical"}
      <SelectButtons bind:value={sourceDetailed.settings.location} name="ical_location" placeholder={"File Location"} editable={editMode} options={[
        {
          value: "remote",
          name: "Internet Link",
        },
        {
          value: "database",
          name: "Upload File",
        },
        {
          value: "local",
          name: "Server Filepath",
        },
      ]}/>
    {/if}

    {#if sourceDetailed.type === "caldav"}
      <TextInput bind:value={sourceDetailed.settings.url} name="caldav_url" placeholder="CalDav URL" editable={editMode} validation={isValidUrl} bind:validity={caldavLinkValidity} />
    {/if}
    {#if sourceDetailed.type === "ical"}
      {#if sourceDetailed.settings.location === "remote"}
        <TextInput bind:value={sourceDetailed.settings.url} name="ical_url" placeholder="iCal URL" editable={editMode} validation={isValidUrl} bind:validity={icalLinkValidity} />
      {:else if sourceDetailed.settings.location === "database"}
          <FileUpload bind:files={sourceDetailed.settings.file} bind:fileId={sourceDetailed.settings.fileId} name="ical_file" placeholder="iCal File" editable={editMode} validation={isValidIcalFile} bind:validity={icalFileValidity} />
      {:else if sourceDetailed.settings.location === "local"}
        <TextInput bind:value={sourceDetailed.settings.path} name="ical_path" placeholder="iCal Path" editable={editMode} validation={isValidPath} bind:validity={icalPathValidity} />
      {/if}
    {/if}
    
    {#if !(sourceDetailed.type === "ical" && sourceDetailed.settings.location !== "remote")}
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

    <!-- TODO: show id debug prefence -->
    <!--
    <TextInput bind:value={source.id} name="id" placeholder="ID" editable={false} />
    -->
  {/if}
</EditableModal>
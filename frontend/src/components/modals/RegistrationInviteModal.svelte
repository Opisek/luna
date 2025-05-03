<script lang="ts">
  import EditableModal from "./EditableModal.svelte";
  import TextInput from "../forms/TextInput.svelte";

  import { EmptyRegistrationInvite, NoOp } from "$lib/client/placeholders";
  import { deepCopy } from "$lib/common/misc";
  import { getSettings } from "$lib/client/data/settings.svelte";
  import { UserSettingKeys } from "../../types/settings";
  import DateTimeInput from "../forms/DateTimeInput.svelte";
  import SelectInput from "../forms/SelectInput.svelte";
  import { page } from "$app/state";
  import { getRegistrationInvites } from "../../lib/client/data/invites.svelte";
  import Image from "../layout/Image.svelte";
  import Horizontal from "../layout/Horizontal.svelte";
  
  interface Props {
    showCreateModal?: () => Promise<RegistrationInvite>;
    showModal?: (session: RegistrationInvite, editable: boolean) => Promise<RegistrationInvite>;
  }

  let {
    showCreateModal = $bindable(),
    showModal = $bindable(),
  }: Props = $props();

  const settings = getSettings();
  const invites = getRegistrationInvites();

  let invite: RegistrationInvite = $state(EmptyRegistrationInvite);
  let originalInvite: RegistrationInvite = $state(EmptyRegistrationInvite);

  let promiseResolve: (value: RegistrationInvite | PromiseLike<RegistrationInvite>) => void = $state(NoOp);
  let promiseReject: (reason?: any) => void = $state(NoOp);

  showModal = async (original: RegistrationInvite, edit: boolean = false): Promise<RegistrationInvite> => {
    promiseReject();

    editMode = edit;
    invite = await deepCopy(original);
    originalInvite = await deepCopy(original);

    if (editMode) setTimeout(showCreateModalInternal(), 0);
    else setTimeout(showModalInternal(), 0);

    // TODO: what if we only show but the modal and the close? memory leak?
    return new Promise((resolve, reject) => {
      promiseResolve = ((res) => {
        invite = EmptyRegistrationInvite;
        resolve(res);
      });
      promiseReject = ((err) => {
        invite = EmptyRegistrationInvite;
        reject(err);
      });
    })
  };
  showCreateModal = () => {
    return showModal(EmptyRegistrationInvite, true);
  }

  let showCreateModalInternal: () => any = $state(NoOp);
  let showModalInternal: () => any = $state(NoOp);

  let editMode: boolean = $state(false);
  let title: string = $derived(invite.invite_id ? "Registration Invite" : "Invite User");

  const onDelete = async () => {
    await invites.revokeInvite(invite.invite_id).catch(err => {
      throw new Error(`Could not delete registration invite: ${err.message}`);
    });

    promiseResolve(originalInvite);
  };
  const onEdit = async () => {
    if (invite.invite_id === "") {
      const newInvite = await invites.createInvite(duration).catch(err => {
        throw new Error(`Could not create registration invite: ${err.message}`);
      });

      editMode = false;

      promiseResolve(newInvite);

      setTimeout(() => {
        showModal(newInvite, false);
      }, 50);
    } else  {
      throw new Error("Not implemented");
    }
  };

  let duration: number = $state(3600);

  let inviteCode: string = $derived(invite.code);
  let inviteLink: string = $derived(`${page.url.origin}/register?code=${inviteCode}`);
</script>

<EditableModal
  title={title}
  deleteConfirmation={`Are you sure you want to delete the registration invite?`}
  bind:editMode={editMode}
  bind:showCreateModal={showCreateModalInternal}
  bind:showModal={showModalInternal}
  onDelete={onDelete}
  onEdit={onEdit}
  onCancel={promiseReject}
  editable={editMode}
  deletable={true}
  submittable={true}
>
  {#if invite.invite_id == ""}
    <SelectInput
      bind:value={duration}
      name="duration"
      placeholder="Invite Expiration"
      options={[
        { name: "1 Hour", value: 60 * 60 },
        { name: "3 Hours", value: 3 * 60 * 60 },
        { name: "1 Day", value: 24 * 60 * 60 },
        { name: "3 Days", value: 3 * 24 * 60 * 60 },
        { name: "7 Days", value: 7 * 24 * 60 * 60 },
      ]}
    />
  {:else}

    <Horizontal position="center">
      <Image
          src={`/api/invites/${invite.invite_id}/qr`}
          alt="QR Code"
          large={true}
      />
    </Horizontal>

    <TextInput
      value={inviteCode}
      name="invite_id"
      placeholder="Invite Code"
      editable={false}
      displayCopyButton={true}
      mono={true}
    />

    <TextInput
      value={inviteLink}
      name="user_id"
      placeholder="Invite Link"
      editable={false}
      displayCopyButton={true}
    />
      
    <DateTimeInput value={invite.created_at} allDay={false} placeholder="Creation Date" name="created_at" editable={false}/>
    <DateTimeInput value={invite.expires_at} allDay={false} placeholder="Expiry Date" name="expires_at" editable={false}/>

    {#if settings.userSettings[UserSettingKeys.DebugMode]}
      <TextInput value={invite.invite_id} name="id" placeholder="Invite ID" editable={false} />
    {/if}
  {/if}
</EditableModal>
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
    showModal: (initial?: RegistrationInvite, edit?: boolean) => Promise<RegistrationInvite>;
  }

  let {
    showModal = $bindable(),
  }: Props = $props();

  let showModalInternal: (initial?: RegistrationInvite, edit?: boolean) => Promise<RegistrationInvite> = $state(Promise.reject);

  const settings = getSettings();
  const invites = getRegistrationInvites();

  let invite: RegistrationInvite = $state(EmptyRegistrationInvite);
  let originalInvite: RegistrationInvite = $state(EmptyRegistrationInvite);

  showModal = async (initial?: RegistrationInvite, edit?: boolean): Promise<RegistrationInvite> => {
    if (!initial) initial = EmptyRegistrationInvite;
    invite = await deepCopy(initial);
    originalInvite = await deepCopy(initial);
    return showModalInternal(initial, edit);
  };

  let editMode: boolean = $state(false);
  let title: string = $derived(invite.id? "Registration Invite" : "Invite User");

  const onDelete = async () => {
    return invites.revokeInvite(invite.id).then(() => invite).catch(err => {
      throw new Error(`Could not delete registration invite: ${err.message}`);
    });
  };
  const onEdit = async () => {
    if (invite.id === "") {
      return invites.createInvite(duration).then((newInvite) => {
        setTimeout(() => {
          showModal(newInvite, false);
        }, 50);
        return newInvite;
      }).catch(err => {
        throw new Error(`Could not create registration invite: ${err.message}`);
      });
    } else {
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
  bind:showModal={showModalInternal}
  onDelete={onDelete}
  onEdit={onEdit}
  editable={editMode}
  deletable={true}
  submittable={true}
>
  {#if invite.id == ""}
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
          src={`/api/invites/${invite.id}/qr`}
          alt="QR Code"
          large={true}
      />
    </Horizontal>

    <TextInput
      value={inviteCode}
      name="id"
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
      <TextInput value={invite.id} name="id" placeholder="Invite ID" editable={false} />
    {/if}
  {/if}
</EditableModal>
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
  import { t } from "@sveltia/i18n";
  
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
  let title: string = $derived(invite.id ? t("invite.title.view") : t("invite.title.create"));

  const onDelete = async () => {
    return invites.revokeInvite(invite.id).then(() => invite).catch(err => {
      throw new Error(t("invite.error.delete", { values: { msg: err.message } }));
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
        throw new Error(t("invite.error.create", { values: { msg: err.message } }));
      });
    } else {
      throw new Error(t("error.unimplemented"));
    }
  };

  let duration: number = $state(3600);

  let inviteCode: string = $derived(invite.code);
  let inviteLink: string = $derived(`${page.url.origin}/register?code=${inviteCode}`);
</script>

<EditableModal
  title={title}
  deleteConfirmation={t("invite.confirm.delete")}
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
      placeholder={t("invite.expiry.display")}
      options={[
        { name: t("invite.expiry.preset.1h"), value: 60 * 60 },
        { name: t("invite.expiry.preset.3h"), value: 3 * 60 * 60 },
        { name: t("invite.expiry.preset.1d"), value: 24 * 60 * 60 },
        { name: t("invite.expiry.preset.3d"), value: 3 * 24 * 60 * 60 },
        { name: t("invite.expiry.preset.7d"), value: 7 * 24 * 60 * 60 },
      ]}
    />
  {:else}

    <Horizontal position="center">
      <Image
          src={`/api/invites/${invite.id}/qr`}
          alt={t("invite.qr")}
          large={true}
      />
    </Horizontal>

    <TextInput
      value={inviteCode}
      name="invite_code"
      placeholder={t("invite.code")}
      editable={false}
      displayCopyButton={true}
      mono={true}
    />

    <TextInput
      value={inviteLink}
      name="invite_link"
      placeholder={t("invite.link")}
      editable={false}
      displayCopyButton={true}
    />
      
    <DateTimeInput value={invite.created_at} allDay={false} placeholder={t("invite.date.creation")} name="created_at" editable={false}/>
    <DateTimeInput value={invite.expires_at} allDay={false} placeholder={t("invite.date.expiry")} name="expires_at" editable={false}/>

    {#if settings.userSettings[UserSettingKeys.DebugMode]}
      <TextInput value={invite.id} name="id" placeholder={t("invite.id.label")} editable={false} />
    {/if}
  {/if}
</EditableModal>
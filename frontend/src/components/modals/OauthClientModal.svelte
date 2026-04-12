<script lang="ts">
  import EditableModal from "./EditableModal.svelte";

  import { EmptyOauthClient, NoOp } from "$lib/client/placeholders";
  import { deepCopy } from "$lib/common/misc";
  import TextInput from "../forms/TextInput.svelte";
  import { getOauthClients } from "../../lib/client/data/oauth.svelte";
  import { queueNotification } from "../../lib/client/notifications";
  import { ColorKeys } from "../../types/colors";
  import { getSettings } from "$lib/client/data/settings.svelte";
  import { UserSettingKeys } from "../../types/settings";
  import { t } from "@sveltia/i18n";
  
  interface Props {
    showModal: (initial?: OauthClientModel, edit?: boolean) => Promise<OauthClientModel>;
  }

  let {
    showModal = $bindable(),
  }: Props = $props();

  let showModalInternal: (initial?: OauthClientModel, edit?: boolean) => Promise<OauthClientModel> = $state(Promise.reject);

  const settings = getSettings();
  const clients = getOauthClients();

  let client: OauthClientModel = $state(EmptyOauthClient);
  let originalClient: OauthClientModel = $state(EmptyOauthClient);

  showModal = async (initial?: OauthClientModel, edit?: boolean): Promise<OauthClientModel> => {
    if (!initial) initial = EmptyOauthClient;

    client = await deepCopy(initial);
    originalClient = await deepCopy(initial);

    if (client.id !== "") {
      const details = await clients.getClientDetails(client.id).catch((err) => {
        queueNotification(ColorKeys.Danger, t("auth.oauth.client.error.details", { values: { name: client.name, msg: err.message } }));
        return Promise.reject();
      });
      client.client_secret = details.client_secret;
    }

    return showModalInternal();
  };

  let editMode: boolean = $state(false);
  let title: string = $derived(client.id ? t("auth.oauth.client.title.edit") : t("auth.oauth.client.title.create"));

  const onDelete = async () => {
    return clients.deleteClient(client.id).then(() => client).catch(err => {
      throw new Error(t("auth.oauth.client.error.delete", { values: { name: client.name, msg: err.message } }));
    });
  };
  const onEdit = async () => {
    if (client.id === "") {
      return clients.registerClient(client).then(newClient => newClient).catch(err => {
        throw new Error(t("auth.oauth.client.error.create", { values: { name: client.name, msg: err.message } }));
      });
    } else {
      return clients.updateClient(client).then(updatedClient => updatedClient).catch(err => {
        throw new Error(t("auth.oauth.client.error.delete", { values: { name: client.name, msg: err.message } }));
      });
    }
  };
</script>

<EditableModal
  title={title}
  deleteConfirmation={t("auth.oauth.client.confirm.delete")}
  bind:editMode={editMode}
  bind:showModal={showModalInternal}
  onDelete={onDelete}
  onEdit={onEdit}
  editable={editMode}
  deletable={true}
  submittable={true}
>
  <TextInput bind:value={client.name} name="name" placeholder={t("form.name")} editable={editMode}/>
  <TextInput bind:value={client.client_id} name="client_id" placeholder={t("auth.oauth.client.extid")} editable={editMode}/>
  <TextInput bind:value={client.client_secret} name="client_secret" placeholder={t("auth.oauth.client.secret")} editable={editMode} password={true}/>
  <TextInput bind:value={client.base_url} name="base_url" placeholder={t("auth.oauth.client.base")} editable={editMode}/>
  <TextInput bind:value={client.scope} name="scope" placeholder={t("auth.oauth.client.scope")} editable={editMode}/>

  {#if settings.userSettings[UserSettingKeys.DebugMode]}
    <TextInput value={client.id} name="id" placeholder={t("invite.id.label")} editable={false} />
  {/if}
</EditableModal>
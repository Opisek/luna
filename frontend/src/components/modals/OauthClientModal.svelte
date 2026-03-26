<script lang="ts">
  import EditableModal from "./EditableModal.svelte";

  import { EmptyOauthClient, NoOp } from "$lib/client/placeholders";
  import { deepCopy } from "$lib/common/misc";
  import TextInput from "../forms/TextInput.svelte";
  import { getOauthClients } from "../../lib/client/data/oauth.svelte";
  import { queueNotification } from "../../lib/client/notifications";
  import { ColorKeys } from "../../types/colors";
  
  interface Props {
    showModal: (initial?: OauthClientModel, edit?: boolean) => Promise<OauthClientModel>;
  }

  let {
    showModal = $bindable(),
  }: Props = $props();

  let showModalInternal: (initial?: OauthClientModel, edit?: boolean) => Promise<OauthClientModel> = $state(Promise.reject);

  const clients = getOauthClients();

  let client: OauthClientModel = $state(EmptyOauthClient);
  let originalClient: OauthClientModel = $state(EmptyOauthClient);

  showModal = async (initial?: OauthClientModel, edit?: boolean): Promise<OauthClientModel> => {
    if (!initial) initial = EmptyOauthClient;

    client = await deepCopy(initial);
    originalClient = await deepCopy(initial);

    if (client.id !== "") {
      const details = await clients.getClientDetails(client.id).catch((err) => {
        queueNotification(ColorKeys.Danger, `Could not get OAuth 2.0 client details: ${err.message}`);
        return Promise.reject();
      });
      client.client_secret = details.client_secret;
    }

    return showModalInternal();
  };

  let showCreateModalInternal: () => any = $state(NoOp);

  let editMode: boolean = $state(false);
  let title: string = $derived(client.id ? "Edit OAuth 2.0 Client" : "Register OAuth 2.0 Client");

  const onDelete = async () => {
    return clients.deleteClient(client.id).then(() => client).catch(err => {
      throw new Error(`Could not delete OAuth 2.0 client: ${err.message}`);
    });
  };
  const onEdit = async () => {
    if (client.id === "") {
      return clients.registerClient(client).then(newClient => newClient).catch(err => {
        throw new Error(`Could not register OAuth 2.0 client: ${err.message}`);
      });
    } else {
      return clients.updateClient(client).then(updatedClient => updatedClient).catch(err => {
        throw new Error(`Could not update OAuth 2.0 client: ${err.message}`);
      });
    }
  };
</script>

<EditableModal
  title={title}
  deleteConfirmation={`Are you sure you want to delete the OAuth 2.0 client?`}
  bind:editMode={editMode}
  bind:showModal={showModalInternal}
  onDelete={onDelete}
  onEdit={onEdit}
  editable={editMode}
  deletable={true}
  submittable={true}
>
  <TextInput bind:value={client.name} name="name" placeholder="Name" editable={editMode}/>
  <TextInput bind:value={client.client_id} name="client_id" placeholder="Client ID" editable={editMode}/>
  <TextInput bind:value={client.client_secret} name="client_secret" placeholder="Client Secret" editable={editMode} password={true}/>
  <TextInput bind:value={client.base_url} name="base_url" placeholder="Base URL" editable={editMode}/>
  <TextInput bind:value={client.scope} name="scope" placeholder="Scope" editable={editMode}/>
</EditableModal>
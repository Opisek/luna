<script lang="ts">
  import EditableModal from "./EditableModal.svelte";

  import { EmptyOauthClient, NoOp } from "$lib/client/placeholders";
  import { deepCopy } from "$lib/common/misc";
  import TextInput from "../forms/TextInput.svelte";
  import { getOauthClients } from "../../lib/client/data/oauth.svelte";
  import { queueNotification } from "../../lib/client/notifications";
  import { ColorKeys } from "../../types/colors";
  
  interface Props {
    showCreateModal?: () => Promise<OauthClientModel>;
    showModal?: (session: OauthClientModel, editable: boolean) => Promise<OauthClientModel>;
  }

  let {
    showCreateModal = $bindable(),
    showModal = $bindable(),
  }: Props = $props();

  const clients = getOauthClients();

  let client: OauthClientModel = $state(EmptyOauthClient);
  let originalClient: OauthClientModel = $state(EmptyOauthClient);

  let promiseResolve: (value: OauthClientModel | PromiseLike<OauthClientModel>) => void = $state(NoOp);
  let promiseReject: (reason?: any) => void = $state(NoOp);

  showModal = async (original: OauthClientModel, edit: boolean = false): Promise<OauthClientModel> => {
    promiseReject();

    editMode = edit;
    client = await deepCopy(original);
    originalClient = await deepCopy(original);

    const details = await clients.getClientDetails(client.id).catch((err) => {
      queueNotification(ColorKeys.Danger, `Could not get OAuth 2.0 client details: ${err.message}`);
        return Promise.reject();
      });
    client.client_secret = details.client_secret;

    if (editMode) setTimeout(showCreateModalInternal(), 0);
    else setTimeout(showModalInternal(), 0);

    // TODO: what if we only show but the modal and the close? memory leak?
    return new Promise((resolve, reject) => {
      promiseResolve = ((res) => {
        client = EmptyOauthClient;
        resolve(res);
      });
      promiseReject = ((err) => {
        client = EmptyOauthClient;
        reject(err);
      });
    })
  };
  showCreateModal = () => {
    return showModal(EmptyOauthClient, true);
  }

  let showCreateModalInternal: () => any = $state(NoOp);
  let showModalInternal: () => any = $state(NoOp);

  let editMode: boolean = $state(false);
  let title: string = $derived(client.id ? "Edit OAuth 2.0 Client" : "Register OAuth 2.0 Client");

  const onDelete = async () => {
    await clients.deleteClient(client.id).catch(err => {
      throw new Error(`Could not delete OAuth 2.0 client: ${err.message}`);
    });

    promiseResolve(originalClient);
  };
  const onEdit = async () => {
    if (client.id === "") {
      const newClient = await clients.registerClient(client).catch(err => {
        throw new Error(`Could not register OAuth 2.0 client: ${err.message}`);
      });

      editMode = false;

      promiseResolve(newClient);

      setTimeout(() => {
        showModal(newClient, false);
      }, 50);
    } else {
      const updatedClient = await clients.updateClient(client).catch(err => {
        throw new Error(`Could not update OAuth 2.0 client: ${err.message}`);
      });

      promiseResolve(updatedClient);

      setTimeout(() => {
        showModal(updatedClient, false);
      }, 50);
    }
  };
</script>

<EditableModal
  title={title}
  deleteConfirmation={`Are you sure you want to delete the OAuth 2.0 client?`}
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
  <TextInput bind:value={client.name} name="name" placeholder="Name" editable={editMode}/>
  <TextInput bind:value={client.client_id} name="client_id" placeholder="Client ID" editable={editMode}/>
  <TextInput bind:value={client.client_secret} name="client_secret" placeholder="Client Secret" editable={editMode} password={true}/>
  <TextInput bind:value={client.authorization_url} name="authorization_url" placeholder="Authorization URL" editable={editMode}/>

</EditableModal>
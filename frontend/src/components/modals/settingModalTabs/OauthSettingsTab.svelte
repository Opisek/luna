<script lang="ts">
  import { Pencil, Trash2 } from "lucide-svelte";
  import { queueNotification } from "../../../lib/client/notifications";
  import { ColorKeys } from "../../../types/colors";
  import { UserSettingKeys  } from "../../../types/settings";
  import List from "../../forms/List.svelte";
  import Button from "../../interactive/Button.svelte";
  import IconButton from "../../interactive/IconButton.svelte";
  import type { Settings } from "../../../lib/client/data/settings.svelte";
  import type { OauthClients } from "../../../lib/client/data/oauth.svelte";
  import OauthClientModal from "../OauthClientModal.svelte";

  interface Props {
    settings: Settings;
    clients: OauthClients;
  }

  let {
    settings,
    clients,
  }: Props = $props();

  let showOauthClient = $state<(session: OauthClientModel, editable: boolean) => Promise<OauthClientModel>>(Promise.reject);
  let registerOauthClient = $state<() => Promise<OauthClientModel>>(Promise.reject);

  function deleteClient(id: string) {
    clients.deleteClient(id).catch((err) => {
      queueNotification(ColorKeys.Danger, err);
    });
  }
</script>

<style lang="scss">
  @use "../../../styles/colors.scss";
  @use "../../../styles/dimensions.scss";
  @use "../../../styles/text.scss";

  .client {
    padding: dimensions.$gapMiddle;
    background-color: colors.$backgroundSecondary;
    color: colors.$foregroundSecondary;
    border-radius: dimensions.$borderRadius;

    display: grid;
    gap: dimensions.$gapSmall;
    row-gap: 0;
    grid-template-columns: 1fr auto;
    grid-template-rows: auto auto;
    grid-template-areas: "name buttons" "clientId buttons";
    justify-content: center;
    align-items: center;
  }

  .client.showId {
    grid-template-rows: auto auto auto;
    grid-template-areas: "name buttons" "clientId buttons" "internalId buttons";
  }

  .client > .name {
    grid-area: name;
    display: flex;
    flex-direction: row;
    justify-content: start;
    align-items: center;
    gap: dimensions.$gapTiny;
  }
  .client > .clientId {
    grid-area: clientId;
    font-size: text.$fontSizeSmall;
  }
  .client > .buttons {
    grid-area: buttons;
    display: flex;
    flex-direction: row;
    gap: dimensions.$gapSmall;
  }
  .client > .internalId {
    grid-area: internalId;
    font-size: text.$fontSizeSmall;
  }
  .client:not(.showId) > .internalId {
    display: none;
  }
</style>

<Button color={ColorKeys.Accent} onClick={registerOauthClient}>Register an OAuth 2.0 Client</Button>

{#if clients.clients.length !== 0}
  <List
    label="Registered OAuth 2.0 Client"
    items={clients.clients}
    id={item => item.id}
    template={clientTemplate}
  />
{/if}

{#snippet clientTemplate(client: OauthClientModel)}
  <div class="client" class:showId={settings.userSettings[UserSettingKeys.DebugMode]}>
    <span class="name">
      {client.name}
    </span>

    <span class="clientId">
      {client.client_id}
    </span>

    <div class="buttons">
      <IconButton click={() => showOauthClient(client, true)} alt="Edit">
        <Pencil size={20}/>
      </IconButton>
      <IconButton click={() => deleteClient(client.id)} alt="Delete">
        <Trash2 size={20}/>
      </IconButton>
    </div>

    <span class="internalId">
      ID: {client.id}
    </span>
  </div>
{/snippet}

<OauthClientModal
  bind:showModal={showOauthClient}
  bind:showCreateModal={registerOauthClient}
/>
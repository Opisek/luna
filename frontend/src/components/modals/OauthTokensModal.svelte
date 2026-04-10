<script lang="ts">
  import Modal from "./Modal.svelte";

  import { NoOp } from "$lib/client/placeholders";
  import { getSettings } from "$lib/client/data/settings.svelte";
  import { UserSettingKeys } from "../../types/settings";
  import { getOauthClients } from "$lib/client/data/oauth.svelte";
  import List from "../forms/List.svelte";
  import { fetchJson } from "$lib/client/net";
  import { queueNotification } from "$lib/client/notifications";
  import { ColorKeys } from "../../types/colors";
  import { Check } from "lucide-svelte";
  import IconButton from "../interactive/IconButton.svelte";
  import Paragraph from "../forms/Paragraph.svelte";
  import Button from "../interactive/Button.svelte";
  import Spinner from "../decoration/Spinner.svelte";
  
  interface Props {
    authorize?: (clientId: string) => Promise<string>;
    abort?: () => void;
  }

  let {
    authorize = $bindable(),
    abort = $bindable(),
  }: Props = $props();

  const settings = getSettings();
  const oauthClients = getOauthClients();

  let showModalInternal: () => Promise<string> = $state(Promise.reject);
  let success: (result: string) => void = $state(NoOp);
  let failure: () => void = $state(NoOp);
  let externalAuthPromiseResolve: (result: string) => void = $state(NoOp);
  let externalAuthPromiseReject: () => void = $state(NoOp);

  let selectedClientId = $state("");
  let selectedClient: OauthClientModel | null = $derived(oauthClients.clients.find(client => client.id === selectedClientId) ?? null);
  let oauthRequestId = $state("");

  authorize = async (clientId: string): Promise<string> => {
    await oauthClients.fetch();
    await oauthClients.fetchTokens();

    selectedClientId = clientId;

    if ((oauthClients.clientTokens.get(selectedClientId) || []).length == 0) {
      return authorizeWithExternalProvider().finally(() => {
        window.removeEventListener("storage", oauthAuthorizationResponseListener);
      })
    } else {
      return showModalInternal().finally(() => {
        window.removeEventListener("storage", oauthAuthorizationResponseListener);
      })
    }
  };

  abort = () => {
    externalAuthPromiseReject();
  };

  async function authorizeWithExternalProvider() {
    const json = await fetchJson(`/api/oauth/authorization/${selectedClientId}`, { method: "PUT" }).catch((err: Error) => {
      if (err.message.includes("Service unavailable")) err.message = "Please try again";
      queueNotification(ColorKeys.Danger, err.message)
    });
    if (!json || !json.url || !json.request?.request_id) return Promise.reject();

    oauthRequestId = json.request.request_id;
    localStorage.setItem(`oauth/${oauthRequestId}/expiry`, json.request.expires_at);

    window.addEventListener("storage", oauthAuthorizationResponseListener);

    window.open(json.url, "_blank")?.focus();

    externalAuthPromiseReject();
    return new Promise<string>((resolve, reject) => {
      externalAuthPromiseResolve = resolve;
      externalAuthPromiseReject = reject;
    });
  }
  async function oauthAuthorizationResponseListener() {
    const rawResponse = localStorage.getItem(`oauth/${oauthRequestId}/response`);
  
    if (!rawResponse) return;

    const response = JSON.parse(rawResponse);

    if (!response) return;

    window.removeEventListener("storage", oauthAuthorizationResponseListener);

    localStorage.removeItem(`oauth/${oauthRequestId}/response`);
    localStorage.removeItem(`oauth/${oauthRequestId}/expiry`);

    if (response?.status === "ok") {
      if (response.warnings) {
        for (const warning of response.warnings) {
          queueNotification(ColorKeys.Warning, warning);
        }
      }
      else queueNotification(ColorKeys.Success, `Logged into ${selectedClient?.name} successfully`);
      await oauthClients.fetchTokens();
      externalAuthPromiseResolve(response.token);
    } else if ((response?.error as string || "").toLowerCase().includes("service unavailable")) {
      queueNotification(ColorKeys.Warning, "Please try again");
      externalAuthPromiseReject();
    } else {
      queueNotification(ColorKeys.Danger, response?.error || "Unknown error");
      externalAuthPromiseReject();
    }
  }
</script>

<style lang="scss">
  @use "../../styles/colors.scss";
  @use "../../styles/dimensions.scss";
  @use "../../styles/text.scss";

  .tokens {
    padding: dimensions.$gapMiddle;
    background-color: colors.$backgroundSecondary;
    color: colors.$foregroundSecondary;
    border-radius: dimensions.$borderRadius;

    display: grid;
    gap: dimensions.$gapSmall;
    row-gap: 0;
    grid-template-columns: 1fr auto;
    grid-template-rows: auto;
    grid-template-areas: "name buttons";
    justify-content: center;
    align-items: center;
  }

  .tokens.showId {
    grid-template-rows: auto auto;
    grid-template-areas: "name buttons" "internalId buttons";
  }

  .tokens > .name {
    grid-area: name;
    display: flex;
    flex-direction: row;
    justify-content: start;
    align-items: center;
    gap: dimensions.$gapTiny;
  }
  .tokens > .buttons {
    grid-area: buttons;
    display: flex;
    flex-direction: row;
    gap: dimensions.$gapSmall;
  }
  .tokens > .internalId {
    grid-area: internalId;
    font-size: text.$fontSizeSmall;
  }
  .tokens:not(.showId) > .internalId {
    display: none;
  }
</style>

<Modal
  title={"Choose account"}
  bind:showModal={showModalInternal}
  bind:success
  bind:failure
>
  <Paragraph>
    You are already signed in with {selectedClient?.name}.<br>
    You can choose one of your existing accounts or authorize a new account.
  </Paragraph>

  <List
    label="Authorized Accounts"
    items={oauthClients.clientTokens.get(selectedClientId) || []}
    id={item => item.id}
    template={tokensTemplate}
  />

  <Button color={ColorKeys.Accent} onClick={async () => authorizeWithExternalProvider().then(res => success(res)).catch(NoOp)}>
    Sign into a different account
  </Button>
</Modal>

{#snippet tokensTemplate(tokens: OauthTokensModel)}
  <div class="tokens" class:showId={settings.userSettings[UserSettingKeys.DebugMode]}>
    <span class="name">
      {tokens.account_name}
    </span>

    <div class="buttons">
      <IconButton onClick={() => success(tokens.id)} alt="Use account">
        <Check size={20}/>
      </IconButton>
    </div>

    <span class="internalId">
      ID: {tokens.id}
    </span>
  </div>
{/snippet}
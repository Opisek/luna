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

  let promiseResolve: (value: string | PromiseLike<string>) => void = $state(NoOp);
  let promiseReject: (reason?: any) => void = $state(NoOp);

  let showModalInternal: () => any = $state(NoOp);

  let selectedClientId = $state("");
  let selectedClient: OauthClientModel | null = $derived(oauthClients.clients.find(client => client.id === selectedClientId) ?? null);
  let oauthRequestId = $state("");
  let extAuthPending = $state(false);
  let extAuthRetried = $state(false);
  let rejectOnFail = $state(false);

  authorize = async (clientId: string): Promise<string> => {
    promiseReject();

    await oauthClients.fetch();
    await oauthClients.fetchTokens();

    selectedClientId = clientId;

    if ((oauthClients.clientTokens.get(selectedClientId) || []).length == 0) {
      rejectOnFail = true
      authorizeWithExternalProvider();
    } else {
      rejectOnFail = false;
      setTimeout(showModalInternal(), 0);
    }

    return new Promise((resolve, reject) => {
      promiseResolve = ((res) => {
        window.removeEventListener("storage", oauthAuthorizationResponseListener);
        resolve(res);
      });
      promiseReject = ((err) => {
        window.removeEventListener("storage", oauthAuthorizationResponseListener);
        reject(err);
      });
    })
  };

  abort = () => {
    promiseReject();
  };

  async function authorizeWithExternalProvider() {
    if (extAuthPending) return;
    extAuthPending = true;

    const json = await fetchJson(`/api/oauth/authorization/${selectedClientId}`, { method: "PUT" }).catch((err: Error) => {
      extAuthPending = false;
      if (err.message.includes("Service unavailable")) err.message = "Please try again";
      queueNotification(ColorKeys.Danger, err.message)
    });
    if (!json || !json.url || !json.request?.request_id) return;

    oauthRequestId = json.request.request_id;
    localStorage.setItem(`oauth/${oauthRequestId}/expiry`, json.request.expires_at);

    window.addEventListener("storage", oauthAuthorizationResponseListener);

    window.open(json.url, "_blank")?.focus();
  }

  async function oauthAuthorizationResponseListener() {
    const rawResponse = localStorage.getItem(`oauth/${oauthRequestId}/response`);
  
    if (!rawResponse) return;

    const response = JSON.parse(rawResponse);

    if (!response) return;

    window.removeEventListener("storage", oauthAuthorizationResponseListener);

    localStorage.removeItem(`oauth/${oauthRequestId}/response`);
    localStorage.removeItem(`oauth/${oauthRequestId}/expiry`);

    extAuthPending = false;

    if (response?.status === "ok") {
      if (response.warnings) {
        for (const warning of response.warnings) {
          queueNotification(ColorKeys.Warning, warning);
        }
      }
      else queueNotification(ColorKeys.Success, `Logged into ${selectedClient?.name} successfully`);
      await oauthClients.fetchTokens();
      promiseResolve(response.token);
    } else if (!extAuthRetried && (response?.error as string || "").toLowerCase().includes("service unavailable")) {
      extAuthRetried = true;
      authorizeWithExternalProvider();
    } else {
      queueNotification(ColorKeys.Danger, response?.error || "Unknown error");
      if (rejectOnFail) {
        promiseReject();
      }
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
  onModalHide={() => {
    promiseReject();
  }}
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

  <Button color={ColorKeys.Accent} onClick={authorizeWithExternalProvider} enabled={!extAuthPending}>
    {#if extAuthPending}
      <Spinner/>
    {:else}
      Sign into a different account
    {/if}
  </Button>
</Modal>

{#snippet tokensTemplate(tokens: OauthTokensModel)}
  <div class="tokens" class:showId={settings.userSettings[UserSettingKeys.DebugMode]}>
    <span class="name">
      {tokens.account_name}
    </span>

    <div class="buttons">
      <IconButton click={() => promiseResolve(tokens.id)}>
        <Check size={20}/>
      </IconButton>
    </div>

    <span class="internalId">
      ID: {tokens.id}
    </span>
  </div>
{/snippet}
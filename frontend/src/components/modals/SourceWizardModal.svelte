<script lang="ts">
  import { getContext } from "svelte";
  import Modal from "./Modal.svelte";
  import SelectButtons from "../forms/SelectButtons.svelte";
  import { NoOp } from "../../lib/client/placeholders";
  import TextInput from "../forms/TextInput.svelte";
  import FileUpload from "../forms/FileUpload.svelte";
  import { isValidIcalFile, isValidUrl, valid } from "../../lib/client/validation";
  import Button from "../interactive/Button.svelte";
  import Loader from "../decoration/Loader.svelte";
  import Link from "../forms/Link.svelte";
  import Horizontal from "../layout/Horizontal.svelte";
  import { fetchJson } from "../../lib/client/net";
  import { queueNotification } from "../../lib/client/notifications";
  import { getRepository } from "../../lib/client/data/repository.svelte";
  import { ColorKeys } from "../../types/colors";
  import Paragraph from "../forms/Paragraph.svelte";
  import { getOauthClients } from "../../lib/client/data/oauth.svelte";
  import { sleep } from "../../lib/common/misc";
  import OauthTokensModal from "./OauthTokensModal.svelte";
  import Spinner from "../decoration/Spinner.svelte";
  import IconButton from "../interactive/IconButton.svelte";
  import { Check, X } from "lucide-svelte";
  import { t } from "@sveltia/i18n";

  interface Props {
    showModal?: () => Promise<SourceModel>;
  }

  let {
    showModal = $bindable(),
  }: Props = $props();

  const oauthClients = getOauthClients();

  let showSourceModal: () => Promise<SourceModel> = getContext("showSourceModal");
  let showModalInternal: () => Promise<SourceModel> = $state(Promise.reject);
  let success: (result: SourceModel) => void = $state(NoOp);
  let failure: () => void = $state(NoOp);

  let urlValid: Validity = $state(valid);
  let fileValid: Validity = $state(valid);

  let name: string = $state("");
  let inputType: "link" | "file" | "google" = $state("link");
  let url: string = $state("");
  let urlType: "ical" | "caldav" | "google" | "unknown" = $state("unknown");
  let files: FileList | null = $state(null);
  let needAuth: boolean = $state(false);
  let authType: "none" | "basic" | "bearer" = $state("none");
  let auth: Record<string, string> = $state({
    username: "",
    password: "",
    token: "",
  });

  let googleOauthClient = $derived(oauthClients.clients.find(x => new URL(x.base_url).host === "accounts.google.com"));
  let googleOauthClientAuthorized: boolean = $derived(
    googleOauthClient != null &&
    auth.client_id != null &&
    auth.client_id == googleOauthClient.id &&
    auth.tokens_id != null &&
    auth.tokens_id != "" &&
    oauthClients.tokenClients.get(auth.tokens_id) == auth.client_id
  );

  let oauthPending = $state(false);
  let performOauthAuhorization: (clientId: string) => Promise<string> = $state(async () => "");
  let abortOauthAuthorization: () => void = $state(NoOp);

  showModal = async () => {
    oauthClients.fetch();

    name = "";
    inputType = "link";
    url = "";
    urlType = "unknown";
    needAuth = false;
    authType = "none";
    auth = {
      username: "",
      password: "",
      token: "",
      client_id: "",
      tokens_id: ""
    };
    files = null;

    cachedChecks = new Map<string, Validity>();
    if (checkingUrl) {
      urlValidityResolve(valid);
      clearTimeout(urlValidityTimeout);
      checkingUrl = false;
    }
    lastUrlValidity = valid;

    return showModalInternal();
  }

  let submittable = $derived.by(() => {
    switch (inputType) {
      case "link":
        return urlValid.valid && urlType !== "unknown" && lastUrlValidity.valid && urlValid.valid && !checkingUrl && name !== "";
      case "file":
        return fileValid.valid && files !== null && files.length === 1 && name !== "";
      case "google":
        return googleOauthClientAuthorized && name !== "";
    }
  });

  function advanced() {
    showSourceModal().then((source) => {
      success(source);
    }).catch(NoOp);
  }

  async function save() {
    const source: SourceModel = {
      id: "",
      name: name,
      type: inputType === "google" ? "google" : urlType,
      auth_type: inputType === "google" ? "oauth" : (needAuth ? authType : "none"),
      auth: (needAuth && authType != "none") || inputType === "google" ? auth : {},
      settings: {},
      can_add_calendars: false
    };

    if (inputType === "link") {
      if (urlType === "ical") source.settings.location = "remote";
      source.settings.url = url;
    } else if (inputType === "file") {
      source.type = "ical";
      source.settings.location = "database";
      source.settings.file = files;
    }

    return getRepository().createSource(source).then(() => {
      success(source);
    }).catch(err => {
      queueNotification(ColorKeys.Danger, t("source.error.create", { values: { name: source.name, msg: err.message } }));
    });
  }

  let cachedChecks = $state(new Map<string, Validity>());
  let checkingUrl = $state(false);
  let urlValidityResolve: (value: Validity) => void = $state(NoOp);
  let urlValidityTimeout: ReturnType<typeof setTimeout> | undefined = $state(undefined);
  let lastUrlValidity: Validity = $state(valid);

  function queueUrlCheck() {
    checkUrl(url);
  }

  async function checkUrl(checkUrl: string): Promise<Validity> {
    if (checkingUrl) {
      clearTimeout(urlValidityTimeout);
      urlValidityResolve(lastUrlValidity);
    }
    const cacheKey = JSON.stringify({ url: checkUrl, auth: needAuth == true && authType != "none" ? auth : null });
    if (checkUrl === "") {
      return {
        valid: false,
        message: t("source.error.url.empty"),
      };
    }
    if (needAuth && ((authType === "basic" && (auth.username === "" || auth.password === "")) || (authType === "bearer" && auth.token === ""))) {
      return {
        valid: false,
        message: t("source.error.url.credentials.empty"),
      };
    }
    if (cachedChecks.has(cacheKey)) {
      return cachedChecks.get(cacheKey)!;
    }
    checkingUrl = true;
    return new Promise((resolve) => {
      urlValidityResolve = resolve;
      urlValidityTimeout = setTimeout(async () => {
        const validity = await (async () => {
          urlType = "unknown";

          const syntacticValidity = await isValidUrl(checkUrl);
          if (!syntacticValidity.valid) return syntacticValidity;

          const formData = new FormData();
          formData.append("url", checkUrl);
          if (!needAuth || authType == "none") {
            formData.append("auth_type", "none");
          } else if (authType == "basic") {
            formData.append("auth_type", "basic");
            formData.append("auth_username", auth.username);
            formData.append("auth_password", auth.password);
          } else if (authType == "bearer") {
            formData.append("auth_type", "bearer");
            formData.append("auth_token", auth.token);
          }
          let newCacheKey;
          return fetchJson(`/api/url`, { method: "POST", body: formData }).then((res) => {
            urlType = res.type;
            switch (res.type) {
              default:
              case "unknown":
                if (res.type === "unknown" && res.status == 401) {
                  const message = (needAuth && authType !== "none") ? t("source.error.url.credentials.incorrect") : t("source.error.url.credentials.empty");
                  needAuth = true;
                  if (authType == "none") authType = "basic";
                  return {
                    valid: false,
                    message: message, 
                  }
                } else {
                  return {
                    valid: false,
                    message: t("source.error.url.calendars"),
                  };
                }
              case "ical":
                needAuth = res.auth !== "none";
                newCacheKey = JSON.stringify({ url: checkUrl, auth: needAuth == true && authType != "none" ? auth : null })
                cachedChecks.set(newCacheKey, { valid: true, message: "" }); // prevent second check
                return {
                  valid: true,
                  message: "",
                };
              case "caldav":
                needAuth = res.auth !== "none";
                newCacheKey = JSON.stringify({ url: res.url, auth: needAuth == true && authType != "none" ? auth : null })
                cachedChecks.set(newCacheKey, { valid: true, message: "" }); // prevent second check
                url = res.url;
                return {
                  valid: true,
                  message: "",
                };
            }
          }).catch((err) => {
            queueNotification(ColorKeys.Danger, t("source.error.calendars.find", { values: { msg: err.message } }));
            return {
              valid: false,
              message: t("source.error.url.calendars"),
            }
          });
        })();
        lastUrlValidity = validity;
        urlValid = validity;
        cachedChecks.set(cacheKey, validity);
        checkingUrl = false;
        resolve(validity);
      }, 1000);
    });
  }

  async function startOauthAuthorization() {
    if (!googleOauthClient) return;
    if (oauthPending) return;
    oauthPending = true;

    await sleep(0);
    
    await performOauthAuhorization(googleOauthClient.id).then((id) => {
      auth.client_id = googleOauthClient.id;
      auth.tokens_id = id;
    }).catch(() => {
      queueNotification(ColorKeys.Danger, t("auth.oauth.abort"));
    }).finally(() => {
      oauthPending = false;
    });
  }
</script>

<Modal
  title={t("source.title.wizard")}
  bind:showModal={showModalInternal}
  bind:success
  bind:failure
  onModalHide={() => {
    abortOauthAuthorization();
  }}
>
  <TextInput bind:value={name} name="name" placeholder="Name"/>
  <SelectButtons bind:value={inputType} name="ical_location" placeholder={t("form.wizard.location")} options={[
    {
      value: "link",
      name: t("ical.location.remote"),
    },
    {
      value: "file",
      name: t("ical.location.database"),
    },
    //{
    //  value: "holidays",
    //  name: "Public Holidays",
    //},
    {
      value: "google",
      name: t("google.display"),
    },
  ]}/>

  {#if inputType === "link"}
    <TextInput bind:value={url} name="url" placeholder="URL" validation={checkUrl} bind:validity={urlValid} />
    {#if needAuth}
        <SelectButtons bind:value={authType} name="auth_type" placeholder={t("auth.type")} options={[
        {
          value: "basic",
          name: t("auth.basic.display"),
        },
        {
          value: "bearer",
          name: t("auth.bearer.display"),
        }
      ]}/>
      {#if authType === "basic"}
        <TextInput bind:value={auth.username} onInput={queueUrlCheck} name="auth_username" placeholder={t("auth.basic.username")}/>
        <TextInput bind:value={auth.password} onInput={queueUrlCheck} name="auth_password" placeholder={t("auth.basic.password")} password={true} />
      {:else if authType === "bearer"}
        <TextInput bind:value={auth.token} onInput={queueUrlCheck} name="auth_token" placeholder={t("auth.bearer.token")} password={true} />
      {/if}
    {/if}
    {#if checkingUrl}
      <Horizontal position="center">
        <Loader/>
      </Horizontal>
    {/if}
  {:else if inputType === "file"}
    <FileUpload bind:files={files} name="file" placeholder={t("file.display")} accept=".ical,.ics,.ifb,.icalendar" validation={isValidIcalFile} bind:validity={fileValid} />
  <!--
  {:else if inputType === "holidays"}
    <Paragraph>
      Feature not yet available
    </Paragraph>
  -->
  {:else if inputType === "google"}
    {#if googleOauthClient}
      <Button color={googleOauthClientAuthorized ? ColorKeys.Success : ColorKeys.Accent} onClick={startOauthAuthorization} enabled={!oauthPending && !googleOauthClientAuthorized}>
        {#if oauthPending}
          <Spinner/>
        {:else if googleOauthClientAuthorized}
          Authorized
        {:else}
          Sign in with {googleOauthClient.name}
        {/if}
      </Button>
      {#if googleOauthClientAuthorized}
        <Horizontal position="right">
          <Link onClick={startOauthAuthorization}>{t("auth.oauth.different")}</Link>
        </Horizontal>
      {/if}
    {:else}
      <Paragraph>
        {t("auth.oauth.google.unconfigured")}
      </Paragraph>
    {/if}
  {/if}

  <Horizontal position="right">
    <Link onClick={advanced}>{t("form.advanced")}</Link>
  </Horizontal>
  {#snippet buttons()}
    <IconButton onClick={save} color={ColorKeys.Success} enabled={submittable} type="submit" alt={t("button.save")} canRenderAsButton={true}>
      <Check/>
    </IconButton>
    <IconButton onClick={failure} color={ColorKeys.Danger} alt={t("button.cancel")} canRenderAsButton={true}><X/></IconButton>
  {/snippet}
</Modal>

{#if oauthPending}
  <OauthTokensModal bind:authorize={performOauthAuhorization} bind:abort={abortOauthAuthorization} />
{/if}
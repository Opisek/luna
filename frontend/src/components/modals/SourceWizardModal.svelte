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
  import { getRepository } from "../../lib/client/repository";

  interface Props {
    showModal?: () => Promise<SourceModel>;
  }

  let {
    showModal = $bindable(),
  }: Props = $props();

  let showNewSourceModal: () => Promise<SourceModel> = getContext("showNewSourceModal");
  let showModalInternal = $state(NoOp)
  let hideModalInternal = $state(NoOp)

  let urlValid: Validity = $state(valid);
  let fileValid: Validity = $state(valid);

  let name: string = $state("");
  let inputType: "link" | "file" | "holidays" = $state("link");
  let url: string = $state("");
  let urlType: "ical" | "caldav" | "unknown" = $state("unknown");
  let files: FileList | null = $state(null);
  let needAuth: boolean = $state(false);
  let authType: "none" | "basic" | "bearer" = $state("none");
  let auth: Record<string, string> = $state({
    username: "",
    password: "",
    token: "",
  });

  let promiseResolve: (value: SourceModel) => void = $state(NoOp);
  let promiseReject: (reason?: any) => void = $state(NoOp);

  showModal = async () => {
    awaitingEdit = false;

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
    };
    files = null;

    cachedChecks = new Map<string, Validity>();
    if (checkingUrl) {
      urlValidityResolve(valid);
      clearTimeout(urlValidityTimeout);
      checkingUrl = false;
    }
    lastUrlValidity = valid;

    showModalInternal();

    return new Promise((resolve, reject) => {
      promiseResolve = resolve;
      promiseReject = reject;
    });
  }

  let submittable = $derived.by(() => {
    switch (inputType) {
      case "link":
        return urlValid.valid && urlType !== "unknown" && lastUrlValidity.valid && urlValid.valid && !checkingUrl && name !== "";
      case "file":
        return fileValid.valid && files !== null && files.length === 1 && name !== "";
      case "holidays":
        return false && name !== "";
    }
  });

  function advanced() {
    showNewSourceModal().then((source) => {
      saveInternal(source);
    }).catch(NoOp);
  }

  let awaitingEdit = $state(false);
  function save() {
    const source: SourceModel = {
      id: "",
      name: name,
      type: urlType,
      auth_type: needAuth ? authType : "none",
      auth: needAuth && authType != "none" ? auth : {},
      settings: {}
    };

    if (inputType === "link") {
      if (urlType === "ical") source.settings.location = "remote";
      source.settings.url = url;
    } else if (inputType === "file") {
      source.settings.location = "database";
      source.settings.file = files;
    }

    getRepository().createSource(source).then(() => {
      saveInternal(source);
    }).catch(err => {
      queueNotification("failure", `Could not create source ${source.name}: ${err.message}`);
    });
  }
  async function saveInternal(source: SourceModel) {
    promiseResolve(source);
    hideModalInternal();
  }
  function cancel() {
    promiseReject();
    hideModalInternal();
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
    checkingUrl = true;
    if (checkUrl === "") {
      checkingUrl = false;
      return {
        valid: false,
        message: "URL is required.",
      };
    }
    if (needAuth && ((authType === "basic" && (auth.username === "" || auth.password === "")) || (authType === "bearer" && auth.token === ""))) {
      checkingUrl = false;
      return {
        valid: false,
        message: "Credentials are required to access this URL.",
      };
    }
    if (cachedChecks.has(cacheKey)) {
      checkingUrl = false;
      return cachedChecks.get(cacheKey)!;
    } else {
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
            return fetchJson(`/api/url`, { method: "POST", body: formData }).then((res) => {
              urlType = res.type;
              switch (res.type) {
                default:
                case "unknown":
                  if (res.type === "unknown" && res.status == 401) {
                    const message = (needAuth && authType !== "none") ? "Credentials are incorrect." : "Credentials are required to access this URL.";
                    needAuth = true;
                    if (authType == "none") authType = "basic";
                    return {
                      valid: false,
                      message: message, 
                    }
                  } else {
                    return {
                      valid: false,
                      message: "Could not find calendars. Are you sure this URL is correct?",
                    };
                  }
                case "ical":
                  needAuth = res.auth !== "none";
                  return {
                    valid: true,
                    message: "",
                  };
                case "caldav":
                  needAuth = res.auth !== "none";
                  url = res.url;
                  return {
                    valid: true,
                    message: "",
                  };
              }
            }).catch((err) => {
              queueNotification("failure", `Could not find calendars: ${err.message}`);
              return {
                valid: false,
                message: "Could not find calendars. Are you sure this URL is correct?",
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
  }
</script>

<Modal
  title="Source Wizard"
  bind:showModal={showModalInternal}
  bind:hideModal={hideModalInternal}
>
  <TextInput bind:value={name} name="name" placeholder="Name"/>
  <SelectButtons bind:value={inputType} name="ical_location" placeholder={"What do you want to add?"} options={[
    {
      value: "link",
      name: "Internet Link",
    },
    {
      value: "file",
      name: "Upload File",
    },
    {
      value: "holidays",
      name: "Public Holidays",
    },
  ]}/>

  {#if inputType === "link"}
    <TextInput bind:value={url} name="url" placeholder="URL" validation={checkUrl} bind:validity={urlValid} />
    {#if needAuth}
        <SelectButtons bind:value={authType} name="auth_type" placeholder={"Authentication Type"} options={[
        {
          value: "basic",
          name: "Password",
        },
        {
          value: "bearer",
          name: "Token",
        },
      ]}/>
      {#if authType === "basic"}
        <TextInput bind:value={auth.username} onInput={queueUrlCheck} name="auth_username" placeholder="Username"/>
        <TextInput bind:value={auth.password} onInput={queueUrlCheck} name="auth_password" placeholder="Password" password={true} />
      {/if}
      {#if authType === "bearer"}
        <TextInput bind:value={auth.token} onInput={queueUrlCheck} name="auth_token" placeholder="Token" password={true} />
      {/if}
    {/if}
    {#if checkingUrl}
      <Horizontal position="center">
        <Loader/>
      </Horizontal>
    {/if}
  {:else if inputType === "file"}
    <FileUpload bind:files={files} name="file" placeholder="File" validation={isValidIcalFile} bind:validity={fileValid} />
  {:else if inputType === "holidays"}
      Feature not yet available
  {/if}

  <Horizontal position="right">
    <Link onClick={advanced}>Click to enter advanced mode</Link>
  </Horizontal>
  {#snippet buttons()}
    {#if submittable}
      <Button onClick={save} color="success" enabled={submittable} type="submit">
        {#if awaitingEdit}
          <Loader/>
        {:else}
          Save
        {/if}
      </Button>
    {/if}
    <Button onClick={cancel} color="failure">Cancel</Button>
  {/snippet}
</Modal>
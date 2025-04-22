<script lang="ts">
  import EditableModal from "./EditableModal.svelte";
  import TextInput from "../forms/TextInput.svelte";

  import { EmptySession, NoOp } from "$lib/client/placeholders";
  import { deepCopy } from "$lib/common/misc";
  import { getSettings } from "$lib/client/settings.svelte";
  import { UserSettingKeys } from "../../types/settings";
  import DateTimeInput from "../forms/DateTimeInput.svelte";
  import { getActiveSessions } from "../../lib/client/sessions.svelte";
  import PasswordPromptModal from "./PasswordPromptModal.svelte";
  import { queueNotification } from "../../lib/client/notifications";
  import { ColorKeys } from "../../types/colors";
  
  interface Props {
    showCreateModal?: () => Promise<Session>;
    showModal?: (session: Session) => Promise<Session>;
  }

  let {
    showCreateModal = $bindable(),
    showModal = $bindable(),
  }: Props = $props();

  const settings = getSettings();
  const sessions = getActiveSessions();

  let session: Session = $state(EmptySession);
  let originalSession: Session = $state(EmptySession);

  let promiseResolve: (value: Session | PromiseLike<Session>) => void = $state(NoOp);
  let promiseReject: (reason?: any) => void = $state(NoOp);

  showModal = async (original: Session): Promise<Session> => {
    promiseReject();

    editMode = original.session_id === "";
    editMode = true;
    session = await deepCopy(original);
    originalSession = await deepCopy(original);

    if (editMode) setTimeout(showCreateModalInternal(), 0);
    else setTimeout(showModalInternal(), 0);

    return new Promise((resolve, reject) => {
      promiseResolve = ((res) => {
        session = EmptySession;
        resolve(res);
      });
      promiseReject = ((err) => {
        session = EmptySession;
        reject(err);
      });
    })
  };
  showCreateModal = () => {
    return showModal(EmptySession);
  }

  let showCreateModalInternal: () => any = $state(NoOp);
  let showModalInternal: () => any = $state(NoOp);
  let passwordPrompt = $state<() => Promise<string>>(() => Promise.reject(""));

  let editMode: boolean = $state(false);
  let title: string = $derived(session.session_id ? (editMode ? "Edit API token" : "API Token") : "Create API Token");

  const onDelete = async () => {
    await sessions.deauthorizeSession(session.session_id).catch(err => {
      throw new Error(`Could not delete API token ${session.user_agent}: ${err.message}`);
    });
    promiseResolve(originalSession);
  };
  const onEdit = async () => {
    const password = await passwordPrompt().catch(() => {
      throw new Error(`You must provide your password to create an API token.`);
    });

    if (session.session_id === "") {
      const tokenResponse = await sessions.requestToken(session, password).catch(err => {
        throw new Error(`Could not create API token ${session.user_agent}: ${err.message}`);
      });
      if (tokenResponse === null) return;

      if (await navigator.clipboard.writeText(tokenResponse.token).catch(err => {
        throw new Error(`Could not create API token ${session.user_agent}: ${err.message}`);
      }) === null) return;

      queueNotification(ColorKeys.Success, `API token copied to clipboard.`);
      promiseResolve(tokenResponse.session);
    } else  {
      promiseResolve(await sessions.updateSession(session, password).catch(err => {
        throw new Error(`Could not update API token ${session.user_agent}: ${err.message}`);
      }));
    }
  };

  let canSubmit: boolean = $derived(session && session.user_agent !== "");
</script>

<EditableModal
  title={title}
  deleteConfirmation={`Are you sure you want to delete API token "${session ? session.user_agent : ""}"?`}
  bind:editMode={editMode}
  bind:showCreateModal={showCreateModalInternal}
  bind:showModal={showModalInternal}
  onDelete={onDelete}
  onEdit={onEdit}
  onCancel={promiseReject}
  deletable={true}
  submittable={canSubmit}
>
  <TextInput bind:value={session.user_agent} name="name" placeholder="Name" editable={editMode} />
  {#if session.session_id != ""}
    <DateTimeInput value={session.created_at} allDay={false} placeholder="Creation Date" name="created_at" editable={false}/>
    <DateTimeInput value={session.last_seen} allDay={false} placeholder="Last Activity" name="last_seen" editable={false}/>
    {#if session.session_id && settings.userSettings[UserSettingKeys.DebugMode]}
      <TextInput value={session.session_id} name="id" placeholder="Session ID" editable={false} />
    {/if}
  {/if}
</EditableModal>

{#if canSubmit}
  <PasswordPromptModal bind:prompt={passwordPrompt}/>
{/if}
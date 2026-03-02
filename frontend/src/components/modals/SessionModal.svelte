<script lang="ts">
  import EditableModal from "./EditableModal.svelte";
  import TextInput from "../forms/TextInput.svelte";

  import { EmptySession, NoOp } from "$lib/client/placeholders";
  import { deepCopy } from "$lib/common/misc";
  import { getSettings } from "$lib/client/data/settings.svelte";
  import { UserSettingKeys } from "../../types/settings";
  import DateTimeInput from "../forms/DateTimeInput.svelte";
  import { getActiveSessions } from "../../lib/client/data/sessions.svelte";
  import PasswordPromptModal from "./PasswordPromptModal.svelte";
  import { queueNotification } from "../../lib/client/notifications";
  import { ColorKeys } from "../../types/colors";
  import { PermissionKeys } from "../../types/permissions";
  import ToggleInput from "../forms/ToggleInput.svelte";
  import SectionDivider from "../layout/SectionDivider.svelte";
  
  interface Props {
    showCreateModal?: () => Promise<Session>;
    showModal?: (session: Session, editable: boolean) => Promise<Session>;
  }

  let {
    showCreateModal = $bindable(),
    showModal = $bindable(),
  }: Props = $props();

  const settings = getSettings();
  const sessions = getActiveSessions();

  let session: Session = $state(EmptySession);
  let originalSession: Session = $state(EmptySession);

  let permissions = $state(Object.fromEntries(Object.values(PermissionKeys).map(x => [x, false]))) as Record<PermissionKeys, boolean>;

  let promiseResolve: (value: Session | PromiseLike<Session>) => void = $state(NoOp);
  let promiseReject: (reason?: any) => void = $state(NoOp);

  showModal = async (original: Session, editable: boolean = false): Promise<Session> => {
    promiseReject();

    editMode = editable;
    session = await deepCopy(original);
    originalSession = await deepCopy(original);
    permissions = {} as Record<PermissionKeys, boolean>;
    for (const permission of Object.values(PermissionKeys)) {
      permissions[permission] = false
    }
    if (session.session_id !== "") {
      const sessionPermissions = await sessions.getSessionPermissions(session.session_id).catch((err) => {
        queueNotification(ColorKeys.Danger, `Could not fetch permissions for session ${session.user_agent}: ${err.message}`);
        return [];
      });
      for (const permission of sessionPermissions) {
        permissions[permission as PermissionKeys] = true;
      }
    } else {
      for (const permission of session.permissions) {
        permissions[permission as PermissionKeys] = true;
      }
    }

    if (editMode) {
      session.is_api = true;
      setTimeout(showCreateModalInternal(), 0);
    } else setTimeout(showModalInternal(), 0);

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
    return showModal(EmptySession, true);
  }

  let showCreateModalInternal: () => any = $state(NoOp);
  let showModalInternal: () => any = $state(NoOp);
  let passwordPrompt = $state<() => Promise<string>>(() => Promise.reject(""));

  let editMode: boolean = $state(false);
  let title: string = $derived((session.session_id ? (editMode ? "Edit " : "") : "Create ") + (session.is_api ? "API Token" : "Session"));

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

    session.permissions = Object.entries(permissions)
      .filter(([_, permitted]) => permitted)
      .map(([key, _]) => key as PermissionKeys);
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
  deleteConfirmation={`Are you sure you want to ${session.is_api ? "delete API token" : "deauthorize session"} "${session ? " " + session.user_agent : ""}"?`}
  bind:editMode={editMode}
  bind:showCreateModal={showCreateModalInternal}
  bind:showModal={showModalInternal}
  onDelete={onDelete}
  onEdit={onEdit}
  onCancel={promiseReject}
  editable={editMode}
  deletable={true}
  submittable={canSubmit}
>
  <TextInput bind:value={session.user_agent} name="user_agent" placeholder={session.is_api ? "Name" : "User Agent"} editable={editMode} />
  {#if session.session_id != ""}
    <DateTimeInput value={session.created_at} allDay={false} placeholder="Creation Date" name="created_at" editable={false}/>
    <DateTimeInput value={session.last_seen} allDay={false} placeholder="Last Activity" name="last_seen" editable={false}/>
    <TextInput value={session.initial_ip_address} name="initial_ip_address" placeholder="Initial IP Address" editable={false} />
    <TextInput value={session.last_ip_address} name="last_ip_address" placeholder="Last IP Address" editable={false} />
    {#if session.session_id && settings.userSettings[UserSettingKeys.DebugMode]}
      <TextInput value={session.session_id} name="id" placeholder="Session ID" editable={false} />
    {/if}
  {/if}
  {#if session.is_api}
    <SectionDivider title="Permissions"/>
    {#each Object.values(PermissionKeys) as permission}
      {@const permissionName = permission.split('_').map(word => word.charAt(0).toUpperCase() + word.slice(1)).join(' ')}
      <ToggleInput
        name={permission} 
        description={permissionName}
        bind:value={permissions[permission]}
      />
    {/each}
  {/if}
</EditableModal>

{#if canSubmit}
  <PasswordPromptModal bind:prompt={passwordPrompt}/>
{/if}
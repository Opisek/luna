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
  import { t } from "@sveltia/i18n";
  
  interface Props {
    showModal?: (initial?: Session, edit?: boolean) => Promise<Session>;
  }

  let {
    showModal = $bindable(),
  }: Props = $props();

  let showModalInternal: (initial?: Session, edit?: boolean) => Promise<Session> = $state(Promise.reject);
  let passwordPrompt: () => Promise<string> = $state(Promise.reject);

  const settings = getSettings();
  const sessions = getActiveSessions();

  let session: Session = $state(EmptySession);
  let originalSession: Session = $state(EmptySession);

  let permissions = $state(Object.fromEntries(Object.values(PermissionKeys).map(x => [x, false]))) as Record<PermissionKeys, boolean>;

  showModal = async (initial?: Session, edit?: boolean): Promise<Session> => {
    permissions = {} as Record<PermissionKeys, boolean>;
    for (const permission of Object.values(PermissionKeys)) {
      permissions[permission] = false
    }

    if (!initial) {
      session = await deepCopy(EmptySession)
      originalSession = await deepCopy(EmptySession)
      session.is_api = true;

      for (const permission of session.permissions) {
        permissions[permission as PermissionKeys] = true;
      }
    } else {
      session = await deepCopy(initial);
      originalSession = await deepCopy(initial);
      if (edit) {
        session.is_api = true;
        const sessionPermissions = await sessions.getSessionPermissions(session.id).catch((err) => {
          queueNotification(ColorKeys.Danger, t("session.error.permissions", { values: { name: session.user_agent, msg: err.message } }));
          return [];
        });
        for (const permission of sessionPermissions) {
          permissions[permission as PermissionKeys] = true;
        }
      }
    }

    return showModalInternal();
  };

  let editMode: boolean = $state(false);
  let title: string = $derived(t(`session.title.${(session.id ? (editMode ? "edit" : "view") : "create")}.${session.is_api ? "api" : "user"}`));

  const onDelete = async () => {
    return sessions.deauthorizeSession(session.id).then(() => session).catch(err => {
      throw new Error(t("session.error.delete", { values: { name: session.user_agent, msg: err.message } }));
    });
  };
  const onEdit: () => Promise<Session> = async () => {
    const password = await passwordPrompt().catch(() => {
      throw new Error(t("session.error.password"));
    });

    session.permissions = Object.entries(permissions)
      .filter(([_, permitted]) => permitted)
      .map(([key, _]) => key as PermissionKeys);

    if (session.id === "") {
      const tokenResponse = await sessions.requestToken(session, password).catch(err => {
        throw new Error(t("session.error.create", { values: { name: session.user_agent, msg: err.message } }));
      });

      await navigator.clipboard.writeText(tokenResponse.token).catch(err => {
        throw new Error(t("session.error.clipboard", { values: { name: session.user_agent, msg: err.message } }));
      });

      queueNotification(ColorKeys.Success, t("session.clipboard"));
      return session;
    } else {
      return sessions.updateSession(session, password).then(() => session).catch(err => {
        throw new Error(t("session.error.edit", { values: { name: session.user_agent, msg: err.message } }));
      });
    }
  };

  let canSubmit: boolean = $derived(session && session.user_agent !== "");
</script>

<EditableModal
  title={title}
  deleteConfirmation={t(`session.confirm.delete.${session.is_api ? "api" : "user"}`, { values: { name: session.user_agent } })}
  bind:editMode={editMode}
  bind:showModal={showModalInternal}
  onDelete={onDelete}
  onEdit={onEdit}
  editable={editMode}
  deletable={true}
  submittable={canSubmit}
>
  <TextInput bind:value={session.user_agent} name="user_agent" placeholder={session.is_api ? t("session.agent.api") : t("session.agent.user")} editable={editMode} />
  {#if session.id != ""}
    <DateTimeInput value={session.created_at} allDay={false} placeholder={t("session.date.creation")} name="created_at" editable={false}/>
    <DateTimeInput value={session.last_seen} allDay={false} placeholder={t("session.date.latest")} name="last_seen" editable={false}/>
    <TextInput value={session.initial_ip_address} name="initial_ip_address" placeholder={t("session.ip.initial")} editable={false} />
    <TextInput value={session.last_ip_address} name="last_ip_address" placeholder={t("session.ip.latest")} editable={false} />
    {#if session.id && settings.userSettings[UserSettingKeys.DebugMode]}
      <TextInput value={session.id} name="id" placeholder={t("session.id.display")} editable={false} />
    {/if}
  {/if}
  {#if session.is_api}
    <SectionDivider title={t("session.permissions.subtitle")}/>
    {#each Object.values(PermissionKeys) as permission}
      {@const permissionName = t(`perms.${permission}`)}
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
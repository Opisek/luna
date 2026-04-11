<script lang="ts">
  import { BadgeCheck, Ban, Eye, Shield, Trash2, UserRoundPlus, UserRoundX } from "lucide-svelte";
  import type { RegistrationInvites } from "../../../lib/client/data/invites.svelte";
  import type { Users } from "../../../lib/client/data/users.svelte";
  import { queueNotification } from "../../../lib/client/notifications";
  import { ColorKeys } from "../../../types/colors";
  import { UserSettingKeys, type UserData } from "../../../types/settings";
  import List from "../../forms/List.svelte";
  import Button from "../../interactive/Button.svelte";
  import IconButton from "../../interactive/IconButton.svelte";
  import Image from "../../layout/Image.svelte";
  import Tooltip from "../../interactive/Tooltip.svelte";
  import type { Settings } from "../../../lib/client/data/settings.svelte";
  import RegistrationInviteModal from "../RegistrationInviteModal.svelte";
  import { NoOp } from "../../../lib/client/placeholders";
  import { t } from "@sveltia/i18n";

  interface Props {
    today: Date;
    settings: Settings;
    invites: RegistrationInvites;
    users: Users;
    showConfirmation: (message: string, confirmText?: string, cancelText?: string) => Promise<void>;
    deleteAccount: (id: string) => void;
  }

  let {
    today,
    settings,
    invites,
    users,
    showConfirmation,
    deleteAccount
  }: Props = $props();

  let showRegistrationInviteModal: (initial?: RegistrationInvite, edit?: boolean) => Promise<RegistrationInvite> = $state(Promise.reject);

  function deleteInvite(id: string) {
    invites.revokeInvite(id).catch((err) => {
      queueNotification(ColorKeys.Danger, err);
    });
  }

  function deleteInvites() {
    invites.revokeInvites().catch((err) => {
      queueNotification(ColorKeys.Danger, err);
    });
  }

  async function disableAccount(id: string) {
    await showConfirmation(t("settings.users.confirm.disable"), t("settings.users.info.preserve")).then(async () => {
      users.disableUser(id).catch((err) => {
        queueNotification(ColorKeys.Danger, t("settings.users.error.disable", { values: { msg: err.message } }))
      })
    }).catch(NoOp);
  }

  async function enableAccount(id: string) {
    await showConfirmation(t("settings.users.confirm.enable")).then(async () => {
      users.enableUser(id).catch((err) => {
        queueNotification(ColorKeys.Danger, t("settings.users.error.enable", { values: { msg: err.message } }))
      })
    }).catch(NoOp);
  }
</script>

<style lang="scss">
  @use "../../../styles/colors.scss";
  @use "../../../styles/dimensions.scss";
  @use "../../../styles/text.scss";

  .invite {
    padding: dimensions.$gapMiddle;
    background-color: colors.$backgroundSecondary;
    color: colors.$foregroundSecondary;
    border-radius: dimensions.$borderRadius;

    display: grid;
    gap: dimensions.$gapSmall;
    row-gap: 0;
    grid-template-columns: 1fr auto;
    grid-template-rows: auto auto;
    grid-template-areas: "expiry buttons" "details buttons";
    justify-content: center;
    align-items: center;
  }

  .invite.showId {
    grid-template-rows: auto auto auto;
    grid-template-areas: "expiry buttons" "details buttons" "id buttons";
  }

  .invite > .expiry {
    grid-area: expiry;
  }
  .invite > .details {
    grid-area: details;
    font-size: text.$fontSizeSmall;
  }
  .invite > .buttons {
    grid-area: buttons;
    display: flex;
    flex-direction: row;
    gap: dimensions.$gapSmall;
  }
  .invite > .id {
    grid-area: id;
    font-size: text.$fontSizeSmall;
  }
  .invite:not(.showId) > .id {
    display: none;
  }

  .user {
    padding: dimensions.$gapMiddle;
    background-color: colors.$backgroundSecondary;
    color: colors.$foregroundSecondary;
    border-radius: dimensions.$borderRadius;

    display: grid;
    gap: dimensions.$gapSmall;
    row-gap: 0;
    grid-template-columns: auto 1fr auto;
    grid-template-rows: auto auto auto;
    grid-template-areas: "profilePicture username buttons" "profilePicture details buttons" "profilePicture date buttons";
    justify-content: center;
    align-items: center;
  }

  .user.showId {
    grid-template-rows: auto auto auto auto;
    grid-template-areas: "profilePicture username buttons" "profilePicture details buttons" "profilePicture date buttons" "profilePicture id buttons";
  }

  .user > .profilePicture {
    grid-area: profilePicture;
    display: flex;
    height: 100%;
    align-items: center;
  }
  .user > .username {
    grid-area: username;
    display: flex;
    flex-direction: row;
    justify-content: start;
    align-items: center;
    gap: dimensions.$gapTiny;
  }
  .user > .details {
    grid-area: details;
    font-size: text.$fontSizeSmall;
  }
  .user > .date {
    grid-area: date;
    font-size: text.$fontSizeSmall;
  }
  .user > .buttons {
    grid-area: buttons;
    display: flex;
    flex-direction: row;
    gap: dimensions.$gapSmall;
  }
  .user > .id {
    grid-area: id;
    font-size: text.$fontSizeSmall;
  }
  .user:not(.showId) > .id {
    display: none;
  }

  .user.active {
    background-color: colors.$backgroundAccent;
    color: colors.$foregroundAccent;
  } 
</style>

<Button color={ColorKeys.Accent} onClick={showRegistrationInviteModal}>{t("settings.users.action.invite.new")}</Button>

{#if invites.activeInvites.length !== 0}
  <List
    label={t("settings.users.list.invites")}
    items={invites.activeInvites}
    id={item => item.id}
    template={inviteTemplate}
  />

  <Button color={ColorKeys.Danger} onClick={deleteInvites}>{t("settings.users.action.invite.delete")}</Button>
{/if}

<List
  label={t("settings.users.list.users")}
  items={users.users}
  id={item => item.id}
  template={userTemplate}
/>

{#snippet inviteTemplate(invite: RegistrationInvite)}
  {@const expiresToday = invite.expires_at.getDate() == today.getDate() && invite.expires_at.getMonth() == today.getMonth() && invite.expires_at.getFullYear() == today.getFullYear()}

  {@const hoursRemaining = Math.floor((invite.expires_at.getTime() - today.getTime()) / (1000 * 60 * 60))}
  {@const minutesRemaining = Math.floor((invite.expires_at.getTime() - today.getTime()) / (1000 * 60)) % 60}

  <div class="invite" class:showId={settings.userSettings[UserSettingKeys.DebugMode]}>
    <span class="expiry">
      {t(`invite.date.expiry.${expiresToday ? "today" : "elsewhen"}`, { values: { date: invite.expires_at, hours: hoursRemaining, minutes: minutesRemaining } })}
    </span>

    <span class="details">
      {t("invite.date.creation.inline", { values: { date: invite.created_at, name: users.users.filter(x => x.id == invite.author)[0]?.username || invite.author } })}
      {#if invite.email != ""}
        •
        {invite.email}
      {/if}
    </span>

    <div class="buttons">
      <IconButton onClick={async () => showRegistrationInviteModal(invite)} alt={t("button.details")}>
        <Eye size={20}/>
      </IconButton>
      <IconButton onClick={async () => deleteInvite(invite.id)} color={ColorKeys.Danger} alt={t("button.delete")}>
        <Trash2 size={20}/>
      </IconButton>
    </div>

    <span class="id">
      {t("invite.id.inline", { values: { id: invite.id } })}
    </span>
  </div>
{/snippet}

{#snippet userTemplate(u: UserData)}
  {@const isActive=u.id === users.currentUser}

  <div class="user" class:active={isActive} class:showId={settings.userSettings[UserSettingKeys.DebugMode]}>
    <div class="profilePicture">
      <Image
        src={u.profile_picture}
        alt={t("user.pfp.alt", { values: { name: u.profile_picture } })}
        small={true}
      />
    </div>

    <span class="username">
      {u.username}
      {#if u.admin}
        <Tooltip inheritColor={true} tight={true}>
          {#snippet icon()}
            <Shield size={12}/>
          {/snippet}
          {t("user.status.admin")}
        </Tooltip>
      {/if}
      {#if u.verified}
        <Tooltip inheritColor={true} tight={true}>
          {#snippet icon()}
            <BadgeCheck size={12}/>
          {/snippet}
          {t("user.status.verified")}
        </Tooltip>
      {/if}
      {#if !u.enabled}
        <Tooltip inheritColor={true} tight={true}>
          {#snippet icon()}
            <Ban size={12}/>
          {/snippet}
          {t("user.status.disabled")}
        </Tooltip>
      {/if}
    </span>

    <span class="details">
      {u.email}
    </span>

    <span class="date">
      {t("user.creation", { values: { date: u.created_at } })}
    </span>

    <div class="buttons">
      {#if !u.admin}
        <IconButton onClick={async () => { u.enabled ? disableAccount(u.id) : enableAccount(u.id) }} alt={ u.enabled ? t("user.action.disable") : t("user.action.enable")}>
          {#if u.enabled}
            <UserRoundX size={20}/>
          {:else}
            <UserRoundPlus size={20}/>
          {/if}
        </IconButton>
      {/if}
      {#if !isActive}
        <IconButton onClick={async () => deleteAccount(u.id)} color={ColorKeys.Danger} alt={t("button.delete")}>
          <Trash2 size={20}/>
        </IconButton>
      {/if}
    </div>

    <span class="id">
      {t("user.id.inline", { values: { id: u.id } })}
    </span>
  </div>
{/snippet}

<RegistrationInviteModal
  bind:showModal={showRegistrationInviteModal}
/>
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

  interface Props {
    today: Date;
    settings: Settings;
    invites: RegistrationInvites;
    users: Users;
    showConfirmation: (message: string, onConfirm: () => Promise<void>, confirmText?: string, onCancel?: () => Promise<void>, cancelText?: string) => void;
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

  let showRegistrationInvite = $state<(session: RegistrationInvite, editable: boolean) => Promise<RegistrationInvite>>(Promise.reject);
  let issueRegistrationInvite = $state<() => Promise<RegistrationInvite>>(Promise.reject);

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

  function disableAccount(id: string) {
    showConfirmation("Are you sure you want to disable this account?\nThe user will no longer be able to log in.\nThe user's account will remain intact.", async () => {
      users.disableUser(id).catch((err) => {
        queueNotification(ColorKeys.Danger, `Could not disable user account: ${err.message}`)
      })
    });
  }

  function enableAccount(id: string) {
    showConfirmation("Are you sure you want to enable this account?\nThe user will be able log in again.", async () => {
      users.enableUser(id).catch((err) => {
        queueNotification(ColorKeys.Danger, `Could not enable user account: ${err.message}`)
      })
    });
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

<Button color={ColorKeys.Accent} onClick={issueRegistrationInvite}>Invite a user</Button>

{#if invites.activeInvites.length !== 0}
  <List
    label="Active Invites"
    items={invites.activeInvites}
    id={item => item.invite_id}
    template={inviteTemplate}
  />

  <Button color={ColorKeys.Danger} onClick={deleteInvites}>Delete all invites</Button>
{/if}

<List
  label="Registered Users"
  items={users.users}
  id={item => item.id}
  template={userTemplate}
/>

{#snippet inviteTemplate(invite: RegistrationInvite)}
  {@const expiresToday = invite.expires_at.getDate() == today.getDate() && invite.expires_at.getMonth() == today.getMonth() && invite.expires_at.getFullYear() == today.getFullYear()}

  {@const hoursRemaining = Math.floor((invite.expires_at.getTime() - today.getTime()) / (1000 * 60 * 60))}
  {@const minutesRemaining = Math.floor((invite.expires_at.getTime() - today.getTime()) / (1000 * 60)) % 60}
  {@const daysRemaining = Math.floor((invite.expires_at.getTime() - today.getTime()) / (1000 * 60 * 60 * 24))}

  {@const expiresString = expiresToday ? `at ${invite.expires_at.toLocaleTimeString()}` : `${invite.expires_at.toLocaleDateString()} at ${invite.expires_at.toLocaleTimeString()}`}
  {@const expiresDetailed = expiresToday ? ` (${hoursRemaining == 0 ? `${minutesRemaining} minutes left` : `${hoursRemaining} ${hoursRemaining == 1 ? "hour" : "hours"} left`})` : ""}

  <div class="invite" class:showId={settings.userSettings[UserSettingKeys.DebugMode]}>
    <span class="expiry">
      Expires {expiresString}{expiresDetailed}
    </span>

    <span class="details">
      Created {invite.created_at.toLocaleDateString()} at {invite.created_at.toLocaleTimeString()} by {users.users.filter(x => x.id == invite.author)[0]?.username || invite.author}
      {#if invite.email != ""}
        •
        {invite.email}
      {/if}
    </span>

    <div class="buttons">
      <IconButton click={() => showRegistrationInvite(invite, false)}>
        <Eye size={20}/>
      </IconButton>
      <IconButton click={() => deleteInvite(invite.invite_id)}>
        <Trash2 size={20}/>
      </IconButton>
    </div>

    <span class="id">
      ID: {invite.invite_id}
    </span>
  </div>
{/snippet}

{#snippet userTemplate(u: UserData)}
  {@const isActive=u.id === users.currentUser}

  <div class="user" class:active={isActive} class:showId={settings.userSettings[UserSettingKeys.DebugMode]}>
    <div class="profilePicture">
      <Image
        src={u.profile_picture}
        alt={`Profile picture of user ${u.username}`}
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
          Administrator
        </Tooltip>
      {/if}
      {#if u.verified}
        <Tooltip inheritColor={true} tight={true}>
          {#snippet icon()}
            <BadgeCheck size={12}/>
          {/snippet}
          Verified E-Mail Address
        </Tooltip>
      {/if}
      {#if !u.enabled}
        <Tooltip inheritColor={true} tight={true}>
          {#snippet icon()}
            <Ban size={12}/>
          {/snippet}
          Account disabled
        </Tooltip>
      {/if}
    </span>

    <span class="details">
      {u.email}
    </span>

    <span class="date">
      Created {u.created_at.toLocaleDateString()} at {u.created_at.toLocaleTimeString()}
    </span>

    <div class="buttons">
      {#if !u.admin}
        <IconButton click={() => { u.enabled ? disableAccount(u.id) : enableAccount(u.id) }} info={ u.enabled ? "Disable account" : "Enable account"}>
          {#if u.enabled}
            <UserRoundX size={20}/>
          {:else}
            <UserRoundPlus size={20}/>
          {/if}
        </IconButton>
      {/if}
      {#if !isActive}
        <IconButton click={() => deleteAccount(u.id)}>
          <Trash2 size={20}/>
        </IconButton>
      {/if}
    </div>

    <span class="id">
      ID: {u.id}
    </span>
  </div>
{/snippet}

<RegistrationInviteModal
  bind:showModal={showRegistrationInvite}
  bind:showCreateModal={issueRegistrationInvite}
/>
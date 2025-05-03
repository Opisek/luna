<script lang="ts">
  import { Bot, Code, Eye, Gamepad2, Info, Laptop, LockKeyhole, LogOut, Microchip, Palette, Pencil, RectangleGoggles, RefreshCw, Smartphone, Tablet, Trash, Trash2, TriangleAlert, TvMinimal, User, Users, Watch } from "lucide-svelte";
  import { NoOp } from "../../lib/client/placeholders";
  import ButtonList from "../forms/ButtonList.svelte";
  import Modal from "./Modal.svelte";
  import TextInput from "../forms/TextInput.svelte";
  import { isValidEmail, isValidPassword, isValidRepeatPassword, isValidUsername, valid } from "../../lib/client/validation";
  import ToggleInput from "../forms/ToggleInput.svelte";
  import SelectButtons from "../forms/SelectButtons.svelte";
  import Image from "../layout/Image.svelte";
  import FileUpload from "../forms/FileUpload.svelte";
  import Horizontal from "../layout/Horizontal.svelte";
  import SelectInput from "../forms/SelectInput.svelte";
  import { GlobalSettingKeys, UserSettingKeys, type GlobalSettings, type UserData, type UserSettings } from "../../types/settings";
  import { getSettings } from "../../lib/client/data/settings.svelte";
  import { getSha256Hash } from "../../lib/common/crypto";
  import { deepCopy, deepEquality } from "$lib/common/misc";
  import Button from "../interactive/Button.svelte";
  import Loader from "../decoration/Loader.svelte";
  import { fetchFileById, fetchJson, fetchResponse } from "$lib/client/net";
  import { queueNotification } from "$lib/client/notifications";
  import { ColorKeys } from "../../types/colors";
  import type { Option } from "../../types/options";
  import SliderInput from "../forms/SliderInput.svelte";
  import ConfirmationModal from "./ConfirmationModal.svelte";
  import { clearSession, getActiveSessions } from "../../lib/client/data/sessions.svelte";
  import Tooltip from "../interactive/Tooltip.svelte";
  import List from "../forms/List.svelte";
  import { UAParser } from "ua-parser-js";
  import IconButton from "../interactive/IconButton.svelte";
  import SessionModal from "./SessionModal.svelte";
  import PasswordPromptModal from "./PasswordPromptModal.svelte";
  import SectionDivider from "../layout/SectionDivider.svelte";
  import { getTheme } from "../../lib/client/data/theme.svelte";
  import RegistrationInviteModal from "./RegistrationInviteModal.svelte";
  import { getRegistrationInvites } from "../../lib/client/data/invites.svelte";

  interface Props {
    showModal?: () => any;
    hideModal?: () => any;
  }

  let {
    showModal = $bindable(),
    hideModal = $bindable(NoOp),
  }: Props = $props();

  const settings = getSettings();
  const sessions = getActiveSessions();
  const invites = getRegistrationInvites();
  const today = new Date();

  // Loading
  let loaderAnimation = $state(false);
  function forceRefresh() {
    loaderAnimation = true;
    fetchThemes();
    fetchFonts();
    sessions.fetch();
    invites.fetch();
    settings.fetchSettings().then(() => {
      snapshotSettings();
      refetchProfilePicture();
    });
  }

  showModal = () => {
    saving = false;
    snapshotSettings();
    forceRefresh();
    showModalInternal();
  };

  hideModal = () => {
    hideModalInternal();
  };

  function refetchProfilePicture() {
    const profilePictureUrl = settings.userData.profile_picture || "";

    switch (profilePictureType) {
      case "gravatar":
        profilePictureGravatarForceDefault = profilePictureUrl.includes("f=y");
        break;
      case "database":
        profilePictureFileId = profilePictureUrl.match(/\/api\/files\/([a-f0-9-]{36})/)?.[1] || "";
        if (profilePictureFileId != "") {
          fetchFileById(profilePictureFileId).then(fileList => {
            profilePictureFiles = fileList;
          }).catch(err => {
            queueNotification(ColorKeys.Danger, `Could not download profile picture: ${err.message}`);
          });
        }
        break;
    }
  }

  let showModalInternal = $state(NoOp);
  let hideModalInternal = $state(NoOp);

  const categoriesAdmin: Option<string>[][] = [
    [
      { name: "Account", value: "account", icon: User },
      { name: "Appearance", value: "appearance", icon: Palette },
      { name: "Developer", value: "developer", icon: Code }
    ],
    [
      { name: "Users", value: "users", icon: Users },
      { name: "Administrative", value: "admin", icon: LockKeyhole },
    ],
    [
      { name: "Danger Zone", value: "danger", icon: TriangleAlert, color: ColorKeys.Danger },
      { name: "Logout", value: "logout", icon: LogOut, color: ColorKeys.Danger },
    ],
  ]
  const categories: Option<string>[][] = [
    [
      { name: "Account", value: "account", icon: User },
      { name: "Appearance", value: "appearance", icon: Palette },
      { name: "Developer", value: "developer", icon: Code }
    ],
    [
      { name: "Danger Zone", value: "danger", icon: TriangleAlert, color: ColorKeys.Danger },
      { name: "Logout", value: "logout", icon: LogOut, color: ColorKeys.Danger },
    ],
  ]

  let selectedCategory = $state("account");
  let previousCategory = $state("account");
  $effect(() => {
    if (selectedCategory === previousCategory) return;
    previousCategory = selectedCategory;
    restoreSettings();
  });

  // Themes and Fonts
  let lightThemes = $state<Option<string>[]>([{ name: "Luna Light", value: "luna-light" }]);
  let darkThemes = $state<Option<string>[]>([{ name: "Luna Dark", value: "luna-dark" }]);
  let fonts = $state<Option<string>[]>([
    { name: "Atkinson Hyperlegible Next", value: "atkinson-hyperlegible-next" },
    { name: "Atkinson Hyperlegible Mono", value: "atkinson-hyperlegible-next" }
  ]);

  function formatInstalledFile(rawName: string): Option<string> {
    const formattedName = rawName
      .split("-")
      .map(x => x.charAt(0).toUpperCase() + x.slice(1))
      .join(" ");

    return { name: formattedName, value: rawName };
  }

  function fetchThemes() {
    fetchJson("/installed/themes").then((response) => {
      lightThemes = Object.keys(response.light).map(formatInstalledFile).sort((a, b) => a.name.localeCompare(b.name));
      darkThemes = Object.keys(response.dark).map(formatInstalledFile).sort((a, b) => a.name.localeCompare(b.name));
    }).catch((err) => {
      queueNotification(ColorKeys.Danger, "Failed to fetch themes: " + err);
    });
  }

  function fetchFonts() {
    fetchJson("/installed/fonts").then((response) => {
      fonts = Object.keys(response).map(formatInstalledFile).sort((a, b) => a.name.localeCompare(b.name));
    }).catch((err) => {
      queueNotification(ColorKeys.Danger, "Failed to fetch fonts: " + err);
    });
  }

  // Account Settings
  let newPassword = $state("");
  let passwordPrompt = $state<() => Promise<string>>(() => Promise.reject(""));

  let profilePictureType = $state("gravatar");
  let profilePictureFiles: FileList | null = $state(null);
  let profilePictureFileId = $state("");
  let profilePictureRemoteUrl = $state("");
  let profilePictureGravatarIsDefault = $state(true);
  let profilePictureGravatarForceDefault = $state(false);
  let profilePictureGravatarUrlTrue = $derived.by(() => {
    const email = settings.userData.email || "";
    const trimmedLowercaseEmail = email.trim().toLowerCase();
    const emailHash = getSha256Hash(trimmedLowercaseEmail);
    return `https://www.gravatar.com/avatar/${emailHash}?d=identicon`;
  })
  let profilePictureGravatarUrl = $derived.by(() => {
    return `${profilePictureGravatarUrlTrue}${profilePictureGravatarForceDefault && !profilePictureGravatarIsDefault ? "&f=y" : ""}`;
  })

  let effectiveProfilePictureSource = $state("");
  
  $effect(() => {
    (async () => {
      if (profilePictureType === "gravatar") return profilePictureGravatarUrl;
      else if (profilePictureType === "database") {
        if (profilePictureFileId != "") return `/api/files/${profilePictureFileId}`
        else if (profilePictureFiles) {
          const file = profilePictureFiles[0];
          const reader = new FileReader();
          return new Promise<string>((resolve) => {
            reader.onload = () => {
              resolve(reader.result as string);
            };
            reader.readAsDataURL(file);
          });
        }
        else return "";
      }
      else if (profilePictureType === "remote") return profilePictureRemoteUrl;
      else return "";
    })().then((url) => {
      effectiveProfilePictureSource = url;
    });
  });

  $effect(() => {
    const loadedProfilePictureUrl = settings.userData.profile_picture || "";

    if (/https:\/\/www\.gravatar\.com\/avatar\/[a-z0-9]{32}(\?.+)?/.test(loadedProfilePictureUrl)) {
      profilePictureType = "gravatar";
    } else if (/\/api\/files\/([a-f0-9-]{36})/.test(loadedProfilePictureUrl)) {
      profilePictureType = "database";
    } else {
      profilePictureType = "remote";
      profilePictureRemoteUrl = loadedProfilePictureUrl;
    }
  });

  let gravatarCheckDefaultTimeout = $state<ReturnType<typeof setTimeout>>();
  $effect(() => {
    if (profilePictureType !== "gravatar") return;
    clearTimeout(gravatarCheckDefaultTimeout);

    setTimeout(async (url: string) => {
      await fetch(url + "&d=404").then((response) => {
        profilePictureGravatarIsDefault = response.status === 404;
      }).catch(() => {
        profilePictureGravatarIsDefault = true;
      });
      if (profilePictureGravatarIsDefault) profilePictureGravatarForceDefault = false;
    }, 100, profilePictureGravatarUrlTrue);
  })

  // Changes

  let userDataSnapshot = $state<UserData | null>(null);
  let userSettingsSnapshot = $state<UserSettings | null>(null);
  let globalSettingsSnapshot = $state<GlobalSettings | null>(null);

  let userDataChanged = $derived(
    !deepEquality(settings.userData, userDataSnapshot) ||
    newPassword != "" ||
    effectiveProfilePictureSource != userDataSnapshot?.profile_picture
  );
  let userSettingsChanged = $derived(!deepEquality(settings.userSettings, userSettingsSnapshot));
  let globalSettingsChanged = $derived(!deepEquality(settings.globalSettings, globalSettingsSnapshot));
  let anyChanged = $derived(userDataChanged || userSettingsChanged || globalSettingsChanged);

  let usernameValidity = $state(valid);
  let emailValidity = $state(valid);
  let passwordValidity = $state(valid);
  let repeatPasswordValidity = $state(valid);

  let requirePasswordForAccountDeletion = $state(false);
  let oldPasswordRequired = $derived(userDataSnapshot && (
    userDataSnapshot.username != settings.userData.username ||
    userDataSnapshot.email != settings.userData.email ||
    newPassword != "" ||
    requirePasswordForAccountDeletion
  ));

  let submittable = $derived(
    !userDataSnapshot || !userDataChanged || (
      (userDataSnapshot.username == settings.userData.username || usernameValidity.valid) &&
      (userDataSnapshot.email == settings.userData.email || emailValidity.valid) &&
      (newPassword == "" || passwordValidity.valid) &&
      (newPassword == "" || repeatPasswordValidity.valid)
    )
  )

  // Interaction with the shared data structures

  async function snapshotSettings() {
    userDataSnapshot = await deepCopy(settings.userData);
    userSettingsSnapshot = await deepCopy(settings.userSettings);
    globalSettingsSnapshot = await deepCopy(settings.globalSettings);
  }

  async function restoreSettings() {
    newPassword = "";

    if (userDataSnapshot) settings.userData = await deepCopy(userDataSnapshot);
    if (userSettingsSnapshot) settings.userSettings = await deepCopy(userSettingsSnapshot);
    if (globalSettingsSnapshot) settings.globalSettings = await deepCopy(globalSettingsSnapshot);

    refetchProfilePicture();
  }

  let saving = $state(false);
  async function saveSettings() {
    if (saving) return;
    saving = true;

    if (userDataChanged && userDataSnapshot) {
      const userDataFormData = new FormData();

      let oldPassword = "";
      if (oldPasswordRequired) {
        oldPassword = await passwordPrompt().catch(() => "");
        if (oldPassword == "") return;
      }

      if (settings.userData.username != userDataSnapshot.username)
        userDataFormData.append("username", settings.userData.username);
      if (settings.userData.email != userDataSnapshot.email)
        userDataFormData.append("email", settings.userData.email);
      if (newPassword != "")
        userDataFormData.append("new_password", newPassword);
      if (oldPasswordRequired)
        userDataFormData.append("password", oldPassword);
      if (settings.userData.searchable != userDataSnapshot.searchable)
        userDataFormData.append("searchable", settings.userData.searchable ? "true" : "false");
      if (profilePictureType !== "database" && effectiveProfilePictureSource != userDataSnapshot.profile_picture)
        userDataFormData.append("pfp_url", effectiveProfilePictureSource);
      if (profilePictureType === "database" && profilePictureFiles && profilePictureFileId === "")
        userDataFormData.append("pfp_file", profilePictureFiles[0]);

      await fetchJson("/api/users/self", {
        method: "PATCH",
        body: userDataFormData,
      }).then(async (response) => {
        if ("profile_picture" in response) {
          settings.userData.profile_picture = response.profile_picture;
          profilePictureFiles = null;
        }
        userDataSnapshot = await deepCopy(settings.userData);
        refetchProfilePicture();
      }).catch((err) => {
        queueNotification(ColorKeys.Danger, "Failed to save user data: " + err);
        saving = false;
        throw err;
      });
    }

    if (userSettingsChanged && userSettingsSnapshot) {
      const userSettingsFormData = new FormData();

      for (const [_, key] of Object.entries(UserSettingKeys)) {
        const originalValue = userSettingsSnapshot[key];
        const newValue = settings.userSettings[key];
        if (originalValue !== newValue)
          userSettingsFormData.append(key, JSON.stringify(newValue));
      }

      await fetchResponse("/api/users/self/settings", {
        method: "PATCH",
        body: userSettingsFormData,
      }).then(async () => {
        userSettingsSnapshot = await deepCopy(settings.userSettings);
      }).catch((err) => {
        queueNotification(ColorKeys.Danger, "Failed to save user settings: " + err);
        saving = false;
        throw err;
      });
    }

    if (globalSettingsChanged && globalSettingsSnapshot) {
      const globalSettingsFormData = new FormData();

      for (const [_, key] of Object.entries(GlobalSettingKeys)) {
        const originalValue = globalSettingsSnapshot[key];
        const newValue = settings.globalSettings[key];
        if (originalValue !== newValue)
          globalSettingsFormData.append(key, JSON.stringify(newValue));
      }

      await fetchResponse("/api/settings", {
        method: "PATCH",
        body: globalSettingsFormData,
      }).then(async () => {
        globalSettingsSnapshot = await deepCopy(settings.globalSettings);
      }).catch((err) => {
        queueNotification(ColorKeys.Danger, "Failed to save global settings: " + err);
        saving = false;
        throw err;
      });
    }

    settings.saveSettings();
    saving = false;
  }

  // Session management actions
  function logout() {
    showConfirmation("Are you sure you want to log out?", async () => {
      await fetchResponse("/api/sessions/current", { method: "DELETE" }); // We don't need to check for errors, because the cookie is deleted either way
      clearSession();
    });
  }
  function deauthorizeSessions() {
    showConfirmation("Are you sure you want to deauthorize all sessions?\nThis will log you out of all your devices.", async () => {
      sessions.deauthorizeUserSessions().catch((err) => {
        queueNotification(ColorKeys.Danger, err);
      });
    }, "Your API tokens will remain valid.\nTo deauthorize those, head to the \"Developer\" tab.");
  }
  function deauthorizeSession(id: string) {
    if (id === sessions.currentSession) return logout();
    sessions.deauthorizeSession(id);
  }
  let editApiToken = $state<(session: Session, editable: boolean) => Promise<Session>>(Promise.reject);
  let createApiToken = $state<() => Promise<Session>>(Promise.reject);

  // Danger zone actions
  function deleteAccount() {
    showConfirmation("Are you sure you want to delete your account?\nThis action is irreversible.", async () => {
      requirePasswordForAccountDeletion = true;
      const password = await new Promise<string>(resolve => setTimeout(async () => resolve(await passwordPrompt().catch(() => "")), 0));
      requirePasswordForAccountDeletion = false;
      if (password == "") return;
      const body = new FormData();
      body.append("password", password);
      await fetchResponse("/api/users/self", { method: "DELETE", body: body });
      clearSession();
    }, "All your data will be deleted.");
  }
  function resetPreferences() {
    showConfirmation("Are you sure you want to reset all your preferences?\nThis action is irreversible.", async () => {
      await fetchResponse("/api/users/self/settings", { method: "DELETE" });
      settings.fetchSettings().then(() => {
        snapshotSettings();
        refetchProfilePicture();
      });
    }, "Your account will remain intact.");
  }
  function resetGlobalSettings() {
    showConfirmation("Are you sure you want to reset all global settings?\nThis action is irreversible.", async () => {
      await fetchResponse("/api/settings", { method: "DELETE" });
      settings.fetchSettings().then(() => {
        snapshotSettings();
        refetchProfilePicture();
      });
    });
  }

  // Appearance
  $effect(() => {
    settings.userSettings[UserSettingKeys.ThemeSynchronize], setTimeout(() => getTheme().refetchTheme(), 0);
  });

  // User Management
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

  // Confirmation dialog
  let internalShowConfirmation = $state(NoOp);
  let confirmationCallback = $state(async () => {});
  let confirmationMessage = $state("");
  let confirmationDetails = $state("");
  function showConfirmation(message: string, callback: () => Promise<void>, details: string = "") {
    confirmationMessage = message;
    confirmationDetails = details;
    confirmationCallback = callback;
    internalShowConfirmation();
  }
</script>

<style lang="scss">
  @use "../../styles/animations.scss";
  @use "../../styles/colors.scss";
  @use "../../styles/dimensions.scss";
  @use "../../styles/text.scss";

  div.container {
    box-sizing: border-box;
    display: grid;
    grid-template-columns: auto 1fr;
    grid-template-rows: 1fr;
    gap: dimensions.$gapMiddle;
    min-width: 30vw;
    height: 60vh;
  }

  main {
    display: flex;
    flex-direction: column;
    flex-wrap: nowrap;
    gap: dimensions.$gapMiddle;
    overflow-y: auto;
    overflow-x: hidden;
    padding-right: dimensions.$gapLarger;
    margin-right: -(dimensions.$gapLarger);
  }

  main > :global(*) {
    flex-shrink: 0;
  }

  div.pfpButtons {
    display: flex;
    flex-direction: column;
    gap: dimensions.$gapMiddle;
    width: 100%;
  }

  .confirmation {
    white-space: pre-wrap;
  }

  .session {
    padding: dimensions.$gapMiddle;
    background-color: colors.$backgroundSecondary;
    color: colors.$foregroundSecondary;
    border-radius: dimensions.$borderRadius;

    display: grid;
    gap: dimensions.$gapSmall;
    row-gap: 0;
    grid-template-columns: auto 1fr auto;
    grid-template-rows: auto auto;
    grid-template-areas: "device agent buttons" "device details buttons";
    justify-content: center;
    align-items: center;
  }

  .session.showId {
    grid-template-rows: auto auto auto;
    grid-template-areas: "device agent buttons" "device details buttons" "device id buttons";
  }

  .session > .device {
    grid-area: device;
    display: flex;
  }
  .session > .agent {
    grid-area: agent;
  }
  .session > .details {
    grid-area: details;
    font-size: text.$fontSizeSmall;
  }
  .session > .buttons {
    grid-area: buttons;
    display: flex;
    flex-direction: row;
    gap: dimensions.$gapSmall;
  }
  .session > .id {
    grid-area: id;
    font-size: text.$fontSizeSmall;
  }
  .session:not(.showId) > .id {
    display: none;
  }

  .session.active {
    background-color: colors.$backgroundAccent;
    color: colors.$foregroundAccent;
  }

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
  
  span.refreshButtonWrapper {
    display: flex;
    align-items: center;
    justify-content: center;
  }

  span.spin {
    animation: spin animations.$animationSpeedSlow animations.$cubic infinite forwards;
  }
</style>

<Modal
  title={"Settings"}
  bind:showModal={showModalInternal}
  bind:hideModal={hideModalInternal}
  onModalHide={restoreSettings}
>
  {#snippet topButtons()}
    <IconButton click={forceRefresh}>
      <span class="refreshButtonWrapper" class:spin={loaderAnimation} onanimationiteration={() => { loaderAnimation = false; }}>
        <RefreshCw/>
      </span>
    </IconButton>
  {/snippet}

  {#snippet buttons()}
    {#if anyChanged}
      <Button type="submit" color={ColorKeys.Success} enabled={submittable} onClick={() => { saveSettings().catch(NoOp); }}>
        {#if saving}
          <Loader/>
        {:else}
          Save
        {/if}
      </Button>
      <Button color={ColorKeys.Danger} onClick={restoreSettings}>Cancel</Button>
    {:else}
      <Button onClick={hideModalInternal}>Close</Button>
    {/if}
  {/snippet}

  <div class="container">
    <ButtonList
      bind:value={selectedCategory}
      options={settings.userData.admin ? categoriesAdmin : categories} 
    />
    <main tabindex="-1">
      {#if selectedCategory === "account"}
        <TextInput
          name="username"
          placeholder="Username"
          bind:value={settings.userData.username}
          validation={isValidUsername}
          bind:validity={usernameValidity}
        />
        <TextInput
          name="email"
          placeholder="Email"
          bind:value={settings.userData.email}
          validation={isValidEmail}
          bind:validity={emailValidity}
        />
        <TextInput
          name="new_password"
          placeholder="New Password"
          password={true}
          bind:value={newPassword}
          validation={isValidPassword}
          bind:validity={passwordValidity}
        />
        {#if newPassword != "" && passwordValidity.valid}
          <TextInput
            name="new_password_confirm"
            placeholder="Confirm New Password"
            password={true}
            validation={isValidRepeatPassword(newPassword)}
            bind:validity={repeatPasswordValidity}
          />
        {/if}
        <ToggleInput
          name="searchable" 
          description="Allow Other Users To Find Me"
          bind:value={settings.userData.searchable}
        />
        <Horizontal position="justify" width="full">
          <div class="pfpButtons">
            <SelectButtons
              name="pfp_type"
              placeholder="Profile Picture"
              bind:value={profilePictureType}
              options={[
                { name: "Gravatar", value: "gravatar" },
                { name: "Upload File", value: "database" },
                { name: "Internet Link", value: "remote" }
              ]}
            />
          </div>
          <Image
            src={effectiveProfilePictureSource}
            alt="Profile Picture"
          />
        </Horizontal>
        {#if profilePictureType === "database"}
          <FileUpload
            name="pfp_file"
            placeholder="Profile Picture File"
            bind:files={profilePictureFiles}
            bind:fileId={profilePictureFileId}
            accept={"image/*"}
          />
          {#if profilePictureFileId != "" && profilePictureFiles && settings.userSettings[UserSettingKeys.DebugMode]}
            <TextInput value={profilePictureFileId} name="id" placeholder="File ID" editable={false} />
          {/if}
        {:else if profilePictureType === "remote"}
          <TextInput
            name="pfp_link"
            placeholder="Profile Picture Link"
            bind:value={profilePictureRemoteUrl}
          />
        {:else if profilePictureType === "gravatar"}
          {#if !profilePictureGravatarIsDefault}
            <ToggleInput
              name="pfp_gravatar_force_default"
              description="Use Default Gravatar Profile Picture"
              bind:value={profilePictureGravatarForceDefault}
            />
          {/if}
        {/if}
        {#if settings.userData.id && settings.userSettings[UserSettingKeys.DebugMode]}
          <TextInput value={settings.userData.id} name="id" placeholder="User ID" editable={false} />
        {/if}
      {:else if selectedCategory === "appearance"}
        <SectionDivider title={"Calendar Appearance"}/>
        <ToggleInput
          name={UserSettingKeys.DisplayAllDayEventsFilled}
          description="Fill All-Day Events"
          bind:value={settings.userSettings[UserSettingKeys.DisplayAllDayEventsFilled]}
        />
        <ToggleInput
          name={UserSettingKeys.DisplayNonAllDayEventsFilled}
          description="Fill Non-All-Day Events"
          bind:value={settings.userSettings[UserSettingKeys.DisplayNonAllDayEventsFilled]}
        />
        <ToggleInput
          name={UserSettingKeys.DisplaySmallCalendar}
          description="Display Small Calendar"
          bind:value={settings.userSettings[UserSettingKeys.DisplaySmallCalendar]}
        />
        <ToggleInput
          name={UserSettingKeys.DynamicCalendarRows}
          description="Dynamic Calendar Row Count"
          bind:value={settings.userSettings[UserSettingKeys.DynamicCalendarRows]}
        />
        {#if settings.userSettings[UserSettingKeys.DisplaySmallCalendar]}
          <ToggleInput
            name={UserSettingKeys.DynamicSmallCalendarRows}
            description="Dynamic Small Calendar Row Count"
            bind:value={settings.userSettings[UserSettingKeys.DynamicSmallCalendarRows]}
          />
        {/if}
        <ToggleInput
          name={UserSettingKeys.DisplayWeekNumbers}
          description="Display Week Numbers"
          bind:value={settings.userSettings[UserSettingKeys.DisplayWeekNumbers]}
        />
        <SelectInput
          name={UserSettingKeys.FirstDayOfWeek}
          placeholder="First Day of Week"
          bind:value={settings.userSettings[UserSettingKeys.FirstDayOfWeek]}
          options={[
            { name: "Monday", value: 1 },
            { name: "Tuesday", value: 2 },
            { name: "Wednesday", value: 3 },
            { name: "Thursday", value: 4 },
            { name: "Friday", value: 5 },
            { name: "Saturday", value: 6 },
            { name: "Sunday", value: 0 }
          ]}
        />
        <SectionDivider title={"Site Appearance"}/>
        <ToggleInput
          name={UserSettingKeys.AppearenceFrostedGlass}
          description="Frosted Glass Effect"
          bind:value={settings.userSettings[UserSettingKeys.AppearenceFrostedGlass]}
        />
        <ToggleInput
          name={UserSettingKeys.DisplayRoundedCorners}
          description="Rounded Corners"
          bind:value={settings.userSettings[UserSettingKeys.DisplayRoundedCorners]}
        />
        <SliderInput
          name={UserSettingKeys.UiScaling}
          title="Scaling"
          bind:value={settings.userSettings[UserSettingKeys.UiScaling]}
          min={0.5}
          max={1.5}
          step={0.05}
          detentTransform={(value) => `${Math.round(value * 100)}%`}
        />
        <SelectInput
          name={UserSettingKeys.ThemeLight}
          placeholder="Light Theme"
          bind:value={settings.userSettings[UserSettingKeys.ThemeLight]}
          options={lightThemes}
        />
        <SelectInput
          name={UserSettingKeys.ThemeDark}
          placeholder="Dark Theme"
          bind:value={settings.userSettings[UserSettingKeys.ThemeDark]}
          options={darkThemes}
        />
        <ToggleInput
          name={UserSettingKeys.ThemeSynchronize}
          description="Synchronize Theme with System"
          bind:value={settings.userSettings[UserSettingKeys.ThemeSynchronize]}
        />
        <SelectInput
          name={UserSettingKeys.FontText}
          placeholder="Text Font"
          bind:value={settings.userSettings[UserSettingKeys.FontText]}
          options={fonts}
        />
        <SelectInput
          name={UserSettingKeys.FontTime}
          placeholder="Monospaced Font"
          bind:value={settings.userSettings[UserSettingKeys.FontTime]}
          options={fonts}
        />
        <SectionDivider title="Animations"/>
        <SliderInput
          name={UserSettingKeys.AnimationDuration}
          title="Animation Duration"
          info={"To disable animations, set the animation duration to 0%."}
          bind:value={settings.userSettings[UserSettingKeys.AnimationDuration]}
          min={0}
          max={2}
          step={0.1}
          detentTransform={(value) => `${Math.round(value * 100)}%`}
        />
        {#if settings.userSettings[UserSettingKeys.AnimationDuration] != 0}
          <ToggleInput
            name={UserSettingKeys.AnimateCalendarSwipe}
            description="Animate Calendar"
            bind:value={settings.userSettings[UserSettingKeys.AnimateCalendarSwipe]}
          />
          {#if settings.userSettings[UserSettingKeys.DisplaySmallCalendar]}
            <ToggleInput
              name={UserSettingKeys.AnimateSmallCalendarSwipe}
              description="Animate Small Calendar"
              bind:value={settings.userSettings[UserSettingKeys.AnimateSmallCalendarSwipe]}
            />
          {/if}
          <ToggleInput
            name={UserSettingKeys.AnimateMonthSelectionSwipe}
            description="Animate Month Selection"
            bind:value={settings.userSettings[UserSettingKeys.AnimateMonthSelectionSwipe]}
          />
        {/if}
      {:else if selectedCategory === "developer"}
        <ToggleInput
          name={UserSettingKeys.DebugMode}
          description="Display IDs"
          bind:value={settings.userSettings[UserSettingKeys.DebugMode]}
        />

        <Button color={ColorKeys.Accent} onClick={() => createApiToken().catch(err => { if (err) queueNotification(ColorKeys.Danger, err.message); } )}>Create an API token</Button>

        {@const apiSessions = sessions.activeSessions.filter(x => x.is_api)}
        {#if apiSessions.length !== 0}
          <List
            label="API Tokens"
            items={apiSessions}
            id={item => item.session_id}
            template={sessionTemplate}
          />
        {/if}
      {:else if selectedCategory === "users"}
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
      {:else if selectedCategory === "admin"}
        <ToggleInput
          name={GlobalSettingKeys.RegistrationEnabled}
          description="Enable Open Registration"
          info={"Allows anyone to create an account.\nIf you just want to invite a few people, head to the \"Users\" tab."}
          bind:value={settings.globalSettings[GlobalSettingKeys.RegistrationEnabled]}
        />
        <!--
        <ToggleInput
          name={GlobalSettingKeys.UseCdnFonts}
          description="Use Google's CDN for fonts"
          bind:value={settings.globalSettings[GlobalSettingKeys.UseCdnFonts]}
        />
        -->
        <SelectButtons
          name={GlobalSettingKeys.LoggingVerbosity}
          bind:value={settings.globalSettings[GlobalSettingKeys.LoggingVerbosity]}
          placeholder="Error Messages Verbosity"
          info={"How much information about errors is returned to the user.\nThis setting applies to all users.\n\"Debug\" should never be used in production."}
          options={[
            { name: "Broad", value: 3 },
            { name: "Plain", value: 2 },
            { name: "Wordy", value: 1 },
            { name: "Debug", value: 0 }
          ]}
        />
        <Button color={ColorKeys.Danger} onClick={resetGlobalSettings}>Reset all global settings</Button>
      {:else if selectedCategory === "danger"}
        <Button color={ColorKeys.Danger} onClick={resetPreferences}>Reset all my preferences</Button>
        <Button color={ColorKeys.Danger} onClick={deleteAccount}>Delete my account</Button>
      {:else if selectedCategory === "logout"}
        <Button color={ColorKeys.Danger} onClick={logout}>Log out of my account</Button>
        <Button color={ColorKeys.Danger} onClick={deauthorizeSessions}>Deauthorize all sessions</Button>

        <List
          label="Active Sessions"
          info={"To see your API sessions, head to the \"Developer\" tab."}
          items={sessions.activeSessions.filter(x => !x.is_api)}
          id={item => item.session_id}
          template={sessionTemplate}
        />
      {/if}
    </main>
  </div>
</Modal>

{#snippet sessionTemplate(s: Session)}
  {@const userAgent=UAParser(s.is_api ? "" : s.user_agent)}
  {@const deviceName=`${userAgent.os.name || ""} ${userAgent.browser.name || ""}`.trim()}
  {@const isActive=s.session_id === sessions.currentSession}

  <div class="session" class:active={isActive} class:showId={settings.userSettings[UserSettingKeys.DebugMode]}>
    <div class="device">
      {#if s.is_api}
        <Bot size={20}/>
      {:else if userAgent.device.type === UAParser.DEVICE.CONSOLE}
        <Gamepad2 size={20}/>
      {:else if userAgent.device.type === UAParser.DEVICE.EMBEDDED}
        <Microchip size={20}/>
      {:else if userAgent.device.type === UAParser.DEVICE.MOBILE}
        <Smartphone size={20}/>
      {:else if userAgent.device.type === UAParser.DEVICE.SMARTTV}
        <TvMinimal size={20}/>
      {:else if userAgent.device.type === UAParser.DEVICE.TABLET}
        <Tablet size={20}/>
      {:else if userAgent.device.type === UAParser.DEVICE.WEARABLE}
        <Watch size={20}/>
      <!--{:else if userAgent.device.type === UAParser.DEVICE.XR}-->
      {:else if userAgent.device.type === "xr"}
        <RectangleGoggles size={20}/>
      {:else if deviceName === ""}
        <Bot size={20}/>
      {:else}
        <Laptop size={20}/>
      {/if}
    </div>

    <span class="agent">
      {#if deviceName === ""}
        {s.user_agent}
      {:else}
        {deviceName}
      {/if}
    </span>

    <span class="details">
      {s.location}
      •
      {#if isActive}
        Current session
      {:else if today.getDate() == s.last_seen.getDate() && today.getMonth() == s.last_seen.getMonth() && today.getFullYear() == s.last_seen.getFullYear()}
        Last active {s.last_seen.toLocaleTimeString()}
      {:else}
        Last active {s.last_seen.toLocaleDateString()} {s.last_seen.toLocaleTimeString()}
      {/if}
    </span>

    <div class="buttons">
      <IconButton click={() => editApiToken(s, s.is_api)}>
        {#if s.is_api}
          <Pencil size={20}/>
        {:else}
          <Info size={20}/>
        {/if}
      </IconButton>
      <IconButton click={() => deauthorizeSession(s.session_id)}>
        <LogOut size={20}/>
      </IconButton>
    </div>

    <span class="id">
      ID: {s.session_id}
    </span>
  </div>
{/snippet}

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
      <!-- TODO: transform invite.author into username -->
      Created {invite.created_at.toLocaleDateString()} at {invite.created_at.toLocaleTimeString()} by {"TODO"}
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

<RegistrationInviteModal
  bind:showModal={showRegistrationInvite}
  bind:showCreateModal={issueRegistrationInvite}
/>

<SessionModal
  bind:showModal={editApiToken}
  bind:showCreateModal={createApiToken}
/>

{#if submittable && oldPasswordRequired}
  <PasswordPromptModal bind:prompt={passwordPrompt}/>
{/if}

<ConfirmationModal
  bind:showModal={internalShowConfirmation}
  confirmCallback={confirmationCallback}
>
  <span class="confirmation">
    {confirmationMessage}
    {#if confirmationDetails != ""}
      <Tooltip inline>{confirmationDetails}</Tooltip>
    {/if}
  </span>
</ConfirmationModal>
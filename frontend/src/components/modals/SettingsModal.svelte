<script lang="ts">
  import { CaseSensitive, Code, LogOut, Moon, Palette, RefreshCw, Shield, Sun, TriangleAlert, User, Users } from "lucide-svelte";
  import { AsyncNoOp, NoOp } from "../../lib/client/placeholders";
  import ButtonList from "../forms/ButtonList.svelte";
  import Modal from "./Modal.svelte";
  import { GlobalSettingKeys, UserSettingKeys, type GlobalSettings, type UserData, type UserSettings } from "../../types/settings";
  import { getSettings } from "../../lib/client/data/settings.svelte";
  import { deepCopy, deepEquality } from "$lib/common/misc";
  import Button from "../interactive/Button.svelte";
  import Loader from "../decoration/Loader.svelte";
  import { fetchJson, fetchResponse } from "$lib/client/net";
  import { queueNotification } from "$lib/client/notifications";
  import { ColorKeys } from "../../types/colors";
  import type { Option } from "../../types/options";
  import ConfirmationModal from "./ConfirmationModal.svelte";
  import { clearSession, getActiveSessions } from "../../lib/client/data/sessions.svelte";
  import Tooltip from "../interactive/Tooltip.svelte";
  import IconButton from "../interactive/IconButton.svelte";
  import PasswordPromptModal from "./PasswordPromptModal.svelte";
  import { getTheme } from "../../lib/client/data/theme.svelte";
  import { getRegistrationInvites } from "../../lib/client/data/invites.svelte";
  import { getUsers } from "$lib/client/data/users.svelte";
  import AccountSettingsTab from "./settingModalTabs/AccountSettingsTab.svelte";
  import AppearanceSettingsTab from "./settingModalTabs/AppearanceSettingsTab.svelte";
  import DeveloperSettingsTab from "./settingModalTabs/DeveloperSettingsTab.svelte";
  import UsersSettingsTab from "./settingModalTabs/UsersSettingsTab.svelte";
  import AdminSettingsTab from "./settingModalTabs/AdminSettingsTab.svelte";
  import DangerSettingsTab from "./settingModalTabs/DangerSettingsTab.svelte";
  import LogoutSettingsTab from "./settingModalTabs/LogoutSettingsTab.svelte";
  import ThemesSettingsTab from "./settingModalTabs/ThemesSettingsTab.svelte";
  import FontsSettingsTab from "./settingModalTabs/FontsSettingsTab.svelte";
  import SessionModal from "./SessionModal.svelte";
  import { getDatabaseFileIdFromUrl } from "../../lib/common/parsing";

  interface Props {
    showModal?: () => any;
    hideModal?: () => any;
  }

  let {
    showModal = $bindable(),
    hideModal = $bindable(NoOp),
  }: Props = $props();

  // Global data structures
  const settings = getSettings();
  const sessions = getActiveSessions();
  const invites = getRegistrationInvites();
  const users = getUsers();

  const today = new Date();

  // Functions and props exported by individual tabs
  let accountSettingsSubmittable = $state(false);
  let accountSettingsReauthenticationRequired = $state(false);
  let accountSettingsNewPassword = $state("");
  let accountSettingsNewProfilePictureChosen = $state(false);
  let accountSettingsNewProfilePictureUrl = $state("");
  let accountSettingsNewProfilePictureFile: File | null = $state(null);
  let refetchProfilePicture = $state<() => void>(NoOp);

  let dangerSettingsReauthenticationRequired = $state(false);

  // Loading
  let loaderAnimation = $state(false);
  function forceRefresh() {
    loaderAnimation = true;
    fetchThemes();
    fetchFonts();
    sessions.fetch();
    settings.fetchSettings().then(() => {
      snapshotSettings();
      refetchProfilePicture();

      if (settings.userData.admin) {
        invites.fetch();
        users.fetchAll();
      }
    });
  }

  // Show and hide hooks
  showModal = () => {
    saving = false;
    snapshotSettings();
    forceRefresh();
    showModalInternal();
  };

  hideModal = () => {
    hideModalInternal();
  };

  let showModalInternal = $state(NoOp);
  let hideModalInternal = $state(NoOp);

  // Settings categories
  const categoriesAdmin: Option<string>[][] = [
    [
      { name: "Account", value: "account", icon: User },
      { name: "Appearance", value: "appearance", icon: Palette },
      { name: "Developer", value: "developer", icon: Code }
    ],
    [
      { name: "Users", value: "users", icon: Users },
      { name: "Themes", value: "themes", icon: Palette },
      { name: "Fonts", value: "fonts", icon: CaseSensitive },
      { name: "Administrative", value: "admin", icon: Shield },
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
  let lightThemes = $state<Option<string>[]>([{ name: "Luna Light", value: "luna-light", icon: Sun }]);
  let darkThemes = $state<Option<string>[]>([{ name: "Luna Dark", value: "luna-dark", icon: Moon }]);
  let fonts = $state<Option<string>[]>([
    { name: "Atkinson Hyperlegible Next", value: "atkinson-hyperlegible-next" },
    { name: "Atkinson Hyperlegible Mono", value: "atkinson-hyperlegible-next" }
  ]);

  function formatInstalledFile(icon: any = null): (rawName: string) => Option<string> {
    return (rawName: string): Option<string> => {
      const formattedName = rawName
        .split("-")
        .map(x => x.charAt(0).toUpperCase() + x.slice(1))
        .join(" ");
      return { name: formattedName, value: rawName, icon: icon };
    }
  }

  function fetchThemes() {
    fetchJson("/installed/themes").then((response) => {
      lightThemes = Object.keys(response.light).map(formatInstalledFile(Sun)).sort((a, b) => a.name.localeCompare(b.name));
      darkThemes = Object.keys(response.dark).map(formatInstalledFile(Moon)).sort((a, b) => a.name.localeCompare(b.name));
    }).catch((err) => {
      queueNotification(ColorKeys.Danger, "Failed to fetch themes: " + err);
    });
  }

  function fetchFonts() {
    fetchJson("/installed/fonts").then((response) => {
      fonts = Object.keys(response).map(formatInstalledFile()).sort((a, b) => a.name.localeCompare(b.name));
    }).catch((err) => {
      queueNotification(ColorKeys.Danger, "Failed to fetch fonts: " + err);
    });
  }

  // Static data change tracking
  let userDataSnapshot = $state<UserData | null>(null);
  let userSettingsSnapshot = $state<UserSettings | null>(null);
  let globalSettingsSnapshot = $state<GlobalSettings | null>(null);

  let userDataChanged = $derived(!deepEquality(settings.userData, userDataSnapshot));
  let userSettingsChanged = $derived(!deepEquality(settings.userSettings, userSettingsSnapshot));
  let globalSettingsChanged = $derived(!deepEquality(settings.globalSettings, globalSettingsSnapshot));
  let anyChanged = $derived(userDataChanged || userSettingsChanged || globalSettingsChanged);

  // Snapshot and restore functions
  async function snapshotSettings() {
    userDataSnapshot = await deepCopy(settings.userData);
    userSettingsSnapshot = await deepCopy(settings.userSettings);
    globalSettingsSnapshot = await deepCopy(settings.globalSettings);
  }

  async function restoreSettings() {
    accountSettingsNewPassword = "";

    if (userDataSnapshot) settings.userData = await deepCopy(userDataSnapshot);
    if (userSettingsSnapshot) settings.userSettings = await deepCopy(userSettingsSnapshot);
    if (globalSettingsSnapshot) settings.globalSettings = await deepCopy(globalSettingsSnapshot);

    refetchProfilePicture();
  }

  // Submittability determination
  let submittable = $derived.by(() => {
    switch (selectedCategory) {
      case "account":
        return accountSettingsSubmittable && (userDataChanged || userSettingsChanged || accountSettingsNewPassword !== "" || accountSettingsNewProfilePictureChosen);
      case "danger":
        return true;
      default:
        return anyChanged;
    }
  })

  let showSaveButton = $derived.by(() => {
    switch (selectedCategory) {
      case "danger":
        return false;
      default:
        return submittable;
    }
  })

  let reauthenticationRequired = $derived.by(() => {
    switch (selectedCategory) {
      case "account":
        return accountSettingsReauthenticationRequired;
      case "danger":
        return dangerSettingsReauthenticationRequired;
      default:
        return false;
    }
  })

  // Account Settings
  let passwordPrompt = $state<() => Promise<string>>(() => Promise.reject(""));
    
  function deleteAccount(id: string = "self") {
    const ownAccount = id === "self" || id === users.currentUser;
    showConfirmation(`Are you sure you want to delete ${ownAccount ? "your" : "this"} account?\nThis action is irreversible.`, async () => {
      dangerSettingsReauthenticationRequired = true;
      const password = await new Promise<string>(resolve => setTimeout(async () => resolve(await passwordPrompt().catch(() => "")), 0));
      dangerSettingsReauthenticationRequired = false;
      if (password == "") return;
      const body = new FormData();
      body.append("password", password);
      await fetchResponse(`/api/users/${id}`, { method: "DELETE", body: body });
      if (ownAccount) clearSession();
    }, `All ${ownAccount ? "your" : "user"} data will be deleted.`);
  }

  let saving = $state(false);
  async function saveSettings() {
    if (saving) return;
    saving = true;

    if (userDataSnapshot && (userDataChanged || accountSettingsNewPassword !== "" || accountSettingsNewProfilePictureChosen)) {
      const userDataFormData = new FormData();

      let oldPassword = "";
      if (reauthenticationRequired) {
        oldPassword = await passwordPrompt().catch(() => "");
        if (oldPassword == "") return;
      }

      if (settings.userData.username != userDataSnapshot.username)
        userDataFormData.append("username", settings.userData.username);
      if (settings.userData.email != userDataSnapshot.email)
        userDataFormData.append("email", settings.userData.email);
      if (accountSettingsNewPassword != "")
        userDataFormData.append("new_password", accountSettingsNewPassword);
      if (reauthenticationRequired)
        userDataFormData.append("password", oldPassword);
      if (settings.userData.searchable != userDataSnapshot.searchable)
        userDataFormData.append("searchable", settings.userData.searchable ? "true" : "false");
      if (accountSettingsNewProfilePictureChosen) {
        userDataFormData.append("pfp_type", settings.userData.profile_picture_type);
        switch (settings.userData.profile_picture_type) {
          case "gravatar":
          case "static":
          case "remote":
            userDataFormData.append("pfp_url", accountSettingsNewProfilePictureUrl);
            break;
          case "database":
            if (accountSettingsNewProfilePictureFile)
              userDataFormData.append("pfp_file", accountSettingsNewProfilePictureFile);
            break;
        }
      }

      await fetchJson("/api/users/self", {
        method: "PATCH",
        body: userDataFormData,
      }).then(async (response) => {
        saving = false;
        if ("profile_picture" in response) {
          settings.userData.profile_picture = response.profile_picture;
          settings.userData.profile_picture_url = accountSettingsNewProfilePictureUrl;
          if (settings.userData.profile_picture_type === "database") {
            settings.userData.profile_picture_file = getDatabaseFileIdFromUrl(response.profile_picture) || "";
          } else {
            settings.userData.profile_picture_file = "";
          }
          accountSettingsNewProfilePictureFile = null;
        }
        userDataSnapshot = await deepCopy(settings.userData);
        refetchProfilePicture();
        if (settings.userData.admin) users.fetchAll();
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

  // Theme synchronization
  $effect(() => {
    settings.userSettings[UserSettingKeys.ThemeSynchronize], setTimeout(() => getTheme().refetchTheme(), 0);
  });

  // Dialogs
  let editApiToken = $state<(session: Session, editable: boolean) => Promise<Session>>(Promise.reject);
  let createApiToken = $state<() => Promise<Session>>(Promise.reject);

  let internalShowConfirmation = $state(NoOp);
  let confirmationCallback = $state(AsyncNoOp);
  let cancellationCallback = $state(AsyncNoOp);
  let confirmationMessage = $state("");
  let confirmationDetails = $state("");
  function showConfirmation(message: string, callback: () => Promise<void>, details: string = "", cancel: () => Promise<void> = AsyncNoOp) {
    confirmationMessage = message;
    confirmationDetails = details;
    confirmationCallback = callback;
    cancellationCallback = cancel;
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

  .confirmation {
    white-space: pre-wrap;
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
    {#if showSaveButton}
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
        <AccountSettingsTab
          settings={settings} 
          userDataSnapshot={userDataSnapshot}
          bind:accountSettingsSubmittable={accountSettingsSubmittable}
          bind:accountSettingsReauthenticationRequired={accountSettingsReauthenticationRequired}
          bind:accountSettingsNewProfilePictureChosen={accountSettingsNewProfilePictureChosen}
          bind:accountSettingsNewProfilePictureUrl={accountSettingsNewProfilePictureUrl}
          bind:accountSettingsNewProfilePictureFile={accountSettingsNewProfilePictureFile}
          bind:newPassword={accountSettingsNewPassword}
          bind:refetchProfilePicture={refetchProfilePicture}
        />
      {:else if selectedCategory === "appearance"}
        <AppearanceSettingsTab
          settings={settings}
          lightThemes={lightThemes}
          darkThemes={darkThemes}
          fonts={fonts} 
        />
      {:else if selectedCategory === "developer"}
        <DeveloperSettingsTab
          settings={settings} 
          sessions={sessions}
          today={today}
          editApiToken={editApiToken}
          createApiToken={createApiToken}
        />
      {:else if selectedCategory === "users"}
        <UsersSettingsTab
          today={today}
          settings={settings}
          invites={invites}
          users={users}
          showConfirmation={showConfirmation}
          deleteAccount={deleteAccount} 
        />
      {:else if selectedCategory === "admin"}
        <AdminSettingsTab
          settings={settings}
          showConfirmation={showConfirmation}
          snapshotSettings={snapshotSettings}
          refetchProfilePicture={refetchProfilePicture} 
        />
      {:else if selectedCategory === "danger"}
        <DangerSettingsTab
          settings={settings}
          bind:requirePasswordForAccountDeletion={dangerSettingsReauthenticationRequired}
          showConfirmation={showConfirmation}
          deleteAccount={deleteAccount}
          refetchProfilePicture={refetchProfilePicture}
          snapshotSettings={snapshotSettings}
        />
      {:else if selectedCategory === "logout"}
        <LogoutSettingsTab
          settings={settings} 
          sessions={sessions}
          today={today}
          editApiToken={editApiToken}
          showConfirmation={showConfirmation} 
        />
      {:else if selectedCategory === "themes"}
        <ThemesSettingsTab
          lightThemes={lightThemes}
          darkThemes={darkThemes}
          fetchThemes={fetchThemes}
          showConfirmation={showConfirmation} 
        />
      {:else if selectedCategory === "fonts"}
        <FontsSettingsTab
          fonts={fonts} 
          fetchFonts={fetchFonts}
          showConfirmation={showConfirmation}
        />
      {/if}
    </main>
  </div>
</Modal>

<SessionModal
  bind:showModal={editApiToken}
  bind:showCreateModal={createApiToken}
/>

{#if submittable && reauthenticationRequired}
  <PasswordPromptModal bind:prompt={passwordPrompt}/>
{/if}

<ConfirmationModal
  bind:showModal={internalShowConfirmation}
  confirmCallback={confirmationCallback}
  cancelCallback={cancellationCallback} 
>
  <span class="confirmation">
    {confirmationMessage}
    {#if confirmationDetails != ""}
      <Tooltip inline>{confirmationDetails}</Tooltip>
    {/if}
  </span>
</ConfirmationModal>
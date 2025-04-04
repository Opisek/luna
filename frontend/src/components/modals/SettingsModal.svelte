<script lang="ts">
  import { Code, LockKeyhole, LogOut, Monitor, TriangleAlert, User, Users } from "lucide-svelte";
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
  import { getSettings } from "../../lib/client/settings.svelte";
  import { getSha256Hash } from "../../lib/common/crypto";
  import { deepCopy, deepEquality } from "$lib/common/misc";
  import Button from "../interactive/Button.svelte";
  import Loader from "../decoration/Loader.svelte";
  import { fetchFileById, fetchJson, fetchResponse } from "$lib/client/net";
  import { queueNotification } from "$lib/client/notifications";
  import { ColorKeys } from "../../types/colors";
  import type { Option } from "../../types/options";

  interface Props {
    showModal?: () => any;
    hideModal?: () => any;
  }

  let {
    showModal = $bindable(),
    hideModal = $bindable(NoOp),
  }: Props = $props();

  const settings = getSettings();

  showModal = () => {
    saving = false;
    snapshotSettings();
    fetchThemes();
    settings.fetchSettings().then(() => {
      snapshotSettings();
      refetchProfilePicture();
    });
    showModalInternal();
  };

  hideModal = () => {
    hideModalInternal();
  };

  function refetchProfilePicture() {
    const profilePictureUrl = settings.userData.profile_picture || "";
    profilePictureFileId = profilePictureUrl.match(/\/api\/files\/([a-f0-9-]{36})/)?.[1] || "";
    if (profilePictureFileId != "") {
      fetchFileById(profilePictureFileId).then(fileList => {
        profilePictureFiles = fileList;
      }).catch(err => {
        queueNotification(ColorKeys.Danger, `Could not download profile picture: ${err.message}`);
      });
    }
  }

  let showModalInternal = $state(NoOp);
  let hideModalInternal = $state(NoOp);

  const categoriesAdmin: Option<string>[][] = [
    [
      { name: "Account", value: "account", icon: User },
      { name: "Appearance", value: "appearance", icon: Monitor },
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
      { name: "Appearance", value: "appearance", icon: Monitor },
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

  function formatThemeOption(theme: string): Option<string> {
    const formattedName = theme
      .split("-")
      .map(x => x.charAt(0).toUpperCase() + x.slice(1))
      .join(" ");

    return { name: formattedName, value: theme };
  }

  function fetchThemes() {
    fetchJson("/installed/themes").then((response) => {
      lightThemes = Object.keys(response.light).map(formatThemeOption);
      darkThemes = Object.keys(response.dark).map(formatThemeOption);
    }).catch((err) => {
      queueNotification(ColorKeys.Danger, "Failed to fetch themes: " + err);
    });
  }

  // Account Settings
  let newPassword = $state("");
  let oldPassword = $state("");

  let profilePictureType = $state("gravatar");
  let profilePictureFiles: FileList | null = $state(null);
  let profilePictureFileId = $state("");
  let profilePictureRemoteUrl = $state("");
  let profilePictureGravatarUrl = $derived.by(() => {
    const email = settings.userData.email || "";
    const trimmedLowercaseEmail = email.trim().toLowerCase();
    const emailHash = getSha256Hash(trimmedLowercaseEmail);
    return `https://www.gravatar.com/avatar/${emailHash}`;
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

    if (/https:\/\/www\.gravatar\.com\/avatar\/[a-z0-9]{32}/.test(loadedProfilePictureUrl)) {
      profilePictureType = "gravatar";
    } else if (/\/api\/files\/([a-f0-9-]{36})/.test(loadedProfilePictureUrl)) {
      profilePictureType = "database";
    } else {
      profilePictureType = "remote";
      profilePictureRemoteUrl = loadedProfilePictureUrl;
    }
  });

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
  let oldPasswordValidity = $state(valid);

  let oldPasswordRequired = $derived(userDataSnapshot && (
    userDataSnapshot.username != settings.userData.username ||
    userDataSnapshot.email != settings.userData.email ||
    newPassword != ""
  ));

  let submittable = $derived(
    !userDataSnapshot || !userDataChanged || (
      (userDataSnapshot.username == settings.userData.username || usernameValidity.valid) &&
      (userDataSnapshot.email == settings.userData.email || emailValidity.valid) &&
      (newPassword == "" || passwordValidity.valid) &&
      (newPassword == "" || repeatPasswordValidity.valid) &&
      (!oldPasswordRequired || oldPasswordValidity.valid) &&
      (!oldPasswordRequired || oldPassword != "")
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
    oldPassword = "";

    if (userDataSnapshot) settings.userData = await deepCopy(userDataSnapshot);
    if (userSettingsSnapshot) settings.userSettings = await deepCopy(userSettingsSnapshot);
    if (globalSettingsSnapshot) settings.globalSettings = await deepCopy(globalSettingsSnapshot);
  }

  let saving = $state(false);
  async function saveSettings() {
    if (saving) return;
    saving = true;

    if (userDataChanged && userDataSnapshot) {
      const userDataFormData = new FormData();

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
        oldPassword = "";
        userDataSnapshot = await deepCopy(settings.userData);
        refetchProfilePicture();
      }).catch((err) => {
        queueNotification(ColorKeys.Danger, "Failed to save user data: " + err);
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
      });
    }

    settings.saveSettings();
    saving = false;
  }
</script>

<style lang="scss">
  @use "../../styles/dimensions.scss";

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
    padding-right: calc(dimensions.$gapLarger);
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
</style>

<Modal
  title={"Settings"}
  bind:showModal={showModalInternal}
  bind:hideModal={hideModalInternal}
  onModalHide={restoreSettings}
>
  {#snippet buttons()}
    {#if anyChanged}
      <Button type="submit" color={ColorKeys.Success} enabled={submittable} onClick={saveSettings}>
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
    <main>
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
        {#if oldPasswordRequired}
          <TextInput
            name="password"
            placeholder="Current Password"
            password={true}
            bind:value={oldPassword}
            validation={isValidPassword}
            bind:validity={oldPasswordValidity}
          />
        {/if}
        <ToggleInput
          name="searchable" 
          description="Allow other users to find me"
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
        {/if}
        {#if settings.userData.id && settings.userSettings[UserSettingKeys.DebugMode]}
          <TextInput value={settings.userData.id} name="id" placeholder="User ID" editable={false} />
        {/if}
      {:else if selectedCategory === "appearance"}
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
        <ToggleInput
          name={UserSettingKeys.DynamicSmallCalendarRows}
          description="Dynamic Small Calendar Row Count"
          bind:value={settings.userSettings[UserSettingKeys.DynamicSmallCalendarRows]}
        />
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
        <ToggleInput
          name={UserSettingKeys.DisplayRoundedCorners}
          description="Rounded Corners"
          bind:value={settings.userSettings[UserSettingKeys.DisplayRoundedCorners]}
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
        <SelectInput
          name={UserSettingKeys.FontText}
          placeholder="Text Font"
          bind:value={settings.userSettings[UserSettingKeys.FontText]}
          options={[
            { name: "Atkinson Hyperlegible Next", value: "Atkinson Hyperlegible Next" },
            { name: "Atkinson Hyperlegible Mono", value: "Atkinson Hyperlegible Mono" }
          ]}
        />
        <SelectInput
          name={UserSettingKeys.FontTime}
          placeholder="Time Font"
          bind:value={settings.userSettings[UserSettingKeys.FontTime]}
          options={[
            { name: "Atkinson Hyperlegible Next", value: "Atkinson Hyperlegible Next" },
            { name: "Atkinson Hyperlegible Mono", value: "Atkinson Hyperlegible Mono" }
          ]}
        />
        TODO: scaling slider
      {:else if selectedCategory === "developer"}
        <ToggleInput
          name={UserSettingKeys.DebugMode}
          description="Display IDs"
          bind:value={settings.userSettings[UserSettingKeys.DebugMode]}
        />
        {#if settings.userData.id && settings.userSettings[UserSettingKeys.DebugMode]}
          <TextInput bind:value={settings.userData.id} name="id" placeholder="User ID" editable={false} />
        {/if}
      {:else if selectedCategory === "users"}
        <Button color={ColorKeys.Accent}>Invite a user</Button>
      {:else if selectedCategory === "admin"}
        <ToggleInput
          name={GlobalSettingKeys.RegistrationEnabled}
          description="Enable Open Registration"
          info={"Allows anyone to create an account.\nIf you just want to invite a few people, head to the \"Users\" tab."}
          bind:value={settings.globalSettings[GlobalSettingKeys.RegistrationEnabled]}
        />
        <ToggleInput
          name={GlobalSettingKeys.UseCdnFonts}
          description="Use Google's CDN for fonts"
          bind:value={settings.globalSettings[GlobalSettingKeys.UseCdnFonts]}
        />
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
      {:else if selectedCategory === "danger"}
        <Button color={ColorKeys.Danger}>Reset all my preferences</Button>
        {#if settings.userData.admin}
          <Button color={ColorKeys.Danger}>Reset all global settings</Button>
        {/if}
        <Button color={ColorKeys.Danger}>Delete my account</Button>
      {:else if selectedCategory === "logout"}
        <Button color={ColorKeys.Danger}>Log out of my account</Button>
        <Button color={ColorKeys.Danger}>Deauthorize all sessions</Button>
      {/if}
    </main>
  </div>
</Modal>
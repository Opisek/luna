<script lang="ts">
  import { Code, LockKeyhole, LogOut, Monitor, User } from "lucide-svelte";
  import { NoOp } from "../../lib/client/placeholders";
  import ButtonList from "../forms/ButtonList.svelte";
  import Modal from "./Modal.svelte";
  import TextInput from "../forms/TextInput.svelte";
  import { isValidEmail, isValidPassword, isValidUsername } from "../../lib/client/validation";
  import ToggleInput from "../forms/ToggleInput.svelte";
  import SelectButtons from "../forms/SelectButtons.svelte";
  import Image from "../layout/Image.svelte";
  import FileUpload from "../forms/FileUpload.svelte";
  import Horizontal from "../layout/Horizontal.svelte";
  import SelectInput from "../forms/SelectInput.svelte";
  import { GlobalSettingKeys, UserSettingKeys, type GlobalSettings, type UserData, type UserSettings } from "../../types/settings";
  import { getSettings } from "../../lib/client/settings.svelte";
  import { getSha256Hash } from "../../lib/common/crypto";
  import { deepCopy } from "$lib/common/misc";

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
    snapshotSettings();
    settings.fetchSettings().then(() => {
      snapshotSettings();
    });
    showModalInternal();
  };

  hideModal = () => {
    hideModalInternal();
  };

  let showModalInternal = $state(NoOp);
  let hideModalInternal = $state(NoOp);

  const categories: Option<string>[][] = [
    [
      { name: "Account", value: "account", icon: User },
      { name: "Appearance", value: "appearance", icon: Monitor },
      { name: "Developer", value: "developer", icon: Code }
    ],
    [
      { name: "Administrative", value: "admin", icon: LockKeyhole },
    ],
    [
      { name: "Logout", value: "logout", icon: LogOut },
    ],
  ]
  let selectedCategory = $state("account");
  let previousCategory = $state("account");
  $effect(() => {
    if (selectedCategory === previousCategory) return;
    previousCategory = selectedCategory;
    restoreSettings();
  });

  // Interaction with the stores
  let userDataSnapshot: UserData | null = $state(null);
  let userSettingsSnapshot: UserSettings | null = $state(null);
  let globalSettingsSnapshot: GlobalSettings | null = $state(null);
  async function snapshotSettings() {
    userDataSnapshot = await deepCopy(settings.userData);
    userSettingsSnapshot = await deepCopy(settings.userSettings);
    globalSettingsSnapshot = await deepCopy(settings.globalSettings);
  }
  async function restoreSettings() {
    if (userDataSnapshot) settings.userData = await deepCopy(userDataSnapshot);
    if (userSettingsSnapshot) settings.userSettings = await deepCopy(userSettingsSnapshot);
    if (globalSettingsSnapshot) settings.globalSettings = await deepCopy(globalSettingsSnapshot);
  }

  // Account Settings
  let profilePictureType = $state("gravatar");
  let profilePictureFiles: FileList | null = $state(null);
  let profilePictureRemoteUrl = $state("");
  let profilePictureGravatarUrl = $derived.by(() => {
    const email = settings.userData.email || "";
    const trimmedLowercaseEmail = email.trim().toLowerCase();
    const emailHash = getSha256Hash(trimmedLowercaseEmail);
    return `https://www.gravatar.com/avatar/${emailHash}`;
  })

  let effectiveProfilePictureSource = $derived.by(() => {
    if (profilePictureType === "gravatar") return profilePictureGravatarUrl;
    else if (profilePictureType === "database") return "TODO"
    else if (profilePictureType === "remote") return profilePictureRemoteUrl;
    else return "";
  });

  $effect(() => {
    const loadedProfilePictureUrl = settings.userData.profile_picture || "";

    if (/https:\/\/www\.gravatar\.com\/avatar\/[a-z0-9]{32}/.test(loadedProfilePictureUrl)) {
      profilePictureType = "gravatar";
    } else if (1 != 1) { // TODO: need current domain to compare
      profilePictureType = "database";
    } else {
      profilePictureType = "remote";
      profilePictureRemoteUrl = loadedProfilePictureUrl;
    }
  });

  // Appearance Settings

  // Developer Settings

  // Admin Settings
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
  <div class="container">
    <ButtonList
      bind:value={selectedCategory}
      options={categories} 
    />
    <main>
      {#if selectedCategory === "account"}
        <TextInput
          name="username"
          placeholder="Username"
          bind:value={settings.userData.username}
          validation={isValidUsername}
        />
        <TextInput
          name="email"
          placeholder="Email"
          bind:value={settings.userData.email}
          validation={isValidEmail}
        />
        <TextInput
          name="password"
          placeholder="New Password"
          password={true}
          validation={isValidPassword}
        />
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
            files={profilePictureFiles}
          />
        {:else if profilePictureType === "remote"}
          <TextInput
            name="pfp_link"
            placeholder="Profile Picture Link"
            bind:value={profilePictureRemoteUrl}
          />
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
          options={[
            { name: "Luna Light", value: "luna-light" },
            { name: "Solarized Light", value: "solarized-light" },
            { name: "Nord Light", value: "nord-light" },
            { name: "High Constrast Light", value: "high-contrast-light" },
          ]}
        />
        <SelectInput
          name={UserSettingKeys.ThemeDark}
          placeholder="Dark Theme"
          bind:value={settings.userSettings[UserSettingKeys.ThemeDark]}
          options={[
            { name: "Luna Dark", value: "luna-dark" },
            { name: "Solarized Dark", value: "solarized-dark" },
            { name: "Nord Dark", value: "Nord Dark" },
            { name: "High Contrast Light", value: "high-contrast-light" },
          ]}
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
      {:else if selectedCategory === "admin"}
        <ToggleInput
          name={GlobalSettingKeys.RegistrationEnabled}
          description="Enable Registration"
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
          options={[
            { name: "Broad", value: 3 },
            { name: "Plain", value: 2 },
            { name: "Wordy", value: 1 },
            { name: "Debug", value: 0 }
          ]}
        />
      {/if}
    </main>
  </div>
</Modal>
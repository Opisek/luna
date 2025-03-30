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
  import { getSettings } from "../../lib/client/setttings";
  import Loader from "../decoration/Loader.svelte";
  import { getSha256Hash } from "../../lib/common/crypto";

  interface Props {
    showModal?: () => any;
    hideModal?: () => any;
  }

  let {
    showModal = $bindable(),
    hideModal = $bindable(NoOp),
  }: Props = $props();

  showModal = () => {
    const settings = getSettings();
    settings.fetchSettings().then(() => {
      userData = settings.getUserData();
      userSettings = settings.getUserSettings();
      globalSettings = settings.getGlobalSettings();
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

  // Setting Data Structures
  let userData: UserData | null = $state(null);
  let userSettings: UserSettings | null = $state(null);
  let globalSettings: GlobalSettings | null = $state(null);

  // Account Settings
  let profilePictureType = $state("gravatar");
  let profilePictureFiles: FileList | null = $state(null);
  let profilePictureRemoteUrl = $state("");
  let profilePictureGravatarUrl = $derived.by(() => {
    if (!userData) return "";
    const email = userData.email || "";
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
    const loadedProfilePictureUrl = userData?.profile_picture || "";

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
>
  <div class="container">
    <ButtonList
      bind:value={selectedCategory}
      options={categories} 
    />
    <main>
      {#if userData != null && userSettings != null && globalSettings != null}
        {#if selectedCategory === "account"}
          <TextInput
            name="username"
            placeholder="Username"
            bind:value={userData.username}
            validation={isValidUsername}
          />
          <TextInput
            name="email"
            placeholder="Email"
            bind:value={userData.email}
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
            bind:value={userData.searchable}
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
            bind:value={userSettings[UserSettingKeys.DisplayAllDayEventsFilled]}
          />
          <ToggleInput
            name={UserSettingKeys.DisplayNonAllDayEventsFilled}
            description="Fill Non-All-Day Events"
            bind:value={userSettings[UserSettingKeys.DisplayNonAllDayEventsFilled]}
          />
          <ToggleInput
            name={UserSettingKeys.DisplaySmallCalendar}
            description="Display Small Calendar"
            bind:value={userSettings[UserSettingKeys.DisplaySmallCalendar]}
          />
          <ToggleInput
            name={UserSettingKeys.DynamicCalendarRows}
            description="Dynamic Calendar Row Count"
            bind:value={userSettings[UserSettingKeys.DynamicCalendarRows]}
          />
          <ToggleInput
            name={UserSettingKeys.DynamicSmallCalendarRows}
            description="Dynamic Small Calendar Row Count"
            bind:value={userSettings[UserSettingKeys.DynamicSmallCalendarRows]}
          />
          <ToggleInput
            name={UserSettingKeys.DisplayWeekNumbers}
            description="Display Week Numbers"
            bind:value={userSettings[UserSettingKeys.DisplayWeekNumbers]}
          />
          <SelectInput
            name={UserSettingKeys.FirstDayOfWeek}
            placeholder="First Day of Week"
            bind:value={userSettings[UserSettingKeys.FirstDayOfWeek]}
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
            bind:value={userSettings[UserSettingKeys.DisplayRoundedCorners]}
          />
          <SelectInput
            name={UserSettingKeys.ThemeLight}
            placeholder="Light Theme"
            bind:value={userSettings[UserSettingKeys.ThemeLight]}
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
            bind:value={userSettings[UserSettingKeys.ThemeDark]}
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
            bind:value={userSettings[UserSettingKeys.FontText]}
            options={[
              { name: "Atkinson Hyperlegible Next", value: "Atkinson Hyperlegible Next" },
              { name: "Atkinson Hyperlegible Mono", value: "Atkinson Hyperlegible Mono" }
            ]}
          />
          <SelectInput
            name={UserSettingKeys.FontTime}
            placeholder="Time Font"
            bind:value={userSettings[UserSettingKeys.FontTime]}
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
            bind:value={userSettings[UserSettingKeys.DebugMode]}
          />
        {:else if selectedCategory === "admin"}
          <ToggleInput
            name={GlobalSettingKeys.RegistrationEnabled}
            description="Enable Registration"
            bind:value={globalSettings[GlobalSettingKeys.RegistrationEnabled]}
          />
          <ToggleInput
            name={GlobalSettingKeys.UseCdnFonts}
            description="Use Google's CDN for fonts"
            bind:value={globalSettings[GlobalSettingKeys.UseCdnFonts]}
          />
          <SelectButtons
            name={GlobalSettingKeys.LoggingVerbosity}
            bind:value={globalSettings[GlobalSettingKeys.LoggingVerbosity]}
            placeholder="Error Messages Verbosity"
            options={[
              { name: "Broad", value: 3 },
              { name: "Plain", value: 2 },
              { name: "Wordy", value: 1 },
              { name: "Debug", value: 0 }
            ]}
          />
        {/if}
      {:else}
        <Horizontal position="center">
          <Loader/>
        </Horizontal>
      {/if}
    </main>
  </div>
</Modal>
<script lang="ts">
  import { Code, GlobeLock, LockKeyhole, LogOut, Monitor, User, UserSearch } from "lucide-svelte";
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
  import { GlobalSettingKeys, UserSettingKeys } from "../../types/settings";

  interface Props {
    showModal?: () => any;
    hideModal?: () => any;
  }

  let {
    showModal = $bindable(),
    hideModal = $bindable(NoOp),
  }: Props = $props();

  showModal = () => {
    showModalInternal();
  };

  hideModal = () => {
    hideModalInternal();
  };

  let showModalInternal = $state(NoOp);
  let hideModalInternal = $state(NoOp);

  const categories: Option[][] = [
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

  // Account Settings
  let profilePictureType = $state("gravatar");
  let profilePictureFiles: FileList | null = $state(null);

  // Appearance Settings
  let firstDayOfWeek = $state("monday");
  let lightTheme = $state("luna-light");
  let darkTheme = $state("luna-dark");
  let fontText = $state("Atkinson Hyperlegible Next");
  let fontTime = $state("Atkinson Hyperlegible Mono");

  // Developer Settings

  // Admin Settings
  let loggingVerbosity = $state("plain");
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
    min-height: 50vh;
  }

  main {
    display: flex;
    flex-direction: column;
    flex-wrap: nowrap;
    gap: dimensions.$gapMiddle;
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
      {#if selectedCategory === "account"}
        <TextInput
          name="username"
          placeholder="Username"
          validation={isValidUsername}
        />
        <TextInput
          name="email"
          placeholder="Email"
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
        />
        <Horizontal position="justify">
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
            src="https://opisek.net/_app/immutable/assets/portrait.2-Ny2279.webp/"
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
          />
        {/if}
      {:else if selectedCategory === "appearance"}
        <ToggleInput
          name={UserSettingKeys.DisplayAllDayEventsFilled}
          description="Fill All-Day Events"
        />
        <ToggleInput
          name={UserSettingKeys.DisplayAllDayEventsFilled}
          description="Fill Non-All-Day Events"
        />
        <ToggleInput
          name={UserSettingKeys.DisplaySmallCalendar}
          description="Display Small Calendar"
        />
        <ToggleInput
          name={UserSettingKeys.DynamicCalendarRows}
          description="Dynamic Calendar Row Count"
        />
        <ToggleInput
          name={UserSettingKeys.DynamicSmallCalendarRows}
          description="Dynamic Small Calendar Row Count"
        />
        <ToggleInput
          name={UserSettingKeys.DisplayWeekNumbers}
          description="Display Week Numbers"
        />
        <SelectInput
          name={UserSettingKeys.FirstDayOfWeek}
          placeholder="First Day of Week"
          bind:value={firstDayOfWeek}
          options={[
            { name: "Monday", value: "1" },
            { name: "Tuesday", value: "2" },
            { name: "Wednesday", value: "3" },
            { name: "Thursday", value: "4" },
            { name: "Friday", value: "5" },
            { name: "Saturday", value: "6" },
            { name: "Sunday", value: "0" }
          ]}
        />
        <ToggleInput
          name={UserSettingKeys.DisplayRoundedCorners}
          description="Rounded Corners"
        />
        <SelectInput
          name={UserSettingKeys.ThemeLight}
          placeholder="Light Theme"
          bind:value={lightTheme}
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
          bind:value={darkTheme}
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
          bind:value={fontText}
          options={[
            { name: "Atkinson Hyperlegible Next", value: "Atkinson Hyperlegible Next" },
            { name: "Atkinson Hyperlegible Mono", value: "Atkinson Hyperlegible Mono" }
          ]}
        />
        <SelectInput
          name={UserSettingKeys.FontTime}
          placeholder="Time Font"
          bind:value={fontTime}
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
        />
      {:else if selectedCategory === "admin"}
        <ToggleInput
          name={GlobalSettingKeys.RegistrationEnabled}
          description="Enable Registration"
        />
        <ToggleInput
          name={GlobalSettingKeys.UseCdnFonts}
          description="Use Google's CDN for fonts"
        />
        <SelectButtons
          name={GlobalSettingKeys.LoggingVerbosity}
          bind:value={loggingVerbosity}
          placeholder="Error Messages Verbosity"
          options={[
            { name: "Broad", value: "broad" },
            { name: "Plain", value: "plain" },
            { name: "Wordy", value: "wordy" },
            { name: "Debug", value: "debug" }
          ]}
        />
      {/if}
    </main>
  </div>
</Modal>
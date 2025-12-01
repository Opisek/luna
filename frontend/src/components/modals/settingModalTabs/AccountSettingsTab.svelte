<script lang="ts">
  import type { Settings } from "../../../lib/client/data/settings.svelte";
  import { isValidEmail, isValidPassword, isValidRepeatPassword, isValidUsername, valid } from "../../../lib/client/validation";
  import { getSha256Hash } from "../../../lib/common/crypto";
  import { GlobalSettingKeys, UserSettingKeys, type UserData } from "../../../types/settings";
  import FileUpload from "../../forms/FileUpload.svelte";
  import SelectButtons from "../../forms/SelectButtons.svelte";
  import TextInput from "../../forms/TextInput.svelte";
  import ToggleInput from "../../forms/ToggleInput.svelte";
  import Horizontal from "../../layout/Horizontal.svelte";
  import Image from "../../layout/Image.svelte";

  interface Props {
    settings: Settings;
    userDataSnapshot: UserData | null;
    accountSettingsSubmittable: boolean;
    accountSettingsReauthenticationRequired: boolean;
    accountSettingsNewProfilePictureChosen: boolean;
    accountSettingsNewProfilePictureUrl: string;
    accountSettingsNewProfilePictureFile: File | null;
    newPassword: string;
    refetchProfilePicture: () => void;
  }

  let {
    settings,
    userDataSnapshot,
    accountSettingsSubmittable = $bindable(),
    accountSettingsReauthenticationRequired = $bindable(),
    accountSettingsNewProfilePictureChosen = $bindable(),
    accountSettingsNewProfilePictureUrl = $bindable(""),
    accountSettingsNewProfilePictureFile = $bindable(null),
    newPassword = $bindable(""),
    refetchProfilePicture = $bindable()
  }: Props = $props();

  let usernameValidity = $state(valid);
  let emailValidity = $state(valid);
  let passwordValidity = $state(valid);
  let repeatPasswordValidity = $state(valid);

  // Submission logic
  // Checks whether the inputted data is valid and different from the snapshot.
  // This is reported back to the settings modal which enables/disables the
  // submit button.
  let submittable = $derived(
    (userDataSnapshot !== null) && (
      (userDataSnapshot.username === settings.userData.username || usernameValidity.valid) &&
      (userDataSnapshot.email === settings.userData.email || emailValidity.valid) &&
      (newPassword === "" || (passwordValidity.valid && repeatPasswordValidity.valid))
    )
  );
  $effect(() => {
    accountSettingsSubmittable = submittable;
  });

  let reauthenticationRequired = $derived(
    (userDataSnapshot !== null) && (
      (userDataSnapshot.username !== settings.userData.username) ||
      (userDataSnapshot.email !== settings.userData.email) ||
      (newPassword !== "")
    )
  );
  $effect(() => {
    accountSettingsReauthenticationRequired = reauthenticationRequired;
  });

  let newProfilePictureChosen = $derived.by(() => {
    if (!userDataSnapshot) return false;
    if (settings.userData.profile_picture_type !== userDataSnapshot.profile_picture_type) return true;
    switch (settings.userData.profile_picture_type) {
      case "static":
        return profilePictureStaticUrl !== userDataSnapshot.profile_picture_url;
      case "remote":
        return profilePictureRemoteUrl !== userDataSnapshot.profile_picture_url;
      case "gravatar":
        return profilePictureGravatarForceDefault !== userDataSnapshot.profile_picture_url.includes("f=y");
      case "database":
        return profilePictureFileId === "" && profilePictureFiles !== null;
      default:
        return false;
    }
  });
  $effect(() => {
    accountSettingsNewProfilePictureChosen = newProfilePictureChosen;
  });

  // Remote profile picture logic
  // Keeps track of the static profile picture URL so the information is not
  // lost when switching types.
  let profilePictureStaticUrl = $state("");

  // Remote profile picture logic
  // Keeps track of the remote profile picture URL so the information is not
  // lost when switching types.
  let profilePictureRemoteUrl = $state("");

  // Uploaded profile picture logic
  // Keeps track of the uploaded profile picture.
  // If the user already uploaded a profile picture, it is referenced by its
  // database ID.
  // If the user chooses to upload a different picture, the ID is reset and a
  // dummy URL is constructed for the preview.
  let profilePictureFiles: FileList | null = $state(null);
  let profilePictureFileId = $state("");

  let uploadedProfilePictureLocalUrl = $state("");
  $effect(() => {
    if (profilePictureFileId != "" || !profilePictureFiles) {
      uploadedProfilePictureLocalUrl = "";
      accountSettingsNewProfilePictureFile = null;
    } else {
      const file = profilePictureFiles[0];
      accountSettingsNewProfilePictureFile = file;
      const reader = new FileReader();
      reader.onload = () => {
        uploadedProfilePictureLocalUrl = reader.result as string;
      };
      reader.readAsDataURL(file);
    }
  });

  // Gravatar logic
  // Keeps track of whether the user chooses to force the default gravatar
  // picture to be displayed instead of their actual one.
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

  let gravatarCheckDefaultTimeout = $state<ReturnType<typeof setTimeout>>();
  $effect(() => {
    if (settings.userData.profile_picture_type !== "gravatar") return;
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

  // Profile picture preview logic
  // Determines the profile picture URL for the preview based on the selected type.
  let effectiveProfilePictureUrl = $derived.by(() => {
    switch (settings.userData.profile_picture_type) {
      case "gravatar":
        return profilePictureGravatarUrl;
      case "remote":
        return profilePictureRemoteUrl;
      case "static":
        return profilePictureStaticUrl;
      case "database":
        if (profilePictureFileId !== "" && profilePictureFileId !== null) return `/api/files/${profilePictureFileId}`
        else return uploadedProfilePictureLocalUrl;
      default:
        return "";
    }
  });
  $effect(() => {
    accountSettingsNewProfilePictureUrl = effectiveProfilePictureUrl;
  });

  // Profile picture initialization logic
  // When the tab is first entered, the user's current profile picture is
  // matched to its type and relevant metadata.
  refetchProfilePicture = (() => {
    if (!userDataSnapshot) return;
    switch (userDataSnapshot.profile_picture_type) {
      case "static":
        profilePictureStaticUrl = userDataSnapshot.profile_picture_url;
        break;
      case "remote":
        profilePictureRemoteUrl = userDataSnapshot.profile_picture_url;
        break;
      case "database":
        profilePictureFileId = userDataSnapshot.profile_picture_file;
        break;
      case "gravatar":
        profilePictureGravatarForceDefault = userDataSnapshot.profile_picture_url.includes("f=y");
        break;
    }
  });
</script>

<style lang="scss">
  @use "../../../styles/dimensions.scss";

  div.pfpButtons {
    display: flex;
    flex-direction: column;
    gap: dimensions.$gapMiddle;
    width: 100%;
  }
</style>

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
      bind:value={settings.userData.profile_picture_type}
      options={(!settings.globalSettings[GlobalSettingKeys.EnableGravatar] ? [] : [
        { name: "Gravatar", value: "gravatar" }
      ]).concat([
        { name: "Upload File", value: "database" },
        { name: "Internet Link", value: "remote" },
        { name: "Luna Art", value: "static" }
      ])}
    />
  </div>
  <Image
    src={effectiveProfilePictureUrl}
    alt="Profile Picture"
  />
</Horizontal>
{#if settings.userData.profile_picture_type === "database"}
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
{:else if settings.userData.profile_picture_type === "remote"}
  <TextInput
    name="pfp_link"
    placeholder="Profile Picture Link"
    bind:value={profilePictureRemoteUrl}
  />
{:else if settings.userData.profile_picture_type === "gravatar"}
  {#if !profilePictureGravatarIsDefault && settings.globalSettings[GlobalSettingKeys.EnableGravatar]}
    <ToggleInput
      name="pfp_gravatar_force_default"
      description="Use Default Gravatar Profile Picture"
      bind:value={profilePictureGravatarForceDefault}
    />
  {/if}
{:else if settings.userData.profile_picture_type === "static"}
  Feature not yet available
{/if}
{#if settings.userData.id && settings.userSettings[UserSettingKeys.DebugMode]}
  <TextInput value={settings.userData.id} name="id" placeholder="User ID" editable={false} />
{/if}
<script lang="ts">
  import { t } from "@sveltia/i18n";
  import type { Settings } from "../../../lib/client/data/settings.svelte";
  import { fetchResponse } from "../../../lib/client/net";
  import { NoOp } from "../../../lib/client/placeholders";
  import { ColorKeys } from "../../../types/colors";
  import { GlobalSettingKeys } from "../../../types/settings";
  import SelectButtons from "../../forms/SelectButtons.svelte";
  import ToggleInput from "../../forms/ToggleInput.svelte";
  import Button from "../../interactive/Button.svelte";

  interface Props {
    settings: Settings;
    showConfirmation: (message: string, details?: string) => Promise<void>;
    refetchProfilePicture: () => void;
    snapshotSettings: () => void;
  }

  let {
    settings,
    showConfirmation,
    refetchProfilePicture,
    snapshotSettings,
  }: Props = $props();

  async function resetGlobalSettings() {
    await showConfirmation(`${t("settings.admin.confirm.reset")}\n${t("confirmation.irreversible")}`).then(async () => {
      await fetchResponse("/api/settings", { method: "DELETE" });
      settings.fetchSettings().then(() => {
        snapshotSettings();
        refetchProfilePicture();
      });
    }).catch(NoOp);
  }
</script>

<ToggleInput
  name={GlobalSettingKeys.RegistrationEnabled}
  description={t("settings.admin.registration.display")}
  info={t("settings.admin.registration.info")}
  bind:value={settings.globalSettings[GlobalSettingKeys.RegistrationEnabled]}
/>
<!--
<ToggleInput
  name={GlobalSettingKeys.UseCdnFonts}
  description="Use Google's CDN for fonts"
  bind:value={settings.globalSettings[GlobalSettingKeys.UseCdnFonts]}
/>
-->
<ToggleInput
  name={GlobalSettingKeys.UseIpGeolocation}
  description={t("settings.admin.geolocation.display")}
  info={t("settings.admin.geolocation.info")}
  bind:value={settings.globalSettings[GlobalSettingKeys.UseIpGeolocation]}
/>
<ToggleInput
  name={GlobalSettingKeys.EnableGravatar}
  description={t("settings.admin.pfp.gravatar.display")}
  info={t("settings.admin.pfp.gravatar.info")}
  bind:value={settings.globalSettings[GlobalSettingKeys.EnableGravatar]}
/>
<ToggleInput
  name={GlobalSettingKeys.CacheProfilePictures}
  description={t("settings.admin.pfp.cache.display")}
  info={t("settings.admin.pfp.cache.info")}
  bind:value={settings.globalSettings[GlobalSettingKeys.CacheProfilePictures]}
/>
<ToggleInput
  name={GlobalSettingKeys.EnableProfilePicturesUpload}
  description={t("settings.admin.pfp.upload.display")}
  info={t("settings.admin.pfp.upload.info")}
  bind:value={settings.globalSettings[GlobalSettingKeys.EnableProfilePicturesUpload]}
/>
<SelectButtons
  name={GlobalSettingKeys.LoggingVerbosity}
  bind:value={settings.globalSettings[GlobalSettingKeys.LoggingVerbosity]}
  placeholder={t("settings.admin.verbosity.display")}
  info={t("settings.admin.verbosity.info")}
  options={[
    { name: t("settings.admin.verbosity.lvl.broad"), value: 3 },
    { name: t("settings.admin.verbosity.lvl.plain"), value: 2 },
    { name: t("settings.admin.verbosity.lvl.wordy"), value: 1 },
    { name: t("settings.admin.verbosity.lvl.debug"), value: 0 }
  ]}
/>
<Button color={ColorKeys.Danger} onClick={resetGlobalSettings}>{t("settings.admin.reset")}</Button>
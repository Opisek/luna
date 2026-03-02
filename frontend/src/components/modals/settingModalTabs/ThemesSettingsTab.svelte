<script lang="ts">
  import { Download, Sun, Trash2 } from "lucide-svelte";
  import { downloadFileToClient, fetchResponse } from "../../../lib/client/net";
  import Spinner from "../../decoration/Spinner.svelte";
  import FileUpload from "../../forms/FileUpload.svelte";
  import List from "../../forms/List.svelte";
  import Button from "../../interactive/Button.svelte";
  import IconButton from "../../interactive/IconButton.svelte";
  import type { Option } from "../../../types/options";
  import { queueNotification } from "../../../lib/client/notifications";
  import { ColorKeys } from "../../../types/colors";

  interface Props {
    lightThemes: Option<string>[];
    darkThemes: Option<string>[];
    fetchThemes: () => void;
    showConfirmation: (message: string, onConfirm: () => Promise<void>, confirmText?: string, onCancel?: () => Promise<void>, cancelText?: string) => void;
  }

  let {
    lightThemes,
    darkThemes,
    fetchThemes,
    showConfirmation,
  }: Props = $props();

  let themeFile: FileList | null = $state(null);
  let themeFileId = $state("");
  let uploadingThemeFile = $state(false);

  async function uploadThemeFile() {
    if (uploadingThemeFile) return;
    
    if (themeFile == null) {
      queueNotification(ColorKeys.Danger, "Missing or corrupted theme file");
      return;
    }

    uploadingThemeFile = true;

    const themeFiles = lightThemes.concat(darkThemes).map(x => x.value.split("/").pop() + ".css");

    if (themeFiles.includes(themeFile[0].name)) {
      if (!(await new Promise<boolean>((resolve) => {
        showConfirmation(
          "A theme with the same file name already exists.\nContinuing will overwrite that theme.\nAre you sure you want to proceed?",
          async () => {resolve(true)}, "", async () => {resolve(false)}
        );
      }))) {
        uploadingThemeFile = false;
        return;
      }
    }

    const formData = new FormData();
    formData.append("file", themeFile[0]);

    await fetchResponse("/installed/themes", {
      method: "PUT",
      body: formData,
    }).then(async () => {
      fetchThemes();
      queueNotification(ColorKeys.Success, "Theme installed successfully");
      themeFile = null;
    }).catch((err) => {
      queueNotification(ColorKeys.Danger, "Failed to install theme: " + err);
    });

    uploadingThemeFile = false;
  }

  async function deleteTheme(theme: string, name: string, isLightTheme: boolean) {
    if (!(await new Promise<boolean>((resolve) => {
      showConfirmation(
        `Are you sure you want to uninstall the theme "${name}"?\nThis action is irreversible.`,
        async () => {resolve(true)}, "", async () => {resolve(false)}
      );
    }))) return;

    await fetchResponse(`/installed/themes/${isLightTheme ? "light" : "dark"}/${theme}`, {
      method: "DELETE",
    }).then(async () => {
      fetchThemes();
      queueNotification(ColorKeys.Success, "Theme uninstalled successfully");
      themeFile = null;
    }).catch((err) => {
      queueNotification(ColorKeys.Danger, "Failed to uninstall theme: " + err);
    });
  }
</script>

<style lang="scss">
  @use "../../../styles/dimensions.scss";
  @use "../../../styles/colors.scss";
  @use "../../../styles/text.scss";

  .installedResource {
    padding: dimensions.$gapMiddle;
    background-color: colors.$backgroundSecondary;
    color: colors.$foregroundSecondary;
    border-radius: dimensions.$borderRadius;

    display: grid;
    gap: dimensions.$gapSmall;
    row-gap: 0;
    grid-template-columns: 1fr auto;
    grid-template-rows: auto auto;
    grid-template-areas: "name buttons" "details buttons";
    justify-content: center;
    align-items: center;
  }
  .installedResource > .name {
    grid-area: name;
    display: flex;
    flex-direction: row;
    justify-content: start;
    align-items: center;
    gap: dimensions.$gapTiny;
  }
  .installedResource > .details {
    grid-area: details;
    font-size: text.$fontSizeSmall;
  }
  .installedResource > .buttons {
    grid-area: buttons;
    display: flex;
    flex-direction: row;
    gap: dimensions.$gapSmall;
  }
</style>

<FileUpload
  name="theme_file"
  placeholder="Install a Theme"
  bind:files={themeFile}
  bind:fileId={themeFileId}
  accept={".css"}
/>
{#if themeFile !== null}
  <Button color={ColorKeys.Success} onClick={uploadThemeFile}>
    {#if uploadingThemeFile}
      <Spinner/> <!-- TODO: Spinner does not have the same height as text -->
    {:else}
      Upload Theme
    {/if}
  </Button>
{/if}
<List
  label="Installed Themes"
  info={"Looking to change your current theme? Head to the \"Appearance\" tab."}
  items={lightThemes.concat(darkThemes)}
  id={item => item.value}
  template={themeTemplate}
/>

{#snippet themeTemplate(theme: Option<string>)}
  {@const Icon = theme.icon}
  {@const isLightTheme = theme.icon === Sun}

  <div class="installedResource">
    <span class="name">
      {theme.name}

      <Icon size={16}/>
    </span>

    <span class="details">
      {theme.value}.css
    </span>

    <div class="buttons">
      <IconButton click={() => { downloadFileToClient(`/themes/${isLightTheme ? "light" : "dark"}/${theme.value}.css`); }}>
        <Download size={20}/>
      </IconButton>
      <IconButton click={() => { deleteTheme(theme.value, theme.name, isLightTheme); }}>
        <Trash2 size={20}/>
      </IconButton>
    </div>
  </div>
{/snippet}
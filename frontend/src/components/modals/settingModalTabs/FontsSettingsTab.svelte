<script lang="ts">
  import { Download, Trash2 } from "lucide-svelte";
  import { downloadFileToClient, fetchResponse } from "../../../lib/client/net";
  import { queueNotification } from "../../../lib/client/notifications";
  import { ColorKeys } from "../../../types/colors";
  import Spinner from "../../decoration/Spinner.svelte";
  import FileUpload from "../../forms/FileUpload.svelte";
  import List from "../../forms/List.svelte";
  import Button from "../../interactive/Button.svelte";
  import IconButton from "../../interactive/IconButton.svelte";
  import type { Option } from "../../../types/options";

  interface Props {
    fonts: Option<string>[];
    fetchFonts: () => void;
    showConfirmation: (message: string, onConfirm: () => Promise<void>, confirmText?: string, onCancel?: () => Promise<void>, cancelText?: string) => void;
  }

  let {
    fonts,
    fetchFonts,
    showConfirmation,
  }: Props = $props();

  let fontFile: FileList | null = $state(null);
  let fontFileId = $state("");
  let uploadingFontFile = $state(false);

  async function uploadFontFile() {
    if (uploadingFontFile) return;
    
    if (fontFile == null) {
      queueNotification(ColorKeys.Danger, "Missing or corrupted font file");
      return;
    }

    uploadingFontFile = true;

    const fontFiles = fonts.map(x => x.value.split("/").pop() + ".ttf");

    if (fontFiles.includes(fontFile[0].name)) {
      if (!(await new Promise<boolean>((resolve) => {
        showConfirmation(
          "A font with the same file name already exists.\nContinuing will overwrite that font.\nAre you sure you want to proceed?",
          async () => {resolve(true)}, "", async () => {resolve(false)}
        );
      }))) {
        uploadingFontFile = false;
        return;
      }
    }

    const formData = new FormData();
    formData.append("file", fontFile[0]);

    await fetchResponse("/installed/fonts", {
      method: "PUT",
      body: formData,
    }).then(async () => {
      fetchFonts();
      queueNotification(ColorKeys.Success, "Font installed successfully");
      fontFile = null;
    }).catch((err) => {
      queueNotification(ColorKeys.Danger, "Failed to install font: " + err);
    });

    uploadingFontFile = false;
  }

  async function deleteFont(font: string, name: string) {
    if (!(await new Promise<boolean>((resolve) => {
      showConfirmation(
        `Are you sure you want to uninstall the font "${name}"?\nThis action is irreversible.`,
        async () => {resolve(true)}, "", async () => {resolve(false)}
      );
    }))) return;

    await fetchResponse(`/installed/fonts/${font}`, {
      method: "DELETE",
    }).then(async () => {
      fetchFonts();
      queueNotification(ColorKeys.Success, "Font uninstalled successfully");
      fontFile = null;
    }).catch((err) => {
      queueNotification(ColorKeys.Danger, "Failed to uninstall font: " + err);
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
  placeholder="Install a Font"
  bind:files={fontFile}
  bind:fileId={fontFileId}
  accept={".ttf"}
/>
{#if fontFile !== null}
  <Button color={ColorKeys.Success} onClick={uploadFontFile}>
    {#if uploadingFontFile}
      <Spinner/> <!-- TODO: Spinner does not have the same height as text -->
    {:else}
      Upload Theme
    {/if}
  </Button>
{/if}
<List
  label="Installed Fonts"
  info={"Looking to change your current font? Head to the \"Appearance\" tab."}
  items={fonts}
  id={item => item.value}
  template={fontTemplate}
/>

{#snippet fontTemplate(font: Option<string>)}
  <div class="installedResource">
    <span class="name">
      {font.name}
    </span>

    <span class="details">
      {font.value}.ttf
    </span>

    <div class="buttons">
      <IconButton click={() => { downloadFileToClient(`/fonts/${font.value}.ttf`); }}>
        <Download size={20}/>
      </IconButton>
      <IconButton click={() => { deleteFont(font.value, font.name); }}>
        <Trash2 size={20}/>
      </IconButton>
    </div>
  </div>
{/snippet}
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
  import { t } from "@sveltia/i18n";

  interface Props {
    fonts: Option<string>[];
    fetchFonts: () => void;
    showConfirmation: (message: string, confirmText?: string, cancelText?: string) => Promise<void>;
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
      queueNotification(ColorKeys.Danger, t("settings.fonts.error.file"));
      return;
    }

    uploadingFontFile = true;

    const fontFiles = fonts.map(x => x.value.split("/").pop() + ".ttf");

    if (fontFiles.includes(fontFile[0].name)) {
      const confirmed = await showConfirmation(t("settings.fonts.confirm.overwrite")).then(() => true).catch(() => false);
      if (!confirmed) {
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
      queueNotification(ColorKeys.Success, t("settings.fonts.success.install"));
      fontFile = null;
    }).catch((err) => {
      queueNotification(ColorKeys.Danger, t("settings.fonts.error.install", { values: { msg: err.message } }));
    });

    uploadingFontFile = false;
  }

  async function deleteFont(font: string, name: string) {
    const confirmed = await showConfirmation(`${t("settings.fonts.confirm.uninstall", { values: { name: name } })}\n${t("confirmation.irreversible")}`).then(() => true).catch(() => false);
    if (!confirmed) return;

    await fetchResponse(`/installed/fonts/${font}`, {
      method: "DELETE",
    }).then(async () => {
      fetchFonts();
      queueNotification(ColorKeys.Success, t("settings.fonts.success.uninstall"));
      fontFile = null;
    }).catch((err) => {
      queueNotification(ColorKeys.Danger, t("settings.fonts.error.uninstall", { values: { msg: err.message } }));
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
  name="font_file"
  placeholder={t("settings.fonts.new")}
  bind:files={fontFile}
  bind:fileId={fontFileId}
  accept={".ttf"}
/>
{#if fontFile !== null}
  <Button color={ColorKeys.Success} onClick={uploadFontFile}>
    {#if uploadingFontFile}
      <Spinner/>
    {:else}
      {t("settings.fonts.upload")}
    {/if}
  </Button>
{/if}
<List
  label={t("settings.fonts.list")}
  info={t("settings.fonts.current")}
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
      <IconButton onClick={() => { downloadFileToClient(`/fonts/${font.value}.ttf`); }} color={ColorKeys.Accent} alt={t("button.download")}>
        <Download size={20}/>
      </IconButton>
      <IconButton onClick={async () => deleteFont(font.value, font.name)} color={ColorKeys.Danger} alt={t("button.delete")}>
        <Trash2 size={20}/>
      </IconButton>
    </div>
  </div>
{/snippet}
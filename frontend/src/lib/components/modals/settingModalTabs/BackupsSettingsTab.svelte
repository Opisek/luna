<script lang="ts">
  import { t } from "@sveltia/i18n";
  import type { Settings } from "$lib/scripts/client/data/settings.svelte";
  import { ColorKeys } from "$lib/types/colors";
  import Button from "../../interactive/Button.svelte";
  import { downloadFileToClient, fetchFile, fetchResponse } from "$lib/scripts/client/net";
  import { queueNotification } from "$lib/scripts/client/notifications";
  import FileUpload from "../../forms/FileUpload.svelte";
  import Modal from "../Modal.svelte";
  import { NoOp } from "$lib/scripts/client/placeholders";
  import IconButton from "../../interactive/IconButton.svelte";
  import { Check, X } from "lucide-svelte";
  import ToggleInput from "../../forms/ToggleInput.svelte";
  import TextInput from "../../forms/TextInput.svelte";

  interface Props {
    settings: Settings;
    showConfirmation: (message: string, details?: string, notice?: boolean) => Promise<void>;
  }

  let {
    settings,
    showConfirmation,
  }: Props = $props();

  let showEncryptionModal: () => Promise<string> = $state(Promise.reject);
  let success: (result: string) => void = $state(NoOp);
  let failure: (reason?: string | Error) => void = $state(NoOp);

  let enableEncryption = $state(true);
  let creatingBackup = $state(false);
  let backupPassword = $state("");

  let backupFile: FileList | null = $state(null);

  async function createBackup() {
    creatingBackup = true;
    backupPassword = "";
    const password = await showEncryptionModal().catch(() => null);
    if (password === null) {
      backupPassword = "";
      return;
    }

    const formData = new FormData();
    if (password != "") formData.append("backup_password", password);

    await fetchFile("/api/backups/create", { method: "POST", body: formData }).then(downloadFileToClient).catch(e => {
      queueNotification(ColorKeys.Danger, e);
    }).finally(() => {
      backupPassword = "";
    })
  }

  async function restoreBackup() {
    if (backupFile === null) {
      queueNotification(ColorKeys.Danger, t("settings.backups.error.file"));
      return;
    }

    enableEncryption = backupFile[0].name.endsWith(".encrypted");
    creatingBackup = false;
    backupPassword = "";
    const password = await showEncryptionModal().catch(() => null);
    if (password === null) {
      backupPassword = "";
      return;
    }

    const confirmed = await showConfirmation(`${t("settings.backups.confirm")}\n${t("confirmation.irreversible")}`).then(() => true).catch(() => false);
    if (!confirmed) return;

    const formData = new FormData();
    formData.append("backup_file", backupFile[0]);
    if (password != "") formData.append("backup_password", password);

    return fetchResponse("/api/backups/restore", { method: "POST", body: formData }).then(() =>
      showConfirmation(t("settings.backups.restart"), "", true)
    ).catch(err => {
      queueNotification(ColorKeys.Danger, t("settings.backups.error.restore", { values: { msg: err.message } }));
    }).finally(() => {
      backupPassword = "";
    })
  }
</script>

<Button color={ColorKeys.Accent} onClick={createBackup}>{t("settings.backups.create")}</Button>

<FileUpload
  name="backup_file"
  placeholder={t("settings.backups.file")}
  bind:files={backupFile}
  accept=".tar.gz,.tar.gz.encrypted"
/>
{#if backupFile !== null}
  <Button color={ColorKeys.Success} onClick={restoreBackup}>
    {t("settings.backups.restore")}
  </Button>
{/if}

<Modal
  title={t("settings.backups.encryption.title")}
  bind:showModal={showEncryptionModal}
  bind:success
  bind:failure
>
  {creatingBackup ? t("settings.backups.encryption.desc.create") : t("settings.backups.encryption.desc.restore")}
  <ToggleInput
    name="backup_encrypted"
    description={creatingBackup ? t("settings.backups.encryption.enable.create") : t("settings.backups.encryption.enable.restore")}
    bind:value={enableEncryption}
  />
  {#if enableEncryption}
    <TextInput bind:value={backupPassword} name="backup_password" placeholder={t("settings.backups.encryption.password")} password={true} />
  {/if}
  {#snippet buttons()}
    <IconButton enabled={!enableEncryption || backupPassword.length > 0} onClick={() => success(enableEncryption ? backupPassword : "")} color={ColorKeys.Success} type="submit" alt={t("button.confirm")} canRenderAsButton={true}>
      <Check/>
    </IconButton>
    <IconButton onClick={failure} color={ColorKeys.Danger} alt={t("button.cancel")} canRenderAsButton={true}><X/></IconButton>
  {/snippet}
</Modal>
<script lang="ts">
  import type { Settings } from "../../../lib/client/data/settings.svelte";
  import { fetchResponse } from "../../../lib/client/net";
  import { ColorKeys } from "../../../types/colors";
  import Button from "../../interactive/Button.svelte";

  interface Props {
    settings: Settings;
    requirePasswordForAccountDeletion: boolean;
    showConfirmation: (message: string, onConfirm: () => Promise<void>, confirmText?: string, onCancel?: () => Promise<void>, cancelText?: string) => void;
    deleteAccount: () => void;
    refetchProfilePicture: () => void;
    snapshotSettings: () => void;
  }

  let {
    settings,
    requirePasswordForAccountDeletion = $bindable(),
    showConfirmation,
    deleteAccount,
    refetchProfilePicture,
    snapshotSettings,
  }: Props = $props();

  function resetPreferences() {
    showConfirmation("Are you sure you want to reset all your preferences?\nThis action is irreversible.", async () => {
      await fetchResponse("/api/users/self/settings", { method: "DELETE" });
      settings.fetchSettings().then(() => {
        snapshotSettings();
        refetchProfilePicture();
      });
    }, "Your account will remain intact.");
  }
</script>

<Button color={ColorKeys.Danger} onClick={resetPreferences}>Reset all my preferences</Button>
<Button color={ColorKeys.Danger} onClick={() => deleteAccount()}>Delete my account</Button>
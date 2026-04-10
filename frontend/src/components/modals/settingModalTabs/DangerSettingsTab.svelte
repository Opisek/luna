<script lang="ts">
  import type { Settings } from "../../../lib/client/data/settings.svelte";
  import { fetchResponse } from "../../../lib/client/net";
  import { NoOp } from "../../../lib/client/placeholders";
  import { ColorKeys } from "../../../types/colors";
  import Button from "../../interactive/Button.svelte";

  interface Props {
    settings: Settings;
    showConfirmation: (message: string, confirmText?: string, cancelText?: string) => Promise<void>;
    deleteAccount: () => Promise<void>;
    refetchProfilePicture: () => void;
    snapshotSettings: () => void;
  }

  let {
    settings,
    showConfirmation,
    deleteAccount,
    refetchProfilePicture,
    snapshotSettings,
  }: Props = $props();

  async function resetPreferences() {
    await showConfirmation(
      "Are you sure you want to reset all your preferences?\nThis action is irreversible.",
      "Your account will remain intact."
    ).then(async () => {
      await fetchResponse("/api/users/self/settings", { method: "DELETE" });
      settings.fetchSettings().then(() => {
        snapshotSettings();
        refetchProfilePicture();
      });
    }).catch(NoOp);
  }
</script>

<Button color={ColorKeys.Danger} onClick={resetPreferences}>Reset all my preferences</Button>
<Button color={ColorKeys.Danger} onClick={() => deleteAccount()}>Delete my account</Button>
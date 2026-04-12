<script lang="ts">
  import Modal from "./Modal.svelte";

  import { NoOp } from "../../lib/client/placeholders";
  import { ColorKeys } from "../../types/colors";
  import { Check, X } from "lucide-svelte";
  import IconButton from "../interactive/IconButton.svelte";
  import { t } from "@sveltia/i18n";

  interface Props {
    showModal: () => Promise<void>;
    children?: import('svelte').Snippet;
  }

  let success: (result: void) => void = $state(NoOp);
  let failure: (reason?: string | Error) => void = $state(NoOp);

  let {
    showModal = $bindable(),
    children
  }: Props = $props();
</script>

<Modal title={t("confirmation.title")} bind:showModal bind:success bind:failure>
  {@render children?.()}
  {#snippet buttons()}
    <IconButton onClick={success} color={ColorKeys.Success} type="submit" alt={t("button.confirm")} canRenderAsButton={true}>
      <Check/>
    </IconButton>
    <IconButton onClick={failure} color={ColorKeys.Danger} alt={t("button.cancel")} canRenderAsButton={true}><X/></IconButton>
  {/snippet}
</Modal>
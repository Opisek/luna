<script lang="ts">
  import Modal from "./Modal.svelte";

  import { NoOp } from "$lib/scripts/client/placeholders";
  import { ColorKeys } from "$lib/types/colors";
  import { Check, X } from "lucide-svelte";
  import IconButton from "../interactive/IconButton.svelte";
  import { t } from "@sveltia/i18n";

  interface Props {
    showModal: () => Promise<void>;
    isNotice?: boolean;
    children?: import('svelte').Snippet;
  }

  let {
    showModal = $bindable(),
    isNotice = false,
    children
  }: Props = $props();

  let success: (result: void) => void = $state(NoOp);
  let failure: (reason?: string | Error) => void = $state(NoOp);
</script>

<Modal title={isNotice ? t("confirmation.notice") : t("confirmation.title")} bind:showModal bind:success bind:failure>
  {@render children?.()}
  {#snippet buttons()}
    <IconButton onClick={success} color={ColorKeys.Success} type="submit" alt={t("button.confirm")} canRenderAsButton={true}>
      <Check/>
    </IconButton>
    {#if !isNotice}
      <IconButton onClick={failure} color={ColorKeys.Danger} alt={t("button.cancel")} canRenderAsButton={true}><X/></IconButton>
    {/if}
  {/snippet}
</Modal>
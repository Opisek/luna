<script lang="ts">
  import { Eye, EyeOff } from "lucide-svelte";

  import IconButton from "./IconButton.svelte";
  import { NoOp } from "../../lib/client/placeholders";

  import { t } from "@sveltia/i18n";

  interface Props {
    visible: boolean;
    momentary?: boolean;
    onClick?: (visible: boolean) => any;
  }

  let {
    visible = $bindable(),
    momentary = false,
    onClick = NoOp,
  }: Props = $props();

  function toggleVisibility() {
    visible = !visible;
    onClick(visible);
  }

  function hide() {
    visible = false;
    onClick(visible);
  }

  function show() {
    visible = true;
    onClick(visible);
  }
</script>

{#if momentary}
  <IconButton down={show} up={hide} tabindex={-1} alt={t("button.show")}>
    {@render icon()}
  </IconButton>
{:else}
  <IconButton onClick={toggleVisibility} alt={visible ? t("button.hide") : t("button.show")}>
    {@render icon()}
  </IconButton>
{/if}

{#snippet icon()}
  {#if visible}
    <Eye size={16}/>
  {:else}
    <EyeOff size={16}/>
  {/if}
{/snippet}
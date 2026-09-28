<script lang="ts">
  import { X } from "lucide-svelte";

  import IconButton from "./IconButton.svelte";

  import { t } from "@sveltia/i18n";

  interface Props {
    onClick: () => void;
  }

  let {
    onClick = $bindable()
  }: Props = $props();
</script>

<IconButton onClick={onClick} alt={t("button.close")}>
  <X/>
</IconButton>
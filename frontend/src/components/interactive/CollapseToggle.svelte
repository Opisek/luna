<script lang="ts">
  import { ChevronDown } from "lucide-svelte";

  import IconButton from "./IconButton.svelte";
  import { t } from "@sveltia/i18n";

  interface Props {
    collapsed: boolean;
  }

  let {
    collapsed = $bindable()
  }: Props = $props();

  function toggleCollapse() {
    collapsed = !collapsed;
  }
</script>

<IconButton onClick={toggleCollapse} style="transform: rotate({collapsed ? "-90deg" : "0deg"})" alt={collapsed ? t("button.hide") : t("button.show")}>
  <ChevronDown size={16}/>
</IconButton>
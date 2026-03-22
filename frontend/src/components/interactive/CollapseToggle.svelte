<script lang="ts">
  import { ChevronDown } from "lucide-svelte";

  import IconButton from "./IconButton.svelte";

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

<IconButton onClick={toggleCollapse} style="transform: rotate({collapsed ? "-90deg" : "0deg"})" alt={collapsed ? "Hide" : "Show"}>
  <ChevronDown size={16}/>
</IconButton>
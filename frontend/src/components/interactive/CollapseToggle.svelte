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

<IconButton click={toggleCollapse} style="transform: rotate({collapsed ? "-90deg" : "0deg"})">
  <ChevronDown size={16}/>
</IconButton>
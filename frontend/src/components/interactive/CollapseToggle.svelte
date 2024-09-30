<script lang="ts">
  import { ChevronDown } from "lucide-svelte";
  import IconButton from "./IconButton.svelte";

  export let collapsed: boolean;

  function toggleCollapse() {
    collapsed = !collapsed;
  }
</script>

<style lang="scss">
  @import "../../styles/animations.scss";

  span {
    transition: all $cubic $animationSpeed;
    width: 100%;
    height: 100%;
  }
</style>

<IconButton click={toggleCollapse} style="transform: rotate({collapsed ? "-90deg" : "0deg"})">
  <ChevronDown size={16}/>
</IconButton>
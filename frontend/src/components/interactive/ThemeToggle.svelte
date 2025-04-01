<script lang="ts">
  import MoonIcon from "lucide-svelte/icons/moon";
  import SunIcon from "lucide-svelte/icons/sun";

  import { getTheme } from "$lib/client/theme.svelte";

  import IconButton from "./IconButton.svelte";

  const theme = getTheme();
</script>

<style lang="scss">
  @use "../../styles/animations.scss";

  // https://github.com/Opisek/opifolio-v2/blob/main/src/components/interactive/Theme.svelte
  span {
    position: relative;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: transform animations.$animationSpeedSlow animations.$cubic; 
  }
  span.light {
    transform: rotate(360deg);
  }
  span.dark {
    transform: rotate(90deg);
  }
  span :global(*) {
    transition: opacity animations.$animationSpeed;
  }
  span :global(> *:first-child) {
    top: 0;
    left: 0;
    position: absolute;
  }
  span.dark :global(> *:first-child) {
    opacity: 0;
  }
  span.light :global(> *:nth-child(2)) {
    opacity: 0;
  }
</style>

<IconButton click={() => theme.toggle()}>
  <span class:dark={theme.isLightMode()} class:light={!theme.isLightMode()}>
    <MoonIcon/>
    <SunIcon/>
  </span>
</IconButton>
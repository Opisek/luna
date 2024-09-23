<script lang="ts">
  import { Eye, EyeOff } from "lucide-svelte";

  export let visible: boolean;

  function toggleVisibility() {
    visible = !visible;
  }
</script>

<style lang="scss">
  span {
    height: 1.25em;
    width: 1.25em;
    cursor: pointer;

    align-items: center;
    display: flex;
    justify-content: center;
    justify-items: center;
  }
</style>

<span on:click={toggleVisibility}>
  {#if visible}
    <Eye/>
  {:else}
    <EyeOff/>
  {/if}
</span>


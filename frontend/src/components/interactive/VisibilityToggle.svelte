<script lang="ts">
  import { Eye, EyeOff } from "lucide-svelte";

  export let visible: boolean;
  export let momentary: boolean = false;

  function toggleVisibility() {
    visible = !visible;
  }

  function hide() {
    visible = false;
  }

  function show() {
    visible = true;
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

{#if momentary}
  <span
    on:mousedown={show}
    on:mouseup={hide} 
    on:mouseleave={hide} 
  >
    {#if visible}
      <Eye/>
    {:else}
      <EyeOff/>
    {/if}
  </span>
{:else}
  <span on:click={toggleVisibility}>
    {#if visible}
      <Eye/>
    {:else}
      <EyeOff/>
    {/if}
  </span>
{/if}

<!-- TODO: svelte snippets once supported -->
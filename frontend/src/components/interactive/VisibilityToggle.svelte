<script lang="ts">
  import { Eye, EyeOff } from "lucide-svelte";
  import IconButton from "./IconButton.svelte";

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

{#if momentary}
  <IconButton down={show} up={hide} tabindex={-1}>
    {#if visible}
      <Eye size={16}/>
    {:else}
      <EyeOff size={16}/>
    {/if}
  </IconButton>
{:else}
  <IconButton click={toggleVisibility}>
    {#if visible}
      <Eye size={16}/>
    {:else}
      <EyeOff size={16}/>
    {/if}
  </IconButton>
{/if}

<!-- TODO: svelte snippets once supported -->
<script lang="ts">
  import { Eye, EyeOff } from "lucide-svelte";

  import IconButton from "./IconButton.svelte";

  interface Props {
    visible: boolean;
    momentary?: boolean;
  }

  let {
    visible = $bindable(),
    momentary = false
  }: Props = $props();

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
    {@render icon()}
  </IconButton>
{:else}
  <IconButton click={toggleVisibility}>
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
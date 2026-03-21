<script lang="ts">
  import { Eye, EyeOff } from "lucide-svelte";

  import IconButton from "./IconButton.svelte";
  import { NoOp } from "../../lib/client/placeholders";

  interface Props {
    visible: boolean;
    momentary?: boolean;
    onClick?: (visible: boolean) => any;
  }

  let {
    visible = $bindable(),
    momentary = false,
    onClick = NoOp,
  }: Props = $props();

  function toggleVisibility() {
    visible = !visible;
    onClick(visible);
  }

  function hide() {
    visible = false;
    onClick(visible);
  }

  function show() {
    visible = true;
    onClick(visible);
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
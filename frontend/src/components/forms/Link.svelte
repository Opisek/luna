<script lang="ts">
  import type { Snippet } from "svelte";

  import { NoOp } from "$lib/client/placeholders";

  interface Props {
    href?: string;
    onClick?: () => any;
    children?: Snippet;
  }

  let {
    href = "",
    onClick = NoOp,
    children
  }: Props = $props();
</script>

<style lang="scss">
  @use "../../styles/colors.scss";

  a, button {
    text-decoration: none;
    color: colors.$foregroundLink;
    width: max-content;
    outline: none;
    border: none;
    background: none;
    cursor: pointer;
    font-size: inherit;
    font: inherit;
    margin: 0;
    padding: 0;
  }
</style>

{#if href === ""}
  <button onclick={onClick}>{@render children?.()}</button>
{:else}
  <a href={href} onclick={onClick}>{@render children?.()}</a>
{/if}
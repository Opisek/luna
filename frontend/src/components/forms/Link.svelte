<script lang="ts">
  import type { Snippet } from "svelte";

  import { NoOp } from "$lib/client/placeholders";

  interface Props {
    href: string;
    onClick?: () => any;
    children?: Snippet;
  }

  let {
    href,
    onClick = NoOp,
    children
  }: Props = $props();
</script>

<style lang="scss">
  @import "../../styles/colors.scss";

  a {
    text-decoration: none;
    color: $foregroundLink;
    width: max-content;
  }
</style>

<a href={href} onclick={onClick()}>{@render children?.()}</a>
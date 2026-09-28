<script lang="ts">
  import type { Snippet } from "svelte";
  import Tooltip from "../interactive/Tooltip.svelte";

  interface Props {
    children?: Snippet;
  }

  let {
    children
  }: Props = $props();
</script>

<style lang="scss">
  @use "$lib/styles/colors.scss";
  @use "$lib/styles/dimensions.scss";
  @use "$lib/styles/text.scss";

  p {
    color: color-mix(in srgb, colors.$foregroundPrimary 50%, transparent);
  }
</style>

<p>
  {@render children?.()}
</p>
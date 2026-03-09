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
  @use "../../styles/colors.scss";
  @use "../../styles/dimensions.scss";
  @use "../../styles/text.scss";

  p {
    color: color-mix(in srgb, colors.$foregroundPrimary 50%, transparent);
  }
</style>

<p>
  {@render children?.()}
</p>
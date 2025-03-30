<script lang="ts">
  import type { Snippet } from "svelte";

  interface Props {
    children?: Snippet;
  }

  let { children }: Props = $props();
</script>

<style lang="scss">
  @use "../../styles/dimensions.scss";
  @use "../../styles/text.scss";

  h2 {
    padding: 0;
    margin: dimensions.$gapSmaller 0;
    font-weight: text.$fontWeightTitle;
    font-size: text.$fontSizeSubtitle;
  }
</style>

<h2>
  {@render children?.()}
</h2>
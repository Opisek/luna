<script lang="ts">
  import type { Snippet } from "svelte";

  interface Props {
    children?: Snippet;
    id?: string;
  }

  let {
    children,
    id
  }: Props = $props();
</script>

<style lang="scss">
  @use "$lib/styles/dimensions.scss";
  @use "$lib/styles/text.scss";

  h1 {
    padding: 0;
    margin: dimensions.$gapSmall 0 dimensions.$gapSmaller 0;
    font-weight: text.$fontWeightTitle;
    font-size: text.$fontSizeTitle;
  }
</style>

<h1 id={id}>
  {@render children?.()}
</h1>
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

  h1 {
    padding: 0;
    margin: dimensions.$gapSmall 0;
    font-weight: text.$fontWeightTitle;
  }
</style>

<h1>
  {@render children?.()}
</h1>
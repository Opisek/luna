<script lang="ts">
  import type { Snippet } from "svelte";

  interface Props {
    children?: Snippet;
  }

  let { children }: Props = $props();
</script>

<style lang="scss">
  @use "../../styles/text.scss";

  p {
    margin: 0;
    padding: 0;
    text-align: justify;
    line-height: text.$lineHeightParagraph;
    white-space: pre-line;
  }
</style>

<p>
  {@render children?.()}
</p>
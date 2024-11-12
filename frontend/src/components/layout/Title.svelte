<script lang="ts">
  import type { Snippet } from "svelte";

  interface Props {
    children?: Snippet;
  }

  let { children }: Props = $props();
</script>

<style lang="scss">
  @import "../../styles/dimensions.scss";

  h1 {
    padding: 0;
    margin: $gapSmall 0;
  }
</style>

<h1>
  {@render children?.()}
</h1>
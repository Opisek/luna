<script lang="ts">
  import type { Snippet } from "svelte";

  interface Props {
    children?: Snippet;
  }

  let { children }: Props = $props();
</script>

<style lang="scss">
  @use "../../styles/dimensions.scss";

  div {
    display: flex;
    width: 100%;
    flex-direction: column;
    flex-wrap: nowrap;
    gap: dimensions.$gap;
  }
</style>

<div>
  {@render children?.()}
</div>
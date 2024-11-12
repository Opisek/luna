<script lang="ts">
  import type { Snippet } from "svelte";

  interface Props {
    position?: "left" | "center" | "right" | "justify";
    children?: Snippet;
  }

  let {
    position = "justify",
    children
  }: Props = $props();
</script>

<style lang="scss">
  @import "../../styles/dimensions.scss";

  div {
    display: flex;
    width: 100%;
    flex-direction: row;
    flex-wrap: nowrap;
    gap: $gapSmall;
  }

  div.justify {
    justify-content: space-between;
    align-items: center;
    justify-items: center;
  }

  div.left {
    justify-content: flex-start;
    align-items: start;
  }

  div.right {
    justify-content: flex-end;
    align-items: end;
  }

  div.center {
    justify-content: space-around;
    align-items: center;
    justify-items: center;
  }
</style>

<div
  class:justify={position == "justify"}
  class:left={position == "left"}
  class:right={position == "right"}
  class:center={position == "center"}
>
  {@render children?.()}
</div>
<script lang="ts">
  import type { Snippet } from "svelte";

  interface Props {
    position?: "left" | "center" | "right" | "justify";
    width?: "full" | "auto";
    children?: Snippet;
  }

  let {
    position = "justify",
    width = "full",
    children
  }: Props = $props();
</script>

<style lang="scss">
  @use "../../styles/dimensions.scss";

  div {
    display: flex;
    flex-direction: row;
    flex-wrap: nowrap;
    gap: dimensions.$gapSmall;
  }

  div.full {
    width: 100%;
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
  class:full={width == "full"}
>
  {@render children?.()}
</div>
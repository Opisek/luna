<script lang="ts">
  import type { Snippet } from "svelte";

  interface Props {
    position?: "top" | "center" | "bottom" | "justify";
    height?: "full" | "auto";
    children?: Snippet;
  }

  let {
    position = "justify",
    height = "full",
    children
  }: Props = $props();
</script>

<style lang="scss">
  @use "../../styles/dimensions.scss";

  div {
    display: flex;
    flex-direction: column;
    flex-wrap: nowrap;
    gap: dimensions.$gapSmall;
  }

  div.full {
    height: 100%;
  }

  div.justify {
    justify-content: space-between;
    align-items: center;
    justify-items: center;
  }

  div.top {
    justify-content: flex-start;
    align-items: center;
  }

  div.bottom {
    justify-content: flex-end;
    align-items: center;
  }

  div.center {
    justify-content: space-around;
    align-items: center;
    justify-items: center;
  }
</style>

<div
  class:justify={position == "justify"}
  class:top={position == "top"}
  class:bottom={position == "bottom"}
  class:center={position == "center"}
  class:full={height == "full"}
>
  {@render children?.()}
</div>
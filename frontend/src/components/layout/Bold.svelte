<script lang="ts">
  import type { Snippet } from "svelte";

  interface Props {
    color?: "inherit" | "success" | "failure" | "accent";
    children?: Snippet;
  }

  let {
    color = "inherit",
    children
  }: Props = $props();
</script>

<style lang="scss">
  @use "sass:map";

  @use "../../styles/colors.scss";
  @use "../../styles/text.scss";

  b {
    font-weight: text.$fontWeightBold;
  }

  @each $key, $val in colors.$specialColors {
    b.#{$key} {
      color: map.get($val, "background");
    }
  }
</style>

<b
  class:success={color == "success"}
  class:failure={color == "failure"}
  class:accent={color == "accent"}
>
  {@render children?.()}
</b>
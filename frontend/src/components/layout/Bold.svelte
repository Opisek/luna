<script lang="ts">
  import type { Snippet } from "svelte";
  import { ColorKeys } from "../../types/colors";

  interface Props {
    color?: ColorKeys; 
    children?: Snippet;
  }

  let {
    color = ColorKeys.Inherit,
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
  class:success={color == ColorKeys.Success}
  class:danger={color == ColorKeys.Danger}
  class:accent={color == ColorKeys.Accent}
  class:neutral={color == ColorKeys.Neutral}
  class:warning={color == ColorKeys.Warning}
  class:inherit={color == ColorKeys.Inherit}
>
  {@render children?.()}
</b>
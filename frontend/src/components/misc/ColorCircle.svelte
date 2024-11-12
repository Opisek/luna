<script lang="ts">
  import { isValidColor } from "$lib/common/colors";

  interface Props {
    color: string | null;
    size?: "small" | "medium" | "fill";
    shape?: "circle" | "squircle";
  }

  let { color, size = "medium", shape = "circle" }: Props = $props();
</script>

<style lang="scss">
  @import "../../styles/decoration.scss";
  @import "../../styles/dimensions.scss";

  div {
    flex-shrink: 0;
  }

  div.circle {
    border-radius: 50%;
  }

  div.squircle {
    border-radius: $borderRadius;
  }

  div.small {
    width: 0.75em;
    height: 0.75em;
  }

  div.medium {
    width: 1.5em;
    height: 1.5em;
  }

  div.fill {
    height: 0%;
    width: 100%;
    padding-bottom: 100%;
  }

  div.none {
    background: $stripes;
  }
</style>

<div
  class="color"
  class:circle={shape === "circle"}
  class:squircle={shape === "squircle"}
  class:small={size === "small"}
  class:medium={size === "medium"}
  class:fill={size === "fill"}
  class:none={!isValidColor(color)}
  style="{isValidColor(color) ? `background-color:${color}`: ""}"
></div>
<script lang="ts">
  import type { Snippet } from "svelte";

  interface Props {
    name: string;
    ownPositioning?: boolean;
    children?: Snippet;
  }

  let {
    name,
    ownPositioning = true,
    children
  }: Props = $props();
</script>

<style lang="scss">
  @import "../../styles/colors.scss";
  @import "../../styles/dimensions.scss";
  @import "../../styles/text.scss";

  label {
    color: $foregroundFaded;
    font-size: $fontSizeSmall;
    cursor: text;
  }

  .ownPositioning {
    margin-bottom: -2 * $gapSmall;
    margin-left: calc($gapSmall * ($fontSize / $fontSizeSmall));
  }
</style>

<label for={name} tabindex="-1" class:ownPositioning={ownPositioning}>
  {@render children?.()}
</label>
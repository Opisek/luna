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
  @use "../../styles/colors.scss";
  @use "../../styles/dimensions.scss";
  @use "../../styles/text.scss";

  label {
    color: colors.$foregroundDim;
    font-size: text.$fontSizeSmall;
    cursor: text;
  }

  .ownPositioning {
    margin-bottom: calc(-#{dimensions.$gap});
    padding-left: calc(dimensions.$gapSmall * (text.$fontSize / text.$fontSizeSmall));
  }
</style>

<label for={name} tabindex="-1" class:ownPositioning={ownPositioning}>
  {@render children?.()}
</label>
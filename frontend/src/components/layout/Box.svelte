<script lang="ts">
  interface Props {
    children?: import('svelte').Snippet;
  }

  let { children }: Props = $props();
</script>

<style lang="scss">
  @use "../../styles/animations.scss";
  @use "../../styles/colors.scss";
  @use "../../styles/dimensions.scss";
  @use "../../styles/decorations.scss";

  div {
    border-radius: dimensions.$borderRadius;
    max-width: 50vw;
    min-width: 30em;
    padding: dimensions.$gap dimensions.$gapLarge dimensions.$gapLarge dimensions.$gapLarge;
    border-radius: dimensions.$borderRadius;
    display: flex;
    flex-direction: column;
    flex-wrap: nowrap;
    gap: dimensions.$gap;
    box-shadow: decorations.$boxShadow;
  }
</style>

<div>
  {@render children?.()}
</div>
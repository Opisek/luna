<script lang="ts">
  interface Props {
    children?: import('svelte').Snippet;
  }

  let { children }: Props = $props();
</script>

<style lang="scss">
  @use "$lib/styles/animations.scss";
  @use "$lib/styles/colors.scss";
  @use "$lib/styles/dimensions.scss";
  @use "$lib/styles/decorations.scss";

  div {
    border-radius: dimensions.$borderRadius;
    max-width: 50vw;
    min-width: 30rem;
    padding: dimensions.$gapLarge dimensions.$gapLarger dimensions.$gapLarger dimensions.$gapLarger;
    border-radius: dimensions.$borderRadius;
    display: flex;
    flex-direction: column;
    flex-wrap: nowrap;
    gap: dimensions.$gapLarge;
    box-shadow: decorations.$boxShadow;
  }
</style>

<div>
  {@render children?.()}
</div>
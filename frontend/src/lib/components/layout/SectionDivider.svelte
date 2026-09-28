<script lang="ts">
  interface Props {
    title: string;
  }

  let { title }: Props = $props();
</script>

<style lang="scss">
  @use "$lib/styles/colors.scss";
  @use "$lib/styles/dimensions.scss";

  span.wrapper {
    color: color-mix(in srgb, currentColor 50%, transparent);
    width: 100%;
    display: grid;
    grid-template-columns: 1fr auto 1fr;
    align-items: center;
    gap: dimensions.$gapLarge;
  }

  hr {
    width: 100%;
    color: colors.$backgroundSecondary; 
    border: solid .1rem;
    border-radius: .05rem;
  }
</style>

<span class="wrapper">
  <hr/>
  {title}
  <hr/>
</span>
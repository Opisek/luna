<script lang="ts">
  interface Props {
    title: string;
  }

  let { title }: Props = $props();
</script>

<style lang="scss">
  @use "../../styles/colors.scss";
  @use "../../styles/dimensions.scss";

  span.wrapper {
    color: color-mix(in srgb, currentColor 50%, transparent);
    background-color: inherit;
    width: 100%;
    position: relative;
    text-align: center;
  }

  span.text {
    background-color: inherit;
    padding: 0 dimensions.$gapLarge;
    z-index: 2;
    position: relative;
  }

  hr {
    width: 100%;
    color: colors.$backgroundSecondary; 
    border: solid .1em;
    border-radius: dimensions.$borderRadius;
    position: absolute;
    top: .1em;
  }
</style>

<span class="wrapper">
  <span class="text">
    {title}
  </span>
  <hr/>
</span>
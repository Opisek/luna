<script lang="ts">
</script>

<style lang="scss">
  @import "../../../styles/animations.scss";
  @import "../../../styles/colors.scss";
  @import "../../../styles/dimensions.scss";

  div {
    background-color: var(--barFocusIndicatorColor, $barFocusIndicatorColor);
    height: 100%;
    width: $barFocusIndicatorWidth;
    position: absolute;
    left: 0;
    top: 0;
    transform: translateX(-100%);
    transition: transform $animationSpeedFast linear;
  }

  :global(*:focus-within:not(.clicked)) > div, 
  :global(*.error) > div {
    transform: translateX(0);
  }

  :global(*.error:not(:focus-within:not(.clicked))) > div {
    background-color: $backgroundFailure;
  }
</style>

<div></div>
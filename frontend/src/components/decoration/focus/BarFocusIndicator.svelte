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
    transition: transform $animationSpeedFast linear, background-color 0s linear $animationSpeedFast;
  }

  :global(*:focus-within:not(.clicked)) > div, 
  :global(*.error) > div {
    transform: translateX(0);
    transition: transform $animationSpeedFast linear;
  }

  :global(*.error:not(:focus-within:not(.clicked))) > div {
    background-color: $backgroundFailure;
    transition-delay: 0;
    transition: transform $animationSpeedFast linear;
  }
</style>

<div></div>
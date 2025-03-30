<style lang="scss">
  @use "../../../styles/animations.scss";
  @use "../../../styles/colors.scss";
  @use "../../../styles/dimensions.scss";

  div {
    background-color: var(--barFocusIndicatorColor, colors.$barFocusIndicatorColor);
    height: 100%;
    width: dimensions.$barFocusIndicatorWidth;
    position: absolute;
    left: 0;
    top: 0;
    transform: translateX(-100%);
    transition: transform animations.$animationSpeedFast linear, background-color 0s linear animations.$animationSpeedFast;
    z-index: 5;
  }

  :global(*:focus-within:not(.clicked)) > div, 
  :global(*.error) > div {
    transform: translateX(0);
    transition: transform animations.$animationSpeedFast linear;
  }

  :global(*.error:not(:focus-within:not(.clicked))) > div {
    background-color: colors.$backgroundFailure;
    transition-delay: 0;
    transition: transform animations.$animationSpeedFast linear;
  }
</style>

<div></div>
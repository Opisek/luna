<style lang="scss">
  @use "../../../styles/animations.scss";
  @use "../../../styles/colors.scss";
  @use "../../../styles/dimensions.scss";

  div.wrapper {
    position: absolute;
    height: dimensions.$underlineWidth;
    width: 100%;
    bottom: -(dimensions.$underlineWidth);
    bottom: 0;
    transition: padding animations.$animationSpeedFast linear, opacity 0s linear animations.$animationSpeedFast;
    padding: 0 50%;
    opacity: 0;
  }
  
  div.bar {
    background-color: var(--underlineFocusIndicatorColor, colors.$underlineFocusIndicatorColor);
    border-radius: calc(dimensions.$underlineWidth / 2);
    width: 100%;
    height: 100%;
  }

  :global(*:focus-within:not(.clicked)) > div {
    transition-delay: 0s !important;
    padding: 0;
    opacity: 1;
  }
</style>

<div class="wrapper">
  <div class="bar">
  </div>
</div>
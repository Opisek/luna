<style lang="scss">
  @import "../../../styles/animations.scss";
  @import "../../../styles/colors.scss";
  @import "../../../styles/dimensions.scss";

  div.wrapper {
    position: absolute;
    height: $underlineWidth;
    width: 100%;
    bottom: -$underlineWidth;
    bottom: 0;
    transition: padding $animationSpeedFast linear, opacity 0s linear $animationSpeedFast;
    padding: 0 50%;
    opacity: 0;
  }
  
  div.bar {
    background-color: var(--underlineFocusIndicatorColor, $underlineFocusIndicatorColor);
    border-radius: calc($underlineWidth / 2);
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
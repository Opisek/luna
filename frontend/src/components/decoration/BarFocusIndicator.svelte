<script lang="ts">
</script>

<style lang="scss">
  @import "../../styles/animations.scss";
  @import "../../styles/colors.scss";
  @import "../../styles/dimensions.scss";

  div {
    background-color: var(--barFocusIndicatorColor, $backgroundAccent);
    height: 100%;
    width: $borderActiveWidth;
    position: absolute;
    left: 0;
    top: 0;
    transform: translateX(-100%);
    transition: transform $animationSpeedFast linear;
  }

  :global(*:focus-within:not(.clicked)) > div {
    transform: translateX(0);
  }
</style>

<div></div>
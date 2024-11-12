<style lang="scss">
  @import "../../styles/animations.scss";
  @import "../../styles/colors.scss";
  @import "../../styles/dimensions.scss";

  $dotSize: 0.3em;

  .loader {
    display: inline-flex;
    flex-direction: row;
    flex-wrap: nowrap;
    gap: calc($dotSize / 2);
  }

  .dot {
    width: $dotSize;
    height: $dotSize;
    border-radius: 50%;
    background: currentColor;
    animation: loader calc(2 * $animationSpeedSlow) $cubic infinite;
  }

  
  @for $i from 1 through 3 {
    .dot:nth-child(#{$i}) {
      animation-delay: $animationSpeedFast * ($i - 1);
    }
  }

  @keyframes loader {
    0%, 50%, 100% {
      transform: translateY(0) scale(1);
    }
    25% {
      transform: translateY(-100%) scale(1.25);
    }
    40% {
      transform: translateY(10%) scale(0.9);
    }
  }
</style>

<span class="loader">
  {#each Array(3) as _}
    <span class="dot"></span>
  {/each}
</span>
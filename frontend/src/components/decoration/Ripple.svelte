<!-- based on https://github.com/GeekLaunch/button-ripple-effect/ -->
<script lang="ts">
  import { browser } from "$app/environment";

  interface Props {
    event: MouseEvent;
    parent: HTMLElement;
  }

  let { event, parent }: Props = $props();

  let circle: HTMLDivElement;
  let disappear = $state(-1);

  $effect(() => {
    ((circle: HTMLDivElement) => {
      if (!circle) return;

      let diameter = Math.max((parent.clientWidth, parent.clientHeight));
      circle.style.width = circle.style.height = `${diameter}px`;

      let rect = parent.getBoundingClientRect();
      circle.style.left = `${event.clientX - rect.left -diameter/2}px`;
      circle.style.top = `${event.clientY - rect.top -diameter/2}px`;

      if (browser) {
        disappear = 0;
        window.addEventListener("mouseup", () => { disappear++ }, { once: true });
      } else circle.remove();
    })(circle);
  });

  function transitionEnd() {
    disappear++;
    if (disappear >= 3) circle.remove();
  }

</script>

<style lang="scss">
  @use "../../styles/animations.scss";
  @use "../../styles/colors.scss";

  div.ripple {
    border-radius: 50%;
    border-radius: 50%;
    position: absolute;
    pointer-events: none;

    animation: ripple animations.$animationSpeedVerySlow animations.$cubic forwards;

    background-color: colors.$backgroundPrimary;
    opacity: 0.5;
    transition: opacity animations.$animationSpeed;
    transform: scale(0);
  }
  div.ripple.animate {
    opacity: 0.25;
  }
  div.ripple.disappear {
    transition: opacity animations.$animationSpeedSlow !important;
    opacity: 0;
  }
</style>

<div
  class="ripple"
  class:animate={disappear >= 0}
  class:disappear={disappear >= 2}
  bind:this={circle}
  ontransitionend={transitionEnd}
></div>
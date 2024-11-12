<!-- based on https://github.com/GeekLaunch/button-ripple-effect/ -->
<script lang="ts">
  interface Props {
    event: MouseEvent;
    parent: HTMLElement;
  }

  let { event, parent }: Props = $props();

  let circle: HTMLDivElement;

  $effect(() => {
    ((circle: HTMLDivElement) => {
      if (!circle) return;

      let diameter = Math.max((parent.clientWidth, parent.clientHeight));
      circle.style.width = circle.style.height = `${diameter}px`;

      let rect = parent.getBoundingClientRect();
      circle.style.left = `${event.clientX - rect.left -diameter/2}px`;
      circle.style.top = `${event.clientY - rect.top -diameter/2}px`;

      setTimeout(() => circle.remove(), 1500);
    })(circle);
  });

</script>

<style lang="scss">
  @import "../../styles/animations.scss";
  @import "../../styles/colors.scss";

  div.ripple {
    border-radius: 50%;
    border-radius: 50%;
    position: absolute;
    pointer-events: none;

    animation: ripple $animationSpeedVerySlow $cubic;

    background-color: $backgroundPrimary;
    opacity: 0.5;
    transform: scale(0);
  }
</style>

<div
  class="ripple"
  bind:this={circle}
></div>
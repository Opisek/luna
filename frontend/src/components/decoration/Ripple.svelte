<!-- based on https://github.com/GeekLaunch/button-ripple-effect/ -->
<script lang="ts">
  export let event: MouseEvent;
  export let parent: HTMLElement;

  let circle: HTMLDivElement;

  $: ((circle: HTMLDivElement) => {
    if (!circle) return;

    let diameter = Math.max((parent.clientWidth, parent.clientHeight));
    circle.style.width = circle.style.height = `${diameter}px`;

    let rect = parent.getBoundingClientRect();
    circle.style.left = `${event.clientX - rect.left -diameter/2}px`;
    circle.style.top = `${event.clientY - rect.top -diameter/2}px`;

    setTimeout(() => circle.remove(), 500);
  })(circle);

</script>

<style lang="scss">
  @import "../../styles/animations.scss";

  div {
    border-radius: 50%;
    border-radius: 50%;
    background-color: rgba(255, 255, 255, 0.7);
    position: absolute;
    transform: scale(0);
    animation: ripple $animationSpeedVerySlow $cubic;
  }
</style>

<div
  bind:this={circle}
/>
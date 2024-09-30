<script lang="ts">
  export let up: () => void = () => {};
  export let down: () => void = () => {};
  export let click: () => void = () => {};
  export let visible: boolean = true;
</script>

<style lang="scss">
  @import "../../styles/animations.scss";
  @import "../../styles/colors.scss";
  @import "../../styles/dimensions.scss";

  button {
    all: unset;
    border-radius: 50%;
    display: flex;
    align-items: center;
    padding: $paddingTiny;
    cursor: pointer;
    position: relative;
  }

  button.hidden {
    visibility: hidden;
  }

  div.circle {
    position: absolute;
    background-color: $backgroundSecondary;
    z-index: -1;
    border-radius: 50%;
    transition: all $cubic $animationSpeed;
    left: 50%;
    top: 50%;
    width: 0%;
    height: 0%;
  }

  button:hover div.circle {
    left: 0;
    top: 0;
    width: 100%;
    height: 100%;
  }

  button:active div.circle {
    width: 125%;
    height: 125%;
    left: -12.5%;
    top: -12.5%;
  }
</style>

<button on:click={click} on:mousedown={down} on:mouseleave={up} on:mouseup={up} class:hidden={!visible}>
  <div class="circle"></div>
  <slot/>
</button>
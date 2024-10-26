<script lang="ts">
  export let up: () => void = () => {};
  export let down: () => void = () => {};
  export let click: () => void = () => {};
  export let visible: boolean = true;
  export let style: string = "";

  export let tabindex: number = 0;

  function clickInternal(e: MouseEvent) {
    e.stopPropagation();
    click();
  }
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
    transition: all $cubic $animationSpeed;
  }

  button.hidden {
    visibility: hidden;
  }

  div.circle {
    position: absolute;
    background-color: $backgroundSecondary;
    z-index: -1;
    border-radius: 50%;
    left: 50%;
    top: 50%;
    width: 0%;
    height: 0%;
    transition: all $cubic $animationSpeed;
  }

  button:hover div.circle, button:focus div.circle {
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

<button
  on:click={clickInternal}
  on:mousedown={down}
  on:mouseleave={up}
  on:mouseup={up}
  class:hidden={!visible}
  style={style}
  tabindex="{tabindex}"
>
  <div class="circle"></div>
  <slot/>
</button>
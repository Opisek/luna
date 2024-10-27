<script lang="ts">
  import { addRipple } from "$lib/client/decoration";

  export let onClick: () => void = () => {};

  // TODO: could not figure out enums for this, try again later
  export let color: string;
  export let type: "button" | "submit" = "button";

  let button: HTMLButtonElement;
</script>

<style lang="scss">
  @import "../../styles/animations.scss";
  @import "../../styles/colors.scss";
  @import "../../styles/dimensions.scss";

  button {
    // unset props
    outline: none;
    border: 0;
    margin: 0;

    cursor: pointer;
    padding: $gapSmall;
    border-radius: $borderRadius;

    min-width: 5em;
    
    position: relative;
    overflow: hidden; 

    transition: background-color $cubic $animationSpeed;
  }

  @each $key, $val in $specialColors {
    button.#{$key} {
      background-color: map-get($val, "background");
      color: map-get($val, "foreground");
    }
    button.#{$key}:hover, button.#{$key}:focus {
      background-color: map-get($val, "backgroundActive");
    }
  }
</style>

<button
  bind:this={button}
  on:click={onClick}
  on:mouseleave={button.blur}
  class:success={color == "success"}
  class:failure={color == "failure"}
  class:accent={color == "accent"}
  type={type}
>
  <slot/>
</button>
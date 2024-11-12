<script lang="ts">
  import type { Snippet } from "svelte";

  interface Props {
    onClick?: () => void;
    // TODO: could not figure out enums for this, try again later
    color: string;
    type?: "button" | "submit";
    enabled?: boolean;
    children?: Snippet;
  }

  let {
    onClick = () => {},
    color,
    type = "button",
    enabled = true,
    children
  }: Props = $props();

  let button: HTMLButtonElement = $state(new HTMLButtonElement());
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

  .disabled {
    cursor: not-allowed;
  }

  @each $key, $val in $specialColors {
    button.#{$key} {
      background-color: map-get($val, "background");
      color: map-get($val, "foreground");
    }
    button.#{$key}:hover:not(.disabled), button.#{$key}:focus:not(.disabled) {
      background-color: map-get($val, "backgroundActive");
    }
  }
</style>

<button
  bind:this={button}
  onclick={onClick}
  onmouseleave={button.blur}
  class:success={color == "success"}
  class:failure={color == "failure"}
  class:accent={color == "accent"}
  type={type}
  disabled={!enabled}
  class:disabled={!enabled}
>
  {@render children?.()}
</button>
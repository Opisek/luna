<script lang="ts">
  import type { Snippet } from "svelte";

  interface Props {
    onClick?: () => void;
    // TODO: could not figure out enums for this, try again later
    color?: string;
    type?: "button" | "submit";
    enabled?: boolean;
    href?: string;
    children?: Snippet;
  }

  let {
    onClick = () => {},
    color = "neutral",
    type = "button",
    enabled = true,
    href = "",
    children
  }: Props = $props();
</script>

<style lang="scss">
  @use "sass:map";

  @use "../../styles/animations.scss";
  @use "../../styles/colors.scss";
  @use "../../styles/dimensions.scss";
  @use "../../styles/text.scss";

  button, a {
    // unset props
    background: none;
    color: inherit;
    border: none;
    padding: 0;
    font: inherit;
    cursor: pointer;
    outline: inherit;
    text-decoration: none;

    display: inline;

    cursor: pointer;
    padding: dimensions.$gapSmall;
    border-radius: dimensions.$borderRadius;

    min-width: 5em;
    text-align: center;
    
    position: relative;
    overflow: hidden; 

    transition: background-color animations.$cubic animations.$animationSpeed;
  }

  .disabled {
    cursor: not-allowed;
  }

  @each $key, $val in colors.$specialColors {
    button.#{$key}, a.#{$key} {
      background-color: map.get($val, "background");
      color: map.get($val, "foreground");
    }
    button.#{$key}:hover:not(.disabled), button.#{$key}:focus:not(.disabled),
    a.#{$key}:hover:not(.disabled), a.#{$key}:focus:not(.disabled) {
      background-color: map.get($val, "backgroundActive");
    }
  }
</style>

{#if href !== ""}
  <a
    class:success={color == "success"}
    class:failure={color == "failure"}
    class:accent={color == "accent"}
    class:neutral={color == "neutral"}
    onmouseleave={(e) => {(e.target as HTMLButtonElement).blur()}}
    class:disabled={!enabled}
    href={enabled ? href : "#"}
  >
    {@render children?.()}
  </a>
{:else}
  <button
    class:success={color == "success"}
    class:failure={color == "failure"}
    class:accent={color == "accent"}
    class:neutral={color == "neutral"}
    onclick={onClick}
    onmouseleave={(e) => {(e.target as HTMLButtonElement).blur()}}
    type={type}
    disabled={!enabled}
    class:disabled={!enabled}
  >
    {@render children?.()}
  </button>
{/if}
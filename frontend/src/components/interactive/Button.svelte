<script lang="ts">
  import type { Snippet } from "svelte";
  import { ColorKeys } from "../../types/colors";
  import { addRipple, focusIndicator } from "../../lib/client/decoration";
  import Loader from "../decoration/Loader.svelte";

  interface Props {
    onClick?: () => any;
    color?: ColorKeys;
    type?: "button" | "submit";
    compact?: boolean;
    enabled?: boolean;
    href?: string;
    children?: Snippet;
  }

  let {
    onClick = () => {},
    color = ColorKeys.Neutral,
    type = "button",
    compact = false,
    enabled = true,
    href = "",
    children
  }: Props = $props();

  let loading = $state(false);
  async function clickHandler() {
    if (loading) return;
    const result = onClick();
    if (!(result instanceof Promise)) return;
    loading = true;
    await result;
    loading = false;
  }
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

    min-width: dimensions.$buttonMinWidth;
    text-align: center;
    
    position: relative;
    overflow: hidden; 
  }

  // Nasty hack to get <Spinner/> to have the same height as text inside of buttons.
  button > :global(span.spinner) {
    &::before, &::after {
      content: "a";
      visibility: hidden;
    }
  }

  button:not(.neutral) {
    --barFocusIndicatorColor: #{colors.$barFocusIndicatorColorAlt};
  }

  button.compact, a.compact {
    min-width: dimensions.$buttonMinWidthCompact;
  }

  .disabled {
    cursor: not-allowed;
  }

  @each $key, $val in colors.$specialColors {
    button.#{$key}, a.#{$key} {
      background-color: map.get($val, "background");
      color: map.get($val, "foreground");
    }
    button.#{$key}.disabled, a.#{$key}.disabled {
      color: color-mix(in srgb, map.get($val, "foreground") 50%, transparent);
    }
  }
</style>

{#if href !== ""}
  <a
    class:success={color == ColorKeys.Success}
    class:warning={color == ColorKeys.Warning}
    class:danger={color == ColorKeys.Danger}
    class:accent={color == ColorKeys.Accent}
    class:neutral={color == ColorKeys.Neutral}
    class:inherit={color == ColorKeys.Inherit}
    class:compact={compact}
    onmouseleave={(e) => {(e.target as HTMLButtonElement).blur()}}
    class:disabled={!enabled}
    href={enabled ? href : "#"}
    onmousedown={addRipple}
    use:focusIndicator
  >
    {@render children?.()}
  </a>
{:else}
  <button
    class:success={color == ColorKeys.Success}
    class:warning={color == ColorKeys.Warning}
    class:danger={color == ColorKeys.Danger}
    class:accent={color == ColorKeys.Accent}
    class:neutral={color == ColorKeys.Neutral}
    class:inherit={color == ColorKeys.Inherit}
    class:compact={compact}
    onclick={clickHandler}
    onmouseleave={(e) => {(e.target as HTMLButtonElement).blur()}}
    type={type}
    disabled={!enabled}
    class:disabled={!enabled}
    onmousedown={addRipple}
    use:focusIndicator
  >
    {#if loading}
      <Loader/>
    {:else}
      {@render children?.()}
    {/if}
  </button>
{/if}
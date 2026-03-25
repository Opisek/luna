<script lang="ts">
  import type { Snippet } from "svelte";

  import { NoOp } from "$lib/client/placeholders";
  import Tooltip from "./Tooltip.svelte";
  import { ColorKeys } from "../../types/colors";
  import Button from "./Button.svelte";
  import { getSettings } from "../../lib/client/data/settings.svelte";
  import { UserSettingKeys } from "../../types/settings";
  import Spinner from "../decoration/Spinner.svelte";
  import Popup from "../popups/Popup.svelte";

  interface Props {
    alt: string;
    up?: () => void;
    down?: () => void;
    onClick?: () => any;
    visible?: boolean;
    style?: string;
    tabindex?: number;
    href?: string;
    type?: "button" | "submit";
    enabled?: boolean;
    color?: ColorKeys;
    canRenderAsButton?: boolean;
    children?: Snippet;
  }

  let {
    alt,
    up = NoOp,
    down = NoOp,
    onClick = NoOp,
    visible = true,
    style = "",
    tabindex = 0,
    href = "",
    type = "button",
    enabled = true,
    color = ColorKeys.Neutral,
    canRenderAsButton = false,
    children
  }: Props = $props();

  const settings = getSettings();

  let button = $state<HTMLElement | undefined>();

  let showPopover = $state(NoOp);
  let hidePopover = $state(NoOp);

  let loading = $state(false);
  async function clickInternal(e: MouseEvent) {
    e.stopPropagation();
    if (loading) return;
    const result = onClick();
    if (!(result instanceof Promise)) return;
    loading = true;
    await result;
    loading = false;
  }

  function leaveInternal() {
    hidePopover();
    if (!button) return;
    button.blur();
    up();
  }
  function upInternal() {
    if (!button) return;
    button.blur();
    up();
  }
</script>

<style lang="scss">
  @use "../../styles/animations.scss";
  @use "../../styles/colors.scss";
  @use "../../styles/dimensions.scss";

  button, a {
    all: unset;
    border-radius: 50%;
    display: flex;
    align-items: center;
    padding: dimensions.$gapSmaller;
    cursor: pointer;
    position: relative;
    transition: all animations.$cubic animations.$animationSpeed;
  }

  button.hidden, a.hidden, .loading > :global(svg) {
    visibility: hidden;
  }

  button,
  a {
    &::before {
      content: "";
      position: absolute;
      border-radius: 50%;
      left: 50%;
      top: 50%;
      width: 0%;
      height: 0%;
      transition: all animations.$cubic animations.$animationSpeed;
      pointer-events: none;
    }
  }

  button:hover,
  button:focus,
  button.loading,
  a:hover,
  a:focus {
    &::before {
      left: 0;
      top: 0;
      width: 100%;
      height: 100%;
    }
  }

  button:active:not(.loading),
  a:active {
    &::before {
      width: 125%;
      height: 125%;
      left: -12.5%;
      top: -12.5%;
    }
  }

  button:hover,
  button:focus,
  button.loading,
  a:hover,
  a:focus {
    &.neutral {
      color: colors.$foregroundSecondary;
      &::before {
        background-color: colors.$backgroundSecondary;
      }
    } 
    &.success {
      color: colors.$foregroundSuccess;
      &::before {
        background-color: colors.$backgroundSuccess;
      }
    } 
    &.accent {
      color: colors.$foregroundAccent;
      &::before {
        background-color: colors.$backgroundAccent;
      }
    } 
    &.warning {
      color: colors.$foregroundWarning;
      &::before {
        background-color: colors.$backgroundWarning;
      }
    } 
    &.danger {
      color: colors.$foregroundFailure;
      &::before {
        background-color: colors.$backgroundFailure;
      }
    } 
    &.inherit {
      color: inherit;
      &::before {
        background-color: inherit;
      }
    } 
  }

  .disabled {
    cursor: not-allowed;
    opacity: 0.5;
  }

  .loading > :global(:nth-child(2)) {
    position: absolute;
    width: 100%;
    height: 100%;
    left: 0;
    top: 0;
    display: flex;
    align-items: center;
    justify-content: center;
  }
</style>

{#if canRenderAsButton && settings.userSettings[UserSettingKeys.UseTextButtons]}
  <Button
     onClick={onClick}
     color={color}
     type={type}
     enabled={enabled}
     href={href}
  >{alt}</Button>
{:else}
  {@render buttonSnippet()}
{/if}

{#snippet buttonSnippet()}
  {#if href !== ""}
    <a
      bind:this={button}
      class:hidden={!visible}
      href={enabled ? href : "#"}
      style={style}
      tabindex="{tabindex}"
      type={type}
      class:disabled={!enabled}
      class:success={color == ColorKeys.Success}
      class:accent={color == ColorKeys.Accent}
      class:warning={color == ColorKeys.Warning}
      class:danger={color == ColorKeys.Danger}
      class:neutral={color == ColorKeys.Neutral}
      class:inherit={color == ColorKeys.Inherit}
    >
      {@render children?.()}
    </a>
  {:else}
    <button
      bind:this={button}
      onclick={clickInternal}
      onmousedown={down}
      onmouseleave={leaveInternal}
      onmouseenter={showPopover}
      onmouseup={upInternal}
      class:hidden={!visible}
      type={type}
      style={style}
      tabindex="{tabindex}"
      class:disabled={!enabled}
      class:success={color == ColorKeys.Success}
      class:accent={color == ColorKeys.Accent}
      class:warning={color == ColorKeys.Warning}
      class:danger={color == ColorKeys.Danger}
      class:neutral={color == ColorKeys.Neutral}
      class:inherit={color == ColorKeys.Inherit}
      disabled={!enabled}
      class:loading
      aria-label={alt}
    >
      {@render children?.()}
      {#if loading}
        <Spinner/>
      {/if}
      {#if alt != ""}
        <Popup bind:showPopup={showPopover} bind:hidePopup={hidePopover} delayed={true}>
          {alt}
        </Popup>
      {/if}
    </button>
  {/if}
{/snippet}
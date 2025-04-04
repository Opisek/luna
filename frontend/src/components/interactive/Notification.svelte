<script lang="ts">
  import { ColorKeys } from "../../types/colors";
  import type { NotificationModel } from "../../types/notification";

  interface Props {
    notification: NotificationModel;
    height: number;
    shift: number;
  }

  let {
    notification,
    height = $bindable(),
    shift
  }: Props = $props();

  let isNew = $state(true);
  setTimeout(() => {
    isNew = false;
  }, 0);

  let wrapper: HTMLElement;

  function removeGracefully() {
    if (wrapper.matches(":hover")) {
      wrapper.addEventListener("mouseleave", () => {
        notification.remove();
      }, { once: true });
    } else {
      notification.remove();
    }
  }

  let viewDetails = $state(false);
  function showDetails() {
    viewDetails = true;
  }
  function hideDetails() {
    viewDetails = false;
  }

  let previousCount = $state(notification.count);
  let counterVisible = $state(false);
  let counterPop = $state(true);
  $effect(() => {
    if (notification.count == 1 || previousCount == notification.count) return;
    previousCount = notification.count;
    counterPop = false;
    setTimeout(() => {
      counterPop = true;
      counterVisible = true;
    }, 10);
  });
</script>

<style lang="scss">
  @use "sass:map";

  @use "../../styles/animations.scss";
  @use "../../styles/colors.scss";
  @use "../../styles/dimensions.scss";
  @use "../../styles/text.scss";

  div.wrapper {
    bottom: 100%;
    right: 0;
    width: 100%;
    position: absolute;
    padding-top: dimensions.$gapSmaller;
    transition: all animations.$cubic animations.$animationSpeedSlow; 
    z-index: 50; // remove if we switch to popover
  }

  div.box {
    padding: dimensions.$gapLarge;
    border-radius: dimensions.$borderRadius;
    cursor: pointer;
    position: relative;
    overflow: hidden;
    white-space: pre-wrap;
  }

  @each $key, $val in colors.$specialColors {
    .#{$key} {
      background-color: map.get($val, "background");
      color: map.get($val, "foreground");
    }
    .#{$key} .details {
      color: color-mix(in srgb, map.get($val, "foreground"), 50%, transparent);
    }
  }

  .disappear {
    opacity: 0;
  }

  div.timer {
    position: absolute;
    bottom: 0;
    left: 0;
    background-color: black;
    opacity: 0.2;
    content: "";
    width: 100%;
    height: 0.5em;
    animation: notification-timer var(--notificationExpireTime) linear forwards;
  }

  .details {
    width: 100%;
    font-size: text.$fontSizeSmall;
    cursor: pointer;
    display: inline-block;
  }

  .count {
    position: absolute;
    top: dimensions.$gapSmall;
    right: dimensions.$gapSmall;
    background-color: rgba(0, 0, 0, 0.2);
    border-radius: 50%;
    height: 1.4em;
    aspect-ratio: 1/1;
    font-size: text.$fontSizeSmall;
    display: grid;
    text-align: center;
    align-items: center;
    align-content: center;
    opacity: 0;
    transition: opacity animations.$animationSpeedFast;
  }
  .opaque {
    opacity: 1;
  }
  .pop {
    animation: pop animations.$cubic animations.$animationSpeed;
  }
</style>

<div
  class="wrapper"
  class:disappear={notification.disappear}
  style="transform: translateY({isNew ? "100%" : `${shift}px`});"
  bind:clientHeight={height}
  bind:this={wrapper}
>
<!-- svelte-ignore a11y_no_static_element_interactions -->
  <div
    class="box"
    class:success={notification.color === ColorKeys.Success}
    class:warning={notification.color === ColorKeys.Warning}
    class:danger={notification.color === ColorKeys.Danger}
    class:accent={notification.color === ColorKeys.Accent}
    class:neutral={notification.color === ColorKeys.Neutral}
    class:inherit={notification.color === ColorKeys.Inherit}
    onclick={notification.remove}
    onkeypress={notification.remove}
  >
    {notification.message}

    {#if notification.details}
      <span class="details" onmouseenter="{showDetails}" onmouseleave="{hideDetails}">
        {#if viewDetails}
          {notification.details}
        {:else}
          Hover to view details
        {/if}
      </span>
    {/if}

    {#if notification.count > 1}
      <div class="count" class:opaque={counterVisible} class:pop={counterPop}>
        {notification.count}
      </div>
    {/if}

    {#if counterPop}
      <div class="timer" onanimationend={removeGracefully}></div>
    {/if}
  </div>
</div>
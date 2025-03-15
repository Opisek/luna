<script lang="ts">
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
</script>

<style lang="scss">
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

  .success {
    background-color: colors.$backgroundSuccess;
    color: colors.$foregroundSuccess;
  }

  .failure {
    background-color: colors.$backgroundFailure;
    color: colors.$foregroundFailure;
  }

  .info {
    background-color: colors.$backgroundAccent;
    color: colors.$foregroundAccent;
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

  .success .details {
    color: colors.$foregroundSuccessFaded;
  }
  
  .failure .details {
    color: colors.$foregroundFailureFaded;
  }

  .info .details {
    color: colors.$foregroundAccentFaded;
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
    class:success={notification.type === "success"}
    class:failure={notification.type === "failure"}
    class:info={notification.type === "info"}
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
    <div class="timer" onanimationend={removeGracefully}></div>
  </div>
</div>
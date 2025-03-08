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
    popover.showPopover();
    isNew = false;
  }, 0);

  let popover: HTMLElement;
</script>

<style lang="scss">
  @import "../../styles/animations.scss";
  @import "../../styles/colors.scss";
  @import "../../styles/dimensions.scss";

  div.wrapper {
    background-color: transparent;
    outline: none;
    border: none;
    top: calc(100% - $gapSmaller);
    left: calc(100% - $notificationWidth - $gapSmaller);
    width: $notificationWidth;
    // previously (before popover) instead of the above: bottom: 0, right: 0, no -50% subtraction in the style attribute down below
    position: absolute;
    padding-top: $gapSmaller;
    transition: all $cubic $animationSpeedSlow; 
  }

  div.box {
    padding: $gap;
    border-radius: $borderRadius;
    cursor: pointer;
    position: relative;
    overflow: hidden;
    white-space: pre-wrap;
  }

  .success {
    background-color: $backgroundSuccess;
    color: $foregroundSuccess;
  }

  .failure {
    background-color: $backgroundFailure;
    color: $foregroundFailure;
  }

  .info {
    background-color: $backgroundAccent;
    color: $foregroundAccent;
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
</style>

<div
  class="wrapper"
  class:disappear={notification.disappear}
  style="transform: translateY({isNew ? "100%" : `calc(${shift}px - 50%)`});"
  popover="manual"
  bind:clientHeight={height}
  bind:this={popover}
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
    <div class="timer"></div>
  </div>
</div>
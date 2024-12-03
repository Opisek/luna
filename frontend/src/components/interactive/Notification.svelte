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
</script>

<style lang="scss">
  @import "../../styles/animations.scss";
  @import "../../styles/colors.scss";
  @import "../../styles/dimensions.scss";

  div.wrapper {
    bottom: 100%;
    right: 0;
    width: 100%;
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
  style="transform: translateY({isNew ? "100%" : `${shift}px`});"
  bind:clientHeight={height}
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
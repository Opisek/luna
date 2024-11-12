<script lang="ts">
  interface Props {
    notification: NotificationModel;
    shift: number;
  }

  let {
    notification,
    shift
  }: Props = $props();
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
    padding-top: $paddingTiny;
    transition: all $cubic $animationSpeedSlow; 
  }

  div.box {
    padding: $paddingSmall;
    border-radius: $borderRadius;
    cursor: pointer;
    position: relative;
    overflow: hidden;
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
    animation: notification-timer $notificationExpireTime linear forwards;
  }
</style>

<div
  class="wrapper"
  class:disappear={notification.disappear}
  style="transform: translateY({shift * -100}%);"
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
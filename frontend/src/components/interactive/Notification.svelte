<script lang="ts">
  import { browser } from "$app/environment";

  export let notification: NotificationModel;
  export let shift: number;

  let surpressAnimation: boolean;

  let actualShift = 0;
  $: ((requestedShift: number) => {
    actualShift = requestedShift;
  })(shift);

  let lastNotification: NotificationModel | null = null;
  $: ((newNotification: NotificationModel) => {
    if (lastNotification != null && lastNotification.created.getTime() < newNotification.created.getTime()) {
      // TODO: i hate this surpress animation solution, because in certain edge cases, it will still look off.
      // TODO: find a better way to do the exit animation giving svelte's constraints.
      surpressAnimation = true;
      setTimeout(() => {
        surpressAnimation = false;
      }, 10);
    }
    lastNotification = newNotification;
  })(notification);
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

  .surpress {
    transition: none !important;
  }
</style>

<div
  class="wrapper"
  class:disappear={notification.disappear}
  style="transform: translateY({actualShift * -100}%);"
  class:surpress={surpressAnimation}
>
<!-- svelte-ignore a11y-no-static-element-interactions -->
  <div
    class="box"
    class:success={notification.type === "success"}
    class:failure={notification.type === "failure"}
    class:info={notification.type === "info"}
    on:click={notification.remove}
    on:keypress={notification.remove}
  >
    {notification.message}
  </div>
</div>
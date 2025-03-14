<script lang="ts">
  import Notification from "../components/interactive/Notification.svelte";

  import { notificationExpireTime, notifications, queueNotification } from "$lib/client/notifications";

  interface Props {
    children?: import('svelte').Snippet;
  }

  let { children }: Props = $props();

  let notifsWrapper: HTMLDivElement;

  // We store the height of every notification element so we can calculate the
  // proper Y-offsets for each notification.
  let notifsHeights: number[] = $state([]);

  // Here, we calculate the Y-offsets for each notification.
  // The older the notification, the higher up it should be.
  let notifsShifts = $derived((() => {
    let shifts = [];
    let shift = 0;
    for (let i = $notifications.length - 1; i >= 0; i--) {
      shifts.unshift(shift);
      shift -= notifsHeights[i];
    }
    return shifts;
  })());

  // When a notification is removed, we can't wait for the height binding to
  // update - the notifications will jump around. For this reason, we manually
  // update the heights that we know have changed. In particular, if
  // notification `i` disappears, then all notifications `j, j > i` will have
  // their index shifted by `-1' - and that's where their new height entries
  // will be expected.
  let prevNotifs: NotificationModel[] = [];
  notifications.subscribe((notifs) => {
    let skip = 0;
    for (const [i, oldNotif] of prevNotifs.entries()) {
      if (!notifs.includes(oldNotif)) skip++;
      notifsHeights[i] = notifsHeights[i + skip];
    }
    prevNotifs = notifs;

    while(notifsHeights.length < $notifications.length) {
      notifsHeights.push(0);
    }

    if (notifsWrapper) {
      notifsWrapper.hidePopover();
      notifsWrapper.showPopover();
    }
  });
</script>

<style lang="scss">
  @use "../styles/dimensions.scss";
  @use "../styles/text.scss";

  :global(*) {
    box-sizing: border-box;
  }

  :global(body) {
    margin: 0;
    padding: dimensions.$gap;
    gap: dimensions.$gapSmall;

    height: 100vh;
    width: 100vw;

    font-family: text.$fontFamilyText;
    font-size: text.$fontSize;
  }

  div.notifications {
    position: fixed;
    left: calc(100vw - 15em - dimensions.$gapSmall);
    top: calc(100vh - dimensions.$gapSmaller);
    width: 15em;
    height: 0;
    background-color: transparent;
    outline: 0;
    border: 0;
    overflow: visible;
  }
</style>

{@render children?.()}

<div
  bind:this={notifsWrapper}
  class="notifications"
  style="--notificationExpireTime: {notificationExpireTime}ms"
  popover="manual"
>
  {#each $notifications as notification, i (notification.created.getTime())}
    <Notification
      notification={notification}
      bind:height={notifsHeights[i]}
      shift={notifsShifts[i]}
    />
  {/each}
</div>

<!--<Footer/>-->
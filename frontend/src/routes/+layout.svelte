<script lang="ts">
  import Notification from "../components/interactive/Notification.svelte";

  import { notificationExpireTime, notifications } from "$lib/client/notifications";

  interface Props {
    children?: import('svelte').Snippet;
  }

  let { children }: Props = $props();

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
  });

</script>

<style lang="scss">
  @import "../styles/dimensions.scss";
  @import "../styles/text.scss";

  :global(*) {
    box-sizing: border-box;
  }

  :global(body) {
    margin: 0;
    padding: 0;
    height: 100vh;
    font-family: $fontFamilyText;
    font-size: $fontSize;
  }

  div.notifications {
    position: fixed;
    right: 0;
    bottom: 0;
    margin: $gapSmaller;
    width: 15em;
  }
</style>

{@render children?.()}

<div class="notifications" style="--notificationExpireTime: {notificationExpireTime}ms">
  {#each $notifications as notification, i (notification.created.getTime())}
    <Notification
      notification={notification}
      bind:height={notifsHeights[i]}
      shift={notifsShifts[i]}
    />
  {/each}
</div>
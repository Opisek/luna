<script lang="ts">
  import { notifications, notificationCount } from "$lib/client/notifications";
  import Notification from "../components/interactive/Notification.svelte";

  let notifs: NotificationModel[] = [];
  let notifsCount = 0;

  notifications.subscribe((active) => {
    notifs = active;
  });
  notificationCount.subscribe((count) => {
    notifsCount = count;
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
    margin: $paddingTiny;
    width: 15em;
  }
</style>

<slot/>

<div class="notifications">
  {#each notifs as notification, i}
    <Notification
      notification={notification}
      shift={notifsCount - i - 1}
    />
  {/each}
</div>
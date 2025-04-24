<script lang="ts">
  import Notification from "../components/interactive/Notification.svelte";

  import { notificationExpireTime, notifications } from "$lib/client/notifications";
  import { getSettings } from "$lib/client/settings.svelte";
  import { UserSettingKeys, type GlobalSettings, type UserData, type UserSettings } from "../types/settings";
  import { getTheme } from "$lib/client/theme.svelte";
  import type { NotificationModel } from "../types/notification";

  interface PageProps {
    children?: import('svelte').Snippet;
    data: {
      userData: UserData;
      userSettings: UserSettings;
      globalSettings: GlobalSettings;
    }
  }

  let {
    children,
    data
  }: PageProps = $props();

  const settings = getSettings(data);
  const theme = getTheme();

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

  $effect(() => {
    const root = document.documentElement;
    root.setAttribute("data-theme", theme.isLightMode() ? "light" : "dark");
    root.setAttribute("data-frost", settings.userSettings[UserSettingKeys.AppearenceFrostedGlass] ? "true" : "false");
  })

  // Prevent "flashing" by unloading the previous theme/font only after loading the next one.
  let currentThemeLight = $derived(settings.userSettings[UserSettingKeys.ThemeLight]);
  let currentThemeDark = $derived(settings.userSettings[UserSettingKeys.ThemeDark]);
  let previousThemeLight = $state("");
  let previousThemeDark = $state("");
  $effect(() => { setTimeout((val) => previousThemeLight = val, 100, currentThemeLight); });
  $effect(() => { setTimeout((val) => previousThemeDark = val, 100, currentThemeDark); });

  // Font loading
  let currentFontText = $derived(settings.userSettings[UserSettingKeys.FontText]);
  let currentFontTime = $derived(settings.userSettings[UserSettingKeys.FontTime]);
  let currentFontTextName = $derived(currentFontText.split("-").map(x => x.charAt(0).toUpperCase() + x.slice(1)).join(" "));
  let currentFontTimeName = $derived(currentFontTime.split("-").map(x => x.charAt(0).toUpperCase() + x.slice(1)).join(" "));
</script>

<style lang="scss">
  @use "../styles/colors.scss";
  @use "../styles/dimensions.scss";
  @use "../styles/text.scss";
  
  :global(*) {
    box-sizing: border-box;
  }

  :global(body) {
    margin: 0;
    padding: dimensions.$gapLarge;
    gap: dimensions.$gapSmall;

    height: 100vh;
    width: 100vw;

    font-family: text.$fontFamilyText;
    font-size: text.$fontSize;

    background-color: colors.$backgroundPrimary;
    color: colors.$foregroundPrimary;
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

  :global(.lucide-icon, svg:has(title)) {
    scale: var(--uiScaling, 1);
  }
</style>

<svelte:head>
  {#if previousThemeLight != "" && previousThemeLight != currentThemeLight}
    <link rel="stylesheet" href="/themes/light/{previousThemeLight}.css">
  {/if}
  {#if previousThemeDark != "" && previousThemeDark != currentThemeDark}
    <link rel="stylesheet" href="/themes/dark/{previousThemeDark}.css">
  {/if}

  <link rel="stylesheet" href="/themes/light/{currentThemeLight}.css">
  <link rel="stylesheet" href="/themes/dark/{currentThemeDark}.css">

  {@html `
    <style>
      @font-face {
        font-family: "${currentFontTextName}"; 
        src: url("/fonts/${currentFontText}.ttf");
      }
      @font-face {
        font-family: "${currentFontTimeName}"; 
        src: url("/fonts/${currentFontTime}.ttf");
      }

      :root {
        --uiScaling: ${settings.userSettings[UserSettingKeys.UiScaling]};

        ${settings.userSettings[UserSettingKeys.DisplayRoundedCorners] ? "" : "\
          --borderRadiusSmall: 0;\
          --borderRadius: 0;\
          --borderRadiusLarge: 0;\
        "}

        --fontFamilyText: ${currentFontTextName};
        --fontFamilyTime: ${currentFontTimeName};
      }
    </style> 
  `}
</svelte:head>

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
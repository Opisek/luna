<script lang="ts">
  import type { Settings } from "../../../lib/client/data/settings.svelte";
  import type { ActiveSessions } from "../../../lib/client/data/sessions.svelte";
  import { UserSettingKeys } from "../../../types/settings";
  import List from "../../forms/List.svelte";
  import ToggleInput from "../../forms/ToggleInput.svelte";
  import Button from "../../interactive/Button.svelte";
  import { ColorKeys } from "../../../types/colors";
  import { UAParser } from "ua-parser-js";
  import { Bot, Gamepad2, Info, Laptop, LogOut, Microchip, Pencil, RectangleGoggles, Smartphone, Tablet, TvMinimal, Watch } from "lucide-svelte";
  import IconButton from "../../interactive/IconButton.svelte";
  import { queueNotification } from "../../../lib/client/notifications";
  import { t } from "@sveltia/i18n";

  interface Props {
    settings: Settings;
    sessions: ActiveSessions;
    today: Date;
    showSessionModal: (initial?: Session, edit?: boolean) => Promise<Session>;
  }

  let {
    settings,
    sessions,
    today,
    showSessionModal,
  }: Props = $props();

  function deauthorizeSession(id: string) {
    sessions.deauthorizeSession(id);
  }
</script>

<style lang="scss">
  @use "../../../styles/dimensions.scss";
  @use "../../../styles/colors.scss";
  @use "../../../styles/text.scss";

  .session {
    padding: dimensions.$gapMiddle;
    background-color: colors.$backgroundSecondary;
    color: colors.$foregroundSecondary;
    border-radius: dimensions.$borderRadius;

    display: grid;
    gap: dimensions.$gapSmall;
    row-gap: 0;
    grid-template-columns: auto 1fr auto;
    grid-template-rows: auto auto;
    grid-template-areas: "device agent buttons" "device details buttons";
    justify-content: center;
    align-items: center;
  }

  .session.showId {
    grid-template-rows: auto auto auto;
    grid-template-areas: "device agent buttons" "device details buttons" "device id buttons";
  }

  .session > .device {
    grid-area: device;
    display: flex;
  }
  .session > .agent {
    grid-area: agent;
  }
  .session > .details {
    grid-area: details;
    font-size: text.$fontSizeSmall;
  }
  .session > .buttons {
    grid-area: buttons;
    display: flex;
    flex-direction: row;
    gap: dimensions.$gapSmall;
  }
  .session > .id {
    grid-area: id;
    font-size: text.$fontSizeSmall;
  }
  .session:not(.showId) > .id {
    display: none;
  }

  .session.active {
    background-color: colors.$backgroundAccent;
    color: colors.$foregroundAccent;
  }
</style>

<ToggleInput
  name={UserSettingKeys.DebugMode}
  description={t("settings.dev.ids")}
  bind:value={settings.userSettings[UserSettingKeys.DebugMode]}
/>

<Button color={ColorKeys.Accent} onClick={() => showSessionModal().catch(err => { if (err) queueNotification(ColorKeys.Danger, err.message); } )}>{t("settings.dev.api.new")}</Button>

<svelte:boundary>
  {@const apiSessions = sessions.activeSessions.filter(x => x.is_api)}
  {#if apiSessions.length !== 0}
    <List
      label={t("settings.dev.api.list")}
      items={apiSessions}
      id={item => item.id}
      template={sessionTemplate}
    />
  {/if}
</svelte:boundary>

<!-- TODO: reduce code duplication by putting templates in separate files as well -->
{#snippet sessionTemplate(s: Session)}
  {@const userAgent=UAParser(s.is_api ? "" : s.user_agent)}
  {@const deviceName=`${userAgent.os.name || ""} ${userAgent.browser.name || ""}`.trim()}
  {@const isActive=s.id === sessions.currentSession}

  <div class="session" class:active={isActive} class:showId={settings.userSettings[UserSettingKeys.DebugMode]}>
    <div class="device">
      {#if s.is_api}
        <Bot size={20}/>
      {:else if userAgent.device.type === UAParser.DEVICE.CONSOLE}
        <Gamepad2 size={20}/>
      {:else if userAgent.device.type === UAParser.DEVICE.EMBEDDED}
        <Microchip size={20}/>
      {:else if userAgent.device.type === UAParser.DEVICE.MOBILE}
        <Smartphone size={20}/>
      {:else if userAgent.device.type === UAParser.DEVICE.SMARTTV}
        <TvMinimal size={20}/>
      {:else if userAgent.device.type === UAParser.DEVICE.TABLET}
        <Tablet size={20}/>
      {:else if userAgent.device.type === UAParser.DEVICE.WEARABLE}
        <Watch size={20}/>
      <!--{:else if userAgent.device.type === UAParser.DEVICE.XR}-->
      {:else if userAgent.device.type === "xr"}
        <RectangleGoggles size={20}/>
      {:else if deviceName === ""}
        <Bot size={20}/>
      {:else}
        <Laptop size={20}/>
      {/if}
    </div>

    <span class="agent">
      {#if deviceName === ""}
        {s.user_agent}
      {:else}
        {deviceName}
      {/if}
    </span>

    <span class="details">
      {s.location}
      •
      {#if isActive}
        {t("session.current")}
      {:else if today.getDate() == s.last_seen.getDate() && today.getMonth() == s.last_seen.getMonth() && today.getFullYear() == s.last_seen.getFullYear()}
        {t("session.active.today", { values: { date: s.last_seen } })}
      {:else}
        {t("session.active.elsewhen", { values: { date: s.last_seen } })}
      {/if}
    </span>

    <div class="buttons">
      <IconButton onClick={async () => showSessionModal(s, s.is_api)} color={ColorKeys.Accent} alt={t("button.edit")}>
        <Pencil size={20}/>
      </IconButton>
      <IconButton onClick={async () => deauthorizeSession(s.id)} color={ColorKeys.Danger} alt={t("settings.dev.deauthorize")}>
        <LogOut size={20}/>
      </IconButton>
    </div>

    <span class="id">
      {t("session.id.inline", { values: { id: s.id } })}
    </span>
  </div>
{/snippet}


import { GlobalSettingKeys, UserSettingKeys, type GlobalSettings, type UserData, type UserSettings } from "../../types/settings";
import { fetchJson } from "./net";
import { queueNotification } from "./notifications";

class Settings {
  public userData: UserData = $state({
    id: "",
    username: "",
    email: "",
    admin: false,
    searchable: true,
    profile_picture: ""
  });
  public userSettings: UserSettings = $state({
    [UserSettingKeys.DebugMode]: false,
    [UserSettingKeys.DisplayAllDayEventsFilled]: true,
    [UserSettingKeys.DisplayNonAllDayEventsFilled]: false,
    [UserSettingKeys.DisplayRoundedCorners]: true,
    [UserSettingKeys.DisplaySmallCalendar]: true,
    [UserSettingKeys.DisplayWeekNumbers]: false,
    [UserSettingKeys.DynamicCalendarRows]: true,
    [UserSettingKeys.DynamicSmallCalendarRows]: false,
    [UserSettingKeys.FirstDayOfWeek]: 1,
    [UserSettingKeys.FontText]: "Atkinson Hyperlegible Next",
    [UserSettingKeys.FontTime]: "Atkinson Hyperlegible Mono",
    [UserSettingKeys.ThemeLight]: "Luna Light",
    [UserSettingKeys.ThemeDark]: "Luna Dark",
    [UserSettingKeys.UiScaling]: 1
  });
  public globalSettings: GlobalSettings = $state({
    [GlobalSettingKeys.LoggingVerbosity]: 2,
    [GlobalSettingKeys.RegistrationEnabled]: false,
    [GlobalSettingKeys.UseCdnFonts]: false,
  });

  constructor() {
    this.fetchSettings();
  }

  private async fetchUserData() {
    await fetchJson("/api/user").then((data: { user: UserData }) => {
      this.userData = data.user;
    }).catch((err) => {
      throw new Error("Could not get user data: " + err.message);
    });
  }

  private async fetchUserSettings() {
    await fetchJson("/api/settings/user").then((data: UserSettings) => {
      this.userSettings = data;
    }).catch((err) => {
      throw new Error("Could not get user settings: " + err.message);
    });
  }

  private async fetchGlobalSettings() {
    await fetchJson("/api/settings/global").then((data: GlobalSettings) => {
      this.globalSettings = data;
    }).catch((err) => {
      throw new Error("Could not get global settings: " + err.message);
    });
  }

  public async fetchSettings() {
    await Promise.all([
      this.fetchUserData(),
      this.fetchUserSettings(),
      this.fetchGlobalSettings()
    ]).catch((err) => {
      queueNotification("failure", "Could not fetch settings: " + err.message);
    });
  }
}

let settings: Settings | null = null;
export function getSettings() {
  if (settings === null) {
    settings = new Settings();
  }
  return settings;
}
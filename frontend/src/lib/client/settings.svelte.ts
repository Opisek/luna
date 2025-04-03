import { browser } from "$app/environment";
import { ColorKeys } from "../../types/colors";
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
    [UserSettingKeys.ThemeLight]: "luna-light",
    [UserSettingKeys.ThemeDark]: "luna-dark",
    [UserSettingKeys.UiScaling]: 1
  });
  public globalSettings: GlobalSettings = $state({
    [GlobalSettingKeys.LoggingVerbosity]: 2,
    [GlobalSettingKeys.RegistrationEnabled]: false,
    [GlobalSettingKeys.UseCdnFonts]: false,
  });

  constructor() {
    this.fetchFromStorage();
    if (browser) window.addEventListener("storage", () => this.fetchFromStorage());
    this.fetchSettings();
  }

  private async fetchUserData() {
    await fetchJson("/api/users/self").then((data: { user: UserData }) => {
      this.userData = data.user;
    }).catch((err) => {
      throw new Error("Could not get user data: " + err.message);
    });
  }

  private async fetchUserSettings() {
    await fetchJson("/api/users/self/settings").then((data: UserSettings) => {
      this.userSettings = data;
    }).catch((err) => {
      throw new Error("Could not get user settings: " + err.message);
    });
  }

  private async fetchGlobalSettings() {
    await fetchJson("/api/settings").then((data: GlobalSettings) => {
      this.globalSettings = data;
    }).catch((err) => {
      throw new Error("Could not get global settings: " + err.message);
    });
  }

  public async fetchSettings() {
    if (!browser || !document.cookie.includes("tokenPresent")) return;
    await Promise.all([
      this.fetchUserData(),
      this.fetchUserSettings(),
      this.fetchGlobalSettings()
    ]).catch((err) => {
      queueNotification(ColorKeys.Danger, "Could not fetch settings: " + err.message);
    });
  }

  // This only saves the settings to the local storage.
  // Saving to the database must be done separately by the caller.
  public saveSettings() {
    if (!browser) return;
    localStorage.setItem("userData", JSON.stringify(this.userData));
    localStorage.setItem("userSettings", JSON.stringify(this.userSettings));
    localStorage.setItem("globalSettings", JSON.stringify(this.globalSettings));
  }

  private fetchFromStorage() {
    if (!browser) return;
    const userData = localStorage.getItem("userData");
    if (userData != null) {
      this.userData = JSON.parse(userData);
    }
    const userSettings = localStorage.getItem("userSettings");
    if (userSettings != null) {
      this.userSettings = JSON.parse(userSettings);
    }
    const globalSettings = localStorage.getItem("globalSettings");
    if (globalSettings != null) {
      this.globalSettings = JSON.parse(globalSettings);
    }
  }
}

let settings: Settings | null = null;
export function getSettings() {
  if (settings === null) {
    settings = new Settings();
  }
  return settings;
}
import { Dataiku } from "svelte-simples";
import type { GlobalSettings, UserData, UserSettings } from "../../types/settings";
import { fetchJson } from "./net";
import { queueNotification } from "./notifications";

class Settings {
  private userData: UserData;
  private userSettings: UserSettings;
  private globalSettings: GlobalSettings;

  constructor() {
    this.userData = {} as UserData;
    this.userSettings = {} as UserSettings;
    this.globalSettings = {} as GlobalSettings;

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

  public getUserData() {
    return this.userData;
  }

  public getUserSettings() {
    return this.userSettings;
  }

  public getGlobalSettings() {
    return this.globalSettings;
  }
}

let settings: Settings | null = null;
export function getSettings() {
  if (settings === null) {
    settings = new Settings();
  }
  return settings;
}
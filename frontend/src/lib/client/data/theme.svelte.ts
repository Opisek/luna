import { browser } from "$app/environment";
import { UserSettingKeys } from "../../../types/settings";
import { getSettings } from "./settings.svelte";

class Theme {
  private lightMode: boolean = $state(false);
  private settings = getSettings();
  
  constructor() {
    if (!browser) return;
    this.fetchFromStorage();
    this.fetchFromSystem();
    window.addEventListener("storage", () => this.fetchFromStorage());
    window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', () => this.fetchFromSystem());
  }

  private fetchFromStorage() {
    if (this.settings.userSettings[UserSettingKeys.ThemeSynchronize]) return;
    const theme = localStorage.getItem("theme");
    if (theme != null) {
      this.lightMode = theme !== "dark";
    }
  }

  private fetchFromSystem() {
    if (!this.settings.userSettings[UserSettingKeys.ThemeSynchronize]) return;
    this.lightMode = !(window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches);
  }

  private saveToStorage() {
    if (!browser) return;
    localStorage.setItem("theme", this.lightMode ? "light" : "dark");
  }

  public setLightMode() {
    this.lightMode = true;
    this.saveToStorage();
  }

  public setDarkMode() {
    this.lightMode = false;
    this.saveToStorage();
  }

  public toggle() {
    this.lightMode = !this.lightMode;
    this.saveToStorage();
  }

  public isLightMode() {
    return this.lightMode;
  }

  public refetchTheme() {
    if (!browser) return;
    this.fetchFromStorage();
    this.fetchFromSystem();
  }
}

let theme: Theme | null = $state(null);
export function getTheme() {
  if (theme === null) {
    theme = new Theme();
  }
  return theme;
}

export function resetThemes() {
  theme = null;
}
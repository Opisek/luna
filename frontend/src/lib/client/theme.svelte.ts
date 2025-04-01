import { browser } from "$app/environment";

class Theme {
  private lightMode: boolean = $state(false);
  
  constructor() {
    if (!browser) return;
    this.fetchFromStorage();
    window.addEventListener("storage", () => this.fetchFromStorage());
  }

  private fetchFromStorage() {
    const theme = localStorage.getItem("theme");
    if (theme != null) {
      this.lightMode = theme === "light";
    }
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
}

let theme: Theme | null = null;
export function getTheme() {
  if (theme === null) {
    theme = new Theme();
  }
  return theme;
}
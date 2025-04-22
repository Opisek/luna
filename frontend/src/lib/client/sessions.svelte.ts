import { browser } from "$app/environment";
import ipLocation from "iplocation";
import { fetchJson } from "./net";

export function clearSession() {
  if (!browser) return;
  localStorage.clear();
  window.location.href = "/logout";
}

class ActiveSessions {
  public currentSession = $state("");
  public activeSessions = $state<Session[]>([]);

  private locationCache = new Map<string, string>();

  public async fetch() {
    await fetchJson("/api/sessions").then((data: { current: string, sessions: Session[] }) => {
      this.currentSession = data.current;
      data.sessions.forEach(async x => {
        const cached = this.locationCache.get(x.ip_address);
        if (cached) {
          x.location = cached;
        } else {
          const location = await ipLocation(x.ip_address);
          if (location.reserved) x.location = "Local Network";
          else x.location = ((loc => `${loc.country.name} ${loc.region.name} ${loc.city}`)(location as ipLocation.LocationData));
          this.locationCache.set(x.ip_address, x.location);
        }

        x.created_at = new Date(x.created_at);
        x.last_seen = new Date(x.last_seen);
      });
      this.activeSessions = data.sessions;
    });
  }
}

let activeSessions: ActiveSessions | null = null;
export function getActiveSessions() {
  if (activeSessions === null) {
    activeSessions = new ActiveSessions();
  }
  return activeSessions;
}
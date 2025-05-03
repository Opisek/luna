import { browser } from "$app/environment";
import ipLocation from "iplocation";
import { fetchJson, fetchResponse } from "../net";
import { resetSettings } from "./settings.svelte";
import { resetThemes } from "./theme.svelte";
import { resetRepository } from "./repository.svelte";
import { resetNotifications } from "../notifications";
import { resetMetadata } from "./metadata.svelte";
import { resetRegistrationInvites } from "./invites.svelte";

export function clearSession() {
  if (!browser) return;
  localStorage.clear();
  document.cookie = "tokenPresent=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";

  resetRegistrationInvites();
  resetSessions();
  resetSettings();
  resetThemes();
  resetRepository();
  resetNotifications();
  resetMetadata();

  window.location.href = "/login?expired=true";
}

class ActiveSessions {
  public currentSession = $state("");
  public activeSessions = $state<Session[]>([]);

  private locationCache = new Map<string, string>();

  public async fetch() {
    await fetchJson("/api/sessions").then((data: { current: string, sessions: Session[] }) => {
      this.currentSession = data.current;
      data.sessions.forEach(async x => {
        const cached = this.locationCache.get(x.last_ip_address);
        if (cached) {
          x.location = cached;
        } else if (["::1", "127.0.0.1", "localhost", ].includes(x.last_ip_address)) {
          x.location = "Local Host";
          this.locationCache.set(x.last_ip_address, x.location);
        } else {
          const location = await ipLocation(x.last_ip_address);
          if (location.reserved) x.location = "Local Network";
          else x.location = ((loc => `${loc.country.name} ${loc.region.name} ${loc.city}`)(location as ipLocation.LocationData));
          if (x.location.trim() === "") x.location = "Local Network"
          this.locationCache.set(x.last_ip_address, x.location);
        }

        x.created_at = new Date((x.created_at as unknown as string).replace("Z", ""));
        x.last_seen = new Date((x.last_seen as unknown as string).replace("Z", ""));
      });
      this.activeSessions = data.sessions;
    });
  }

  public async requestToken(session: Session, password: string): Promise<{ token: string, session: Session }> {
    const body = new FormData();
    body.append("name", session.user_agent);
    body.append("password", password);

    const token = (await fetchJson(`/api/sessions`, { method: "PUT", body: body })).token;
    await this.fetch();

    return {
      token: token,
      session: this.activeSessions.filter(x => x.is_api).sort((a, b) => a.created_at.getTime() - b.created_at.getTime())[0]
    };
  }

  public async updateSession(session: Session, password: string): Promise<Session> {
    const body = new FormData();
    body.append("name", session.user_agent);
    body.append("password", password);

    const token = (await fetchJson(`/api/sessions/${session.session_id}`, { method: "PATCH", body: body })).token;
    await this.fetch();

    return this.activeSessions.filter(x => x.session_id === session.session_id)[0]
  }

  public async deauthorizeSession(id: string) {
    return fetchResponse(`/api/sessions/${id}`, { method: "DELETE" }).then(() => {
      this.activeSessions = this.activeSessions.filter(x => x.session_id != id);
    });
  }

  public async deauthorizeUserSessions() {
    return fetchResponse("/api/sessions?type=user", { method: "DELETE" }).then(() => {
      clearSession();
    });
  }

  public async deauthorizeApiSessions() {
    return fetchResponse("/api/sessions?type=api", { method: "DELETE" }).then(() => {
      clearSession();
    });
  }
}

let activeSessions: ActiveSessions | null = $state(null);
export function getActiveSessions() {
  if (activeSessions === null) {
    activeSessions = new ActiveSessions();
  }
  return activeSessions;
}

export function resetSessions() {
  activeSessions = null;
}
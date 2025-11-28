import { browser } from "$app/environment";
import { page } from "$app/state";

import ipLocation from "iplocation";

import { fetchJson, fetchResponse } from "../net";
import type { Settings } from "./settings.svelte";
import { GlobalSettingKeys } from "../../../types/settings";

export function clearSession() {
  if (!browser) return;
  localStorage.clear();
  document.cookie = "tokenPresent=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
  window.location.href = "/login?expired=true";
}

export class ActiveSessions {
  public currentSession = $state("");
  public activeSessions = $state<Session[]>([]);

  private locationCache = new Map<string, string>();

  private settings: Settings;

  constructor(settings: Settings) {
    this.settings = settings;
  }

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
          const location = this.settings.globalSettings[GlobalSettingKeys.UseIpGeolocation] ? await ipLocation(x.last_ip_address) : null;
          if (location == null) x.location = "Unknown Location";
          else if (location.reserved) x.location = "Local Network";
          else {
            x.location = ((loc => `${loc.country.name} ${loc.region.name} ${loc.city}`)(location as ipLocation.ReturnType as ipLocation.LocationData)); // TypeScript bug, need to double cast even though it should be able to do the first one implicitly
            if (x.location.trim() === "") x.location = "Local Network"
          }
          this.locationCache.set(x.last_ip_address, x.location);
        }

        x.created_at = new Date((x.created_at as unknown as string).replace("Z", ""));
        x.last_seen = new Date((x.last_seen as unknown as string).replace("Z", ""));
      });
      this.activeSessions = data.sessions;
    });
  }

  public async getSessionPermissions(sessionId: string): Promise<string[]> {
    return fetchJson(`/api/sessions/${sessionId}/permissions`).then((data: { permissions: string[] }) => {
      return data.permissions;
    });
  }

  public async requestToken(session: Session, password: string): Promise<{ token: string, session: Session }> {
    const body = new FormData();
    body.append("name", session.user_agent);
    body.append("permissions", JSON.stringify(session.permissions));
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
    body.append("permissions", JSON.stringify(session.permissions));
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

export function getActiveSessions(): ActiveSessions {
  return page.data.singletons.sessions;
}
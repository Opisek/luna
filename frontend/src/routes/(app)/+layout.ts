import { redirect, type LoadEvent } from "@sveltejs/kit";
import { fetchJsonFromEvent } from "$lib/scripts/client/net";
import type { PageLoad } from "./register/$types";
import { unprivilegedPaths } from "$lib/scripts/common/paths";
import { NoOp } from "$lib/scripts/client/placeholders";
import { isCompatibleWithBackend, VersionCompatibility } from "$lib/scripts/common/version";

import { ActiveSessions } from "$lib/scripts/client/data/sessions.svelte";
import { Connectivity } from "$lib/scripts/client/data/connectivity.svelte";
import { Metadata } from "$lib/scripts/client/data/metadata.svelte";
import { RegistrationInvites } from "$lib/scripts/client/data/invites.svelte";
import { Repository } from "$lib/scripts/client/data/repository.svelte";
import { Theme } from "$lib/scripts/client/data/theme.svelte";
import { Users } from "$lib/scripts/client/data/users.svelte";
import { Settings } from "$lib/scripts/client/data/settings.svelte";
import { OauthClients } from "$lib/scripts/client/data/oauth.svelte";

import "$lib/scripts/common/i18n";
import { UserSettingKeys } from "$lib/types/settings";
import { loadLanguage } from "$lib/scripts/common/i18n";
import { encodeRedirectUrl } from "$lib/scripts/common/url";

function getSingletons(version: string, preloadedSettings: { userData: any, userSettings: any, globalSettings: any } | null = null): {
  connectivity: Connectivity;
  invites: RegistrationInvites;
  metadata: Metadata;
  repository: Repository;
  sessions: ActiveSessions;
  settings: Settings;
  theme: Theme;
  users: Users;
  oauthClients: OauthClients;
} {
  let rep: Repository | null = null;
  let recalculateEvents = (calendarThatBecameVisible: (string | null)) => {
    if (rep == null) return;
    rep.recalculateEvents(calendarThatBecameVisible);
  }

  let connectivity = new Connectivity(version);
  let invites = new RegistrationInvites();
  let metadata = new Metadata(recalculateEvents);
  let repository = new Repository(metadata);
  let settings = new Settings(preloadedSettings);
  let sessions = new ActiveSessions(settings);
  let theme = new Theme(settings);
  let users = new Users();
  let oauthClients = new OauthClients();

  rep = repository;

  return {
    connectivity,
    invites,
    metadata,
    repository,
    sessions,
    settings,
    theme,
    users,
    oauthClients
  };
}

export const load: PageLoad = async (event: LoadEvent) => {
  const response = await fetchJsonFromEvent(event, "/api/version", {}, true).catch(NoOp);
  const version = response?.version;
  if (event.url.pathname !== "/version" && [VersionCompatibility.BackendOutdatedMajor, VersionCompatibility.FrontendOutdatedMajor].includes(isCompatibleWithBackend(version)))
    redirect(302, `/version?redirect=${encodeURIComponent(event.url.pathname)}`);

  for (const path of unprivilegedPaths) {
    if (!event.url.pathname.startsWith(`/${path}`)) continue;
    await loadLanguage(null);
    return {
      version: version,
      singletons: getSingletons(version)
    }; 
  }

  const results = await Promise.all([
    fetchJsonFromEvent(event, "/api/users/self", {}, true),
    fetchJsonFromEvent(event, "/api/users/self/settings", {}, true),
    fetchJsonFromEvent(event, "/api/settings", {}, true)
  ]).catch((err) => {
    if (err.message.includes("Unauthorized") || err.message.includes("Session expired")) {
      redirect(302, `/login?redirect=${encodeRedirectUrl(new URL(document.location.href))}&expired=true`);
    }
    return null; 
  });

  const settingsLoaded = results && results[0].user !== undefined;

  await loadLanguage(((results || {})[1] || {})[UserSettingKeys.Language])

  if (!settingsLoaded) return {
    version: version,
    singletons: getSingletons(version)
  };

	return {
    version: version,
    userData: results[0].user,
    userSettings: results[1],
    globalSettings: results[2],
    singletons: getSingletons(version, { userData: results[0].user, userSettings: results[1], globalSettings: results[2] }),
	};
};
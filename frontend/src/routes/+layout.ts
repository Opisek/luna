import { redirect, type LoadEvent } from "@sveltejs/kit";
import { fetchJsonFromEvent } from "../lib/client/net";
import type { PageLoad } from "./register/$types";
import { unprivilegedPaths } from "../lib/common/paths";
import { NoOp } from "$lib/client/placeholders";
import { isCompatibleWithBackend, VersionCompatibility } from "$lib/common/version";

import { ActiveSessions } from "../lib/client/data/sessions.svelte";
import { Connectivity } from "../lib/client/data/connectivity.svelte";
import { Metadata } from "../lib/client/data/metadata.svelte";
import { RegistrationInvites } from "../lib/client/data/invites.svelte";
import { Repository } from "../lib/client/data/repository.svelte";
import { Theme } from "../lib/client/data/theme.svelte";
import { Users } from "../lib/client/data/users.svelte";
import { Settings } from "../lib/client/data/settings.svelte";

function getSingletons(version: string, preloadedSettings: { userData: any, userSettings: any, globalSettings: any } | null = null): {
  connectivity: Connectivity;
  invites: RegistrationInvites;
  metadata: Metadata;
  repository: Repository;
  sessions: ActiveSessions;
  settings: Settings;
  theme: Theme;
  users: Users;
} {
  let rep: Repository | null = null;
  let recalculateEvents = () => {
    if (rep == null) return;
    rep.recalculateEvents();
  }

  let connectivity = new Connectivity(version);
  let invites = new RegistrationInvites();
  let metadata = new Metadata(recalculateEvents);
  let repository = new Repository(metadata);
  let settings = new Settings(preloadedSettings);
  let sessions = new ActiveSessions(settings);
  let theme = new Theme(settings);
  let users = new Users();

  rep = repository;

  return {
    connectivity,
    invites,
    metadata,
    repository,
    sessions,
    settings,
    theme,
    users
  };
}

export const load: PageLoad = async (event: LoadEvent) => {
  const response = await fetchJsonFromEvent(event, "/api/version", {}, true).catch(NoOp);
  const version = response?.version;
  if (event.url.pathname !== "/version" && [VersionCompatibility.BackendOutdatedMajor, VersionCompatibility.FrontendOutdatedMajor].includes(isCompatibleWithBackend(version)))
    redirect(302, `/version?redirect=${encodeURIComponent(event.url.pathname)}`);

  for (const path of unprivilegedPaths) if (event.url.pathname.startsWith(`/${path}`)) return {
    version: version,
    singletons: getSingletons(version)
  }; 

  const results = await Promise.all([
    fetchJsonFromEvent(event, "/api/users/self", {}, true),
    fetchJsonFromEvent(event, "/api/users/self/settings", {}, true),
    fetchJsonFromEvent(event, "/api/settings", {}, true)
  ]).catch((err) => {
    if (err.message.includes("Unauthorized") || err.message.includes("Session expired")) {
      redirect(302, `/login?redirect=${encodeURIComponent(event.url.pathname)}&expired=true`);
    }
    return null; 
  });

  if (!results || results[0].user === undefined) return {
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
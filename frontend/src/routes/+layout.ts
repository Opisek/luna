import { redirect, type LoadEvent } from "@sveltejs/kit";
import { fetchJsonFromEvent } from "../lib/client/net";
import type { PageLoad } from "./register/$types";
import { unprivilegedPaths } from "../lib/common/paths";
import { NoOp } from "$lib/client/placeholders";
import { isCompatibleWithBackend, VersionCompatibility } from "$lib/common/version";

export const load: PageLoad = async (event: LoadEvent) => {
  const version = (await fetchJsonFromEvent(event, "/api/version", {}, true).catch(NoOp)).version || null;
  if (event.url.pathname !== "/version" && [VersionCompatibility.BackendOutdatedMajor, VersionCompatibility.FrontendOutdatedMajor].includes(isCompatibleWithBackend(version)))
    redirect(302, `/version?redirect=${encodeURIComponent(event.url.pathname)}`);

  for (const path of unprivilegedPaths) if (event.url.pathname.startsWith(`/${path}`)) return { version: version }; 

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

  if (!results || results[0].user === undefined) return { version: version };

	return {
    version: version,
    userData: results[0].user,
    userSettings: results[1],
    globalSettings: results[2]
	};
};
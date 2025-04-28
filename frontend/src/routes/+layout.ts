import { redirect, type LoadEvent } from "@sveltejs/kit";
import { fetchJsonFromEvent } from "../lib/client/net";
import type { PageLoad } from "./register/$types";
import { unprivilegedPaths } from "../lib/common/paths";

export const load: PageLoad = async (event: LoadEvent) => {
  for (const path of unprivilegedPaths) {
    if (event.url.pathname.startsWith(`/${path}`)) {
      return null; 
    }
  }

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

  if (!results || results[0].user === undefined) return null;

	return {
    userData: results[0].user,
    userSettings: results[1],
    globalSettings: results[2]
	};
};
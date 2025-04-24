import type { LoadEvent } from "@sveltejs/kit";
import { fetchJsonFromEvent } from "../lib/client/net";
import type { PageLoad } from "./register/$types";

export const load: PageLoad = async (event: LoadEvent) => {
  const results = await Promise.all([
    fetchJsonFromEvent(event, "/api/users/self", {}, true),
    fetchJsonFromEvent(event, "/api/users/self/settings", {}, true),
    fetchJsonFromEvent(event, "/api/settings", {}, true)
  ])

	return {
    userData: results[0].user,
    userSettings: results[1],
    globalSettings: results[2]
	};
};
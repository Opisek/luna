import type { PageLoad } from "./$types";
import { fetchJsonFromEvent } from "../../lib/client/net";
import type { LoadEvent } from "@sveltejs/kit";

export const load: PageLoad = async (event: LoadEvent) => {
  const response = await fetchJsonFromEvent(event, "/api/register/enabled");

	return {
    registrationEnabled: response.enabled
	};
};
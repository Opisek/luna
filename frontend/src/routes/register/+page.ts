import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async (event) => {
  const response = await event.fetch("/api/register/enabled").catch((err) => {
    if (!err.cause) err.cause = new Error("500");
    error(Number.parseInt(err.cause.message), err.message);
  });
  if (!response.ok) {
    error(response.status, response.statusText);
  }
  const body = await response.json().catch((err) => {
    error(500, err.message);
  });

	return {
    registrationEnabled: body.enabled
	};
};
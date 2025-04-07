import { redirect, type Actions } from "@sveltejs/kit";

import { COOKIE_MAX_AGE } from "$lib/server/constants.server";
import { callApi } from "$lib/server/api.server";
import { getRedirectPage } from "$lib/common/parsing";

export const actions = {
  default: async ({cookies, request, getClientAddress}) => {
    const formData = await request.formData();

    const password = formData.get("password");
    const passwordRepeat = formData.get("passwordRepeat");

    if (password !== passwordRepeat) {
      return {
        status: 400,
        error: "Passwords do not match"
      };
    }

    const res = await callApi(request, getClientAddress(), "register", { method: "POST", body: formData });

    if (res.ok) {
      const body = await res.json();

      const opts = {
        path: "/",
        httpOnly: false,
        maxAge: undefined as number | undefined,
        sameSite: "strict" as boolean | "strict" | "lax" | "none",
      };

      if (formData.get("remember") === "true") {
        opts.maxAge = COOKIE_MAX_AGE;
      }

      cookies.set("tokenPresent", "true", opts);
      cookies.set("token", body.token, opts);

      redirect(302, getRedirectPage(new URL(request.url)));
    } else {
      return await res.json().catch(() => ({ error: `Error ${res.status}: ${res.statusText}` }));
    }
  }
} satisfies Actions;
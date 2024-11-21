import { redirect, type Actions } from "@sveltejs/kit";

import { COOKIE_MAX_AGE } from "$lib/server/constants";
import { callApi } from "$lib/server/api.server";
import { getRedirectPage } from "$lib/common/parsing";

export const actions = {
  default: async ({cookies, request}) => {
    const formData = await request.formData();

    const password = formData.get("password");
    const passwordRepeat = formData.get("passwordRepeat");

    if (password !== passwordRepeat) {
      return {
        status: 400,
        error: "Passwords do not match"
      };
    }

    const res = await callApi("register", { method: "POST", body: formData });

    if (res.ok) {
      const body = await res.json();

      const opts = {
        path: "/",
        httpOnly: false,
        maxAge: undefined as number | undefined,
      };

      if (formData.get("remember") === "true") {
        opts.maxAge = COOKIE_MAX_AGE;
      }

      // TODO: max age
      cookies.set("tokenPresent", "true", opts);
      cookies.set("token", body.token, opts);

      redirect(302, getRedirectPage(new URL(request.url)));
    } else {
      return await res.json();
    }
  }
} satisfies Actions;
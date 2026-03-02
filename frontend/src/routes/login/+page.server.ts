import { fail, redirect, type Actions } from "@sveltejs/kit";

import { COOKIE_MAX_AGE } from "$lib/server/constants.server";
import { apiProxy } from "$lib/server/api.server";
import { getRedirectPage } from "$lib/common/parsing";
import { invalidateAll } from "$app/navigation";

export const actions = {
  default: async ({cookies, request, getClientAddress}) => {
    const formData = await request.formData();

    const res = await apiProxy(request, getClientAddress, "login", { method: "POST", body: formData }, false);

    if (res.ok) {
      const body = await res.json();

      const opts = {
        path: "/",
        httpOnly: false,
        maxAge: undefined as number | undefined,
        sameSite: "strict" as boolean | "strict" | "lax" | "none",
        secure: !(process.env.PUBLIC_URL?.startsWith("http://") || false),
      };

      if (formData.get("remember") === "true") {
        opts.maxAge = COOKIE_MAX_AGE;
      }

      cookies.set("tokenPresent", "true", opts);
      opts.httpOnly = true;
      cookies.set("token", body.token, opts);

      redirect(302, getRedirectPage(new URL(request.url)));
    } else {
      const error = await res.json().catch(() => ({ error: `${res.statusText}` }));
      return fail(res.status, { error: error.error });
    }
  }
} satisfies Actions;
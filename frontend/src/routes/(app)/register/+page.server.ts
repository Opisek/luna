import { fail, redirect, type Actions } from "@sveltejs/kit";

import { COOKIE_MAX_AGE } from "$lib/server/constants.server";
import { apiProxy } from "$lib/server/api.server";
import { getRedirectPage } from "$lib/common/parsing";
import { t } from "@sveltia/i18n";

export const actions = {
  default: async ({cookies, request, getClientAddress}) => {
    const formData = await request.formData();

    const password = formData.get("password");
    const passwordRepeat = formData.get("password_repeat");

    if (password !== passwordRepeat) {
      return {
        status: 400,
        error: t("validation.password.match")
      };
    }

    const res = await apiProxy(request, getClientAddress, "register", { method: "POST", body: formData }, false).catch(() => null);
    if (!res) return fail(500, { error: "The backend server cannot be reached." });

    if (res.ok) {
      const body = await res.json();

      const opts = {
        path: "/",
        httpOnly: false,
        maxAge: undefined as number | undefined,
        sameSite: "lax" as boolean | "strict" | "lax" | "none",
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
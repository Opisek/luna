import { redirect, type Actions } from "@sveltejs/kit";
import { callApi } from "../../lib/server/api.server";
import { getRedirectPage } from "../../lib/common/parsing";

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

      // TODO: max age
      cookies.set("tokenPresent", "true", {
        path: "/",
        httpOnly: false
      });
      cookies.set("token", body.token, {
        path: "/",
        httpOnly: true
      });

      redirect(302, getRedirectPage(new URL(request.url)));
    } else {
      return await res.json();
    }
  }
} satisfies Actions;
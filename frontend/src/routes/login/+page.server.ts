import { redirect, type Actions } from "@sveltejs/kit";
import { callApi } from "../../lib/server/api.server";

export const actions = {
  default: async ({cookies, request}) => {
    const data = await request.formData();

    const res = await callApi("login", { method: "POST", body: data });

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

      const url = new URL(request.url);
      let redirectPage = url.searchParams.get("redirect");
      if (redirectPage == null || redirectPage == "") {
        redirectPage = '/';
      } else {
        redirectPage = decodeURIComponent(redirectPage);
      }

      redirect(302, redirectPage);
    } else {
      return await res.json();
    }
  }
} satisfies Actions;
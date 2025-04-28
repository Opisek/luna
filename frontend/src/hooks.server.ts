import type { Handle } from "@sveltejs/kit";
import { getRedirectPage } from "./lib/common/parsing";

import "dotenv/config"
import { loginPaths, unprivilegedPaths } from "./lib/common/paths";

export const handle: Handle = async ({ event, resolve }) => {
  const tokenPresent = event.cookies.get("tokenPresent");
  const currentUrl = new URL(event.request.url)

  if (tokenPresent) {
    let isLogin = false
    for (const path of loginPaths) {
      if (currentUrl.pathname.startsWith(`/${path}`)) {
        isLogin = true
        break
      }
    }

    const sessionExpired = event.url.searchParams.get("expired") === "true";

    if (isLogin) {
      if (sessionExpired) {
        const opts = {
          path: "/",
          httpOnly: true,
          maxAge: undefined as number | undefined,
          sameSite: "strict" as boolean | "strict" | "lax" | "none",
          secure: !(process.env.PUBLIC_URL?.startsWith("http://") || false),
        };

        event.cookies.delete("tokenPresent", opts);
        opts.httpOnly = true;
        event.cookies.delete("token", opts);
      } else {
        return new Response(null, {
          status: 302,
          headers: {
            location: getRedirectPage(currentUrl),
          },
        })
      }
    }
  } else {
    let isUnprivileged = false
    for (const path of unprivilegedPaths) {
      if (currentUrl.pathname.startsWith(`/${path}`)) {
        isUnprivileged = true
        break
      }
    }
    
    if (!isUnprivileged) {
      return new Response(null, {
        status: 302,
        headers: {
          location: `/login?redirect=${encodeURIComponent(event.request.url)}`,
        },
      })
    }
  }

  return await resolve(event)
}
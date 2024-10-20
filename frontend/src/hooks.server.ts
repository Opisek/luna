import type { Handle } from "@sveltejs/kit";
import { getRedirectPage } from "./lib/common/parsing";

import "dotenv/config"

const loginPaths = [
  "login",
  "register",
  "recover"
]

const unprivilegedPaths = loginPaths.concat([
  "api"
])

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

    if (isLogin) {
      return new Response(null, {
        status: 302,
        headers: {
          location: getRedirectPage(currentUrl),
        },
      })
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
import { redirect, type Actions } from "@sveltejs/kit";

import { getRedirectPage } from "$lib/common/parsing";

// Warning: This page assumes you had already sent a request to the backend to deauthorize this session.
export const GET = async ({cookies}) => {
  const opts = {
    path: "/",
    httpOnly: true,
    maxAge: undefined as number | undefined,
    sameSite: "strict" as boolean | "strict" | "lax" | "none",
    secure: !(process.env.PUBLIC_URL?.startsWith("http://") || false),
  };

  cookies.delete("tokenPresent", opts);
  opts.httpOnly = true;
  cookies.delete("token", opts);

  redirect(302, "/login");
};
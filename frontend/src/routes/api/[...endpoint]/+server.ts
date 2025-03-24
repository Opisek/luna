import { callApi } from "../../../lib/server/api.server";
import type { RequestEvent } from "./$types";
import { error } from "@sveltejs/kit";

const proxy = (async ({ params, request, url }: RequestEvent) => {
  // CSRF protection
  if (request.method !== "GET" && request.method !== "HEAD") {
    const origin = request.headers.get("Origin");
    const publicUrl = process.env.PUBLIC_URL
    if (publicUrl === null || publicUrl === undefined || publicUrl === "") {
    // @ts-ignore
      return new error(403, { message: "Environmental variable PUBLIC_URL was not set." });
    }
    const allowedOrigins = publicUrl.split(",").map(x => x.trim());
    if (!origin || !allowedOrigins.includes(origin)) {
    // @ts-ignore
      return new error(403, { message: "Origin not allowed. Was PUBLIC_URL set correctly?" });
    }
  }

  // API call to the backend
  return callApi(params.endpoint + url.search, request).catch((err) => {
    let errorMessage = "Internal Server Error";

    if (err.cause && err.cause.code === "ECONNREFUSED") {
      errorMessage = "The backend is not reachable";
    }

    // @ts-ignore
    return new error(500, { message: errorMessage } );
  });
})

export const DELETE = proxy
export const GET = proxy;
export const PATCH = proxy;
export const POST = proxy;
export const PUT = proxy;
export const HEAD = proxy;
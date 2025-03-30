import { callApi } from "../../../lib/server/api.server";
import type { RequestEvent } from "./$types";
import { error } from "@sveltejs/kit";

const proxy = (async ({ params, request, url }: RequestEvent) => {
  // CSRF protection
  if (request.method !== "GET" && request.method !== "HEAD") {
    const origin = request.headers.get("Origin");
    const publicUrl = process.env.PUBLIC_URL
    if (publicUrl === null || publicUrl === undefined || publicUrl === "" || origin !== publicUrl) {
    // @ts-ignore
      return new error(403, { message: "Origin not allowed. Was PUBLIC_URL set correctly in the frontend?" });
    }
  }

  // API call to the backend
  const response = await callApi(params.endpoint + url.search, request).catch((err) => {
    let errorMessage = "Internal Server Error";

    if (err.cause && err.cause.code === "ECONNREFUSED") {
      errorMessage = "The backend is not reachable";
    }

    // @ts-ignore
    return new error(500, { message: errorMessage } );
  });

  if (!response.ok) {
    return response;
  }

  // CORS check
  const corsHeader = response.headers.get("Access-Control-Allow-Origin")
  if (corsHeader === null || corsHeader === undefined || corsHeader === "" || corsHeader === "*" || corsHeader !== process.env.PUBLIC_URL) {
    // @ts-ignore
    return new error(403, { message: "Unexpected CORS header. Was PUBLIC_URL set correctly in the backend?" });
  }

  return response;
})

export const DELETE = proxy
export const GET = proxy;
export const PATCH = proxy;
export const POST = proxy;
export const PUT = proxy;
export const HEAD = proxy;
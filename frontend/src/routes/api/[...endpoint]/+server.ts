import { callApi } from "../../../lib/server/api.server";
import type { RequestEvent } from "./$types";
import { error } from "@sveltejs/kit";

const proxy = (async ({ params, request, url }: RequestEvent) => {
  // API call to the backend
  const response = await callApi(request, params.endpoint + url.search, request).catch((err) => {
    let errorMessage = "Internal Server Error";

    if (err.cause && err.cause.code === "ECONNREFUSED") {
      errorMessage = "The backend is not reachable";
    }

    // @ts-ignore
    return new error(err.status || 500, { message: errorMessage } );
  });

  return response;
})

export const DELETE = proxy
export const GET = proxy;
export const PATCH = proxy;
export const POST = proxy;
export const PUT = proxy;
export const HEAD = proxy;
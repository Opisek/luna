import { callApi } from "../../../lib/server/api.server";
import type { RequestEvent } from "./$types";
import { error } from "@sveltejs/kit";

const proxy = (async ({ params, request, url }: RequestEvent) => {
  return await callApi(params.endpoint + url.search, request).catch((err) => {
    let errorMessage = "Internal Server Error";

    if (err.cause && err.cause.code === "ECONNREFUSED") {
      errorMessage = "The backend is not reachable";
    }

    // @ts-ignore
    return new error(500, { message: errorMessage} );
  });
})

export const DELETE = proxy
export const GET = proxy;
export const PATCH = proxy;
export const POST = proxy;
export const PUT = proxy;
export const HEAD = proxy;
import { callApi } from "../../../lib/server/api.server";
import type { RequestEvent } from "./$types";
//import { API_URL } from "$env/static/private";

const proxy = (async ({ params, request }: RequestEvent) => {
  return await callApi(params.endpoint, request);
})

export const DELETE = proxy
export const GET = proxy;
export const PATCH = proxy;
export const POST = proxy;
export const PUT = proxy;
import type { RequestEvent } from "./$types";
import { API_URL } from "$env/static/private";

function proxy(method: string) {
  return (async ({ params }: RequestEvent) => {
    return await fetch(`${API_URL}/api/${params.endpoint}`, { method: method });
  })
}

export const DELETE = proxy("DELETE");
export const GET = proxy("GET");
export const PATCH = proxy("PATCH");
export const POST = proxy("POST");
export const PUT = proxy("PUT");
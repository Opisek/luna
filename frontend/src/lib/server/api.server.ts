import { json } from "@sveltejs/kit";

export async function apiProxy(request: Request, clientAddress: string, endpoint: string, init?: RequestInit, stream?: boolean): Promise<Response> {
  // CSRF protection
  if (request.method !== "GET" && request.method !== "HEAD") {
    const origin = request.headers.get("Origin");
    const publicUrl = process.env.PUBLIC_URL
    if (publicUrl === null || publicUrl === undefined || publicUrl === "" || origin !== publicUrl) {
      return json({ error: "Origin not allowed. Was PUBLIC_URL set correctly in the frontend?" }, { status: 403 });
    }
  }

  const originalHeaders: HeadersInit = [ ...request.headers ];
  if (!init) init = {};
  // @ts-ignore
  if (stream) init.duplex = "half";
  init.headers = [
    ...originalHeaders.filter(entry => !(stream ? [] : ["content-length", "content-type"]).includes(entry[0].toLowerCase())),
    [ "X-Forwarded-For", clientAddress ],
  ];
  const response = await fetch(`${process.env.API_URL}/api/${endpoint}`, init).catch((error) => {
    throw error; // TODO: maybe format as return json(...) too?
  });

  if (!response.ok) return response;

  // CORS check
  const corsHeader = response.headers.get("Access-Control-Allow-Origin")
  if (corsHeader === null || corsHeader === undefined || corsHeader === "" || corsHeader === "*" || corsHeader !== process.env.PUBLIC_URL) {
    return json({ error: "Unexpected CORS header. Was PUBLIC_URL set correctly in the backend?" }, { status: 403 });
  }

  return response;
}
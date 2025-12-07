import { json } from "@sveltejs/kit";

// Secure any API call, whether forwarded to the backend or not.
// This includes checking the Origin header of the request and the Access-Control-Allow-Origin of the response.
export async function secureApiCall(request: Request, func: () => Promise<Response>): Promise<Response> {
  // CSRF protection
  if (request.method !== "GET" && request.method !== "HEAD") {
    const origin = request.headers.get("Origin");
    const publicUrl = process.env.PUBLIC_URL
    if (publicUrl === null || publicUrl === undefined || publicUrl === "" || origin !== publicUrl) {
      return json({ error: "Origin not allowed. Was PUBLIC_URL set correctly in the frontend?" }, { status: 403 });
    }
  }

  // Process the request
  const response = await func();
  if (!response.ok) return response;

  // CORS check
  const corsHeader = response.headers.get("Access-Control-Allow-Origin")
  if (corsHeader === null || corsHeader === undefined || corsHeader === "" || corsHeader === "*" || corsHeader !== process.env.PUBLIC_URL) {
    return json({ error: "Unexpected CORS header. Was PUBLIC_URL set correctly in the backend?" }, { status: 403 });
  }

  return response;
}

// Forward the request to the backend.
export async function apiProxy(request: Request, getClientAddress: () => string, endpoint: string, init?: RequestInit, stream?: boolean): Promise<Response> {
  return secureApiCall(request, async () => {
    const originalHeaders: HeadersInit = [ ...request.headers ];
    if (!init) init = {};
    // @ts-ignore
    if (stream) init.duplex = "half";
    init.headers = [
      ...originalHeaders.filter(entry => !(stream ? [] : ["content-length", "content-type"]).includes(entry[0].toLowerCase())),
      [ "X-Forwarded-For", determineClientAddress(request, getClientAddress) ],
    ];
    const response = await fetch(`${process.env.API_URL}/api/${endpoint}`, init).catch((error) => {
      throw error; // TODO: maybe format as return json(...) too?
    });
    return response;
  });
}

// Try to determine the client IP address
export function determineClientAddress(request: Request, getClientAddress: () => string): string {
  let addr: string;
  if (request.headers.has("X-Real-IP")) {
    addr = request.headers.get("X-Real-IP") as string;
    if (addr.length > 0) return addr;
  }
  if (request.headers.has("True-Client-IP")) {
    addr = request.headers.get("True-Client-IP") as string;
    if (addr.length > 0) return addr;
  }
  if (request.headers.has("X-Forwarded-For")) {
    addr = ((request.headers.get("X-Forwarded-For") as string).split(",")[0] || "").trim();
    if (addr.length > 0) return addr;
  }
  return getClientAddress();
}
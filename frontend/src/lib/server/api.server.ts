import { json } from "@sveltejs/kit";

export async function callApi(request: Request, endpoint: string, init?: RequestInit): Promise<Response> {
  // CSRF protection
  if (request.method !== "GET" && request.method !== "HEAD") {
    const origin = request.headers.get("Origin");
    const publicUrl = process.env.PUBLIC_URL
    if (publicUrl === null || publicUrl === undefined || publicUrl === "" || origin !== publicUrl) {
      return json({ error: "Origin not allowed. Was PUBLIC_URL set correctly in the frontend?" }, { status: 403 });
    }
  }

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
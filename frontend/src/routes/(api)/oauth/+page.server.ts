import { apiProxy } from '$lib/server/api.server.js';
import type { PageServerLoad } from './$types.js';

// This could also be in client's load, but the API call has to travel through
// the frontend anyway, so this approach saves some back-and-forth requests.
export const load: PageServerLoad = async({url, request, getClientAddress}) => {
  const requestId = url.searchParams.get("state");

  if (!requestId || requestId == "") {
    return {
      error: "No state provided"
    }
  }

  let error = url.searchParams.get("error");
  if (error && error != "") {
    apiProxy(request, getClientAddress, `oauth/authorization/${encodeURIComponent(requestId)}`, { method: "DELETE" }, false);
    return { status: "error", request: requestId, error: error };
  }

  const authCode = url.searchParams.get("code");
  if (!authCode || authCode == "") {
    apiProxy(request, getClientAddress, `oauth/authorization/${encodeURIComponent(requestId)}`, { method: "DELETE" }, false);
    return { status: "error", request: requestId, error: "No authorization code was provided" };
  }


  const formData = new FormData();
  formData.append("authorization_code", authCode)

  const res = await apiProxy(request, getClientAddress, `oauth/authorization/${encodeURIComponent(requestId)}`, { method: "POST", body: formData }, false);

  if (res.ok) {
    let response: any = { status: "ok", request: requestId };
    const json = await res.json().catch(() => ([]));
    if (json && json.warnings) response.warnings = json.warnings;
    return response;
  } else {
    const error = await res.json().catch(() => ({ error: `${res.statusText}` }));
    return { status: "error", request: requestId, error: error.error };
  }
};
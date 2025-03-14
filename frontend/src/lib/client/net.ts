export async function fetchResponse(url: string, options: RequestInit = {}): Promise<Response> {
  const response = await fetch(url, options).catch((err) => {
    if (!err) err = new Error("Could not contact server");
    throw err;
  });
  if (response.ok) {
    return response;
  } else {
    const json = await response.json().catch(() => null);
    let err = null;
    if (!err && json != null) err = json.error;
    if (!err && json != null) err = json.message;
    if (!err) err = `${response.statusText ? response.statusText : "Could not contact server"} (${response.status})`;
    throw new Error(err);
  }
}

export async function fetchJson(url: string, options: RequestInit = {}) {
  return (await fetchResponse(url, options).catch(err => { throw err; })).json();
}
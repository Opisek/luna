import { clearSession } from "./sessions.svelte";

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
    if (err.includes("Session expired")) clearSession();
    throw new Error(err);
  }
}

export async function fetchJson(url: string, options: RequestInit = {}) {
  const json = await (await fetchResponse(url, options).catch(err => { throw err; })).json();
  if (json?.warnings) {
    // TODO: show warnings as notifications
  }
  return json;
}

export function downloadFileToClient(file: FileList | string | null) {
  if (file === null) return;

  const a = document.createElement("a");

  let url: string;
  if (typeof file === "string") {
    url = `/api/files/${file}`;
  } else {
    const blob = new Blob([file[0]], { type: file[0].type });
    url = URL.createObjectURL(blob);
    a.download = file[0].name;
  }

  a.href = url;

  a.click();
  URL.revokeObjectURL(url);
  a.remove();
}

export async function fetchFileById(fileId: string) {
  const res = await fetchResponse(`/api/files/${fileId}`, { method: "HEAD" }).catch(err => {
    throw err;
  });

  let filename = `${fileId}`;

  const header = res.headers.get("Content-Disposition")
  if (header) {
    const remoteFilename = header
      .split(";")
      .map(x => x.trim())
      .filter(x => x.startsWith("filename="))
      .map(x => x.split("=")[1]);
    
    if (remoteFilename.length > 0) filename = remoteFilename[0];
  }

  // https://stackoverflow.com/questions/52078853/is-it-possible-to-update-filelist
  const list = new DataTransfer();
  const file = new File([], filename);
  list.items.add(file);
  return list.files;
}
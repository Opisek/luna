import hash from "hash.js";

export function getSha1Hash(data: string): string {
  return hash.sha1().update(data).digest("hex");
}
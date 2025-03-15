import sha1 from "sha1";

// Note: The library object-hash is really cool, but it always adds type prefixes, e.g. string:5:hello, which is not what we want.
// If they ever add an option to disable this, we can consider switching due to its much higher flexibility and configurability.
export function getSha1Hash(data: string): string {
  return sha1(data);
}
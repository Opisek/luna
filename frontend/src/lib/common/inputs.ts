export function passIfEnter(e: KeyboardEvent, callback: () => any): boolean {
  if (!["Enter", " "].includes(e.key)) return false;
  callback();
  return true;
}
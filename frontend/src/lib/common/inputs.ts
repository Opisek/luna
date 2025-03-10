export function passIfEnter(e: KeyboardEvent, callback: () => any) {
  if (e.key === "Enter") {
    callback();
  }
}
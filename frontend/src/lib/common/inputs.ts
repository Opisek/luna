export function passIfEnter(e: KeyboardEvent, callback: () => any) {
  if (["Enter", " "].includes(e.key)) callback();
}
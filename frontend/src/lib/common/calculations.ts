export const calculateOptimalPopupPosition = (el: HTMLElement, horizontalParts: number = 3): { bottom: boolean, right: boolean, center: boolean } => {
  if (!el) return { bottom: false, right: false, center: false };

  const rect = el.getBoundingClientRect();

  const x = rect.left + (rect.right - rect.left) / 2;
  const y = rect.top + (rect.bottom - rect.top) / 2;

  const bottom = y > window.innerHeight / 2;
  const right = x > window.innerWidth / 2;
  const center = horizontalParts <= 2 ? false : (x < window.innerWidth / horizontalParts * (horizontalParts - 1) && x > window.innerWidth / horizontalParts);

  return { bottom, right, center };
}
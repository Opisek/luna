import BarFocusIndicator from "../../components/decoration/focus/BarFocusIndicator.svelte";
import UnderlineFocusIndicator from "../../components/decoration/focus/UnderlineFocusIndicator.svelte";
import Ripple from "../../components/decoration/Ripple.svelte"

export const addRipple = (e: MouseEvent, addToParent: boolean = true) => {
  if (!e.target) return;

  const parent = e.target as HTMLElement;

  new Ripple({ target: addToParent ? parent : e.target as HTMLElement, props: { event: e, parent: parent } });
}

//export const removeRipple = (e: MouseEvent) => {
//  if (!e.target) return;
//
//  const parent = e.target as HTMLElement;
//
//  const ripple = parent.querySelector(".ripple");
//
//  if (ripple) {
//    ripple.remove();
//  }
//
//  return;
//}

export const focusIndicator = (node: HTMLElement, settings: FocusIndicatorSettings = { type: "bar" }) => {
  const mouseDown = () => {
    node.classList.add("clicked");
  }

  const focusOut = (e: FocusEvent) => {
    if (settings.ignoreParent && e.relatedTarget && node.parentElement?.contains(e.relatedTarget as HTMLElement)) return;
    node.classList.remove("clicked");
  }

  switch (settings.type) {
    case "bar":
      new BarFocusIndicator({ target: node });
      break;
    case "underline":
      new UnderlineFocusIndicator({ target: node });
      break;
  }


  node.addEventListener("mousedown", mouseDown);
  node.addEventListener("focusout", focusOut);

  return {
    destroy() {
      node.removeEventListener("mousedown", mouseDown);
      node.removeEventListener("focusout", focusOut);
    }
  }
}
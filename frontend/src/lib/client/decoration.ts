import BarFocusIndicator from "../../components/decoration/BarFocusIndicator.svelte";
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

export const barFocusIndicator = (node: HTMLElement, ignore: string | null = null) =>{
  const mouseDown = (e: MouseEvent) => {
    node.classList.add("clicked");
  }

  const focusOut = (e: FocusEvent) => {
    if (ignore == null || !(e.relatedTarget as HTMLElement).classList.contains(ignore)) node.classList.remove("clicked");
  }

  new BarFocusIndicator({ target: node });

  node.addEventListener("mousedown", mouseDown);
  node.addEventListener("focusout", focusOut);

  return {
    destroy() {
      node.removeEventListener("mousedown", mouseDown);
      node.removeEventListener("focusout", focusOut);
    }
  }
}
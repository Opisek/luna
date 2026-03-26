import { mount } from "svelte";

import BarFocusIndicator from "../../components/decoration/focus/BarFocusIndicator.svelte";
import Ripple from "../../components/decoration/Ripple.svelte"
import UnderlineFocusIndicator from "../../components/decoration/focus/UnderlineFocusIndicator.svelte";
import { isDescendentOf } from "../common/misc";

export const addRipple = (e: MouseEvent, addToParent: boolean = true) => {
  if (!e.target) return;

  const parent = e.target as HTMLElement;

  mount(Ripple, { target: addToParent ? parent : e.target as HTMLElement, props: { event: e, parent: parent } });
}

export const focusIndicator = (node: HTMLElement, settings: FocusIndicatorSettings = { type: "bar" }) => {
  let clicked = 0;
  const mouseDown = () => {
    node.classList.add("clicked");
    clicked = (new Date()).getTime();
  }

  const focusOut = (e: FocusEvent) => {
    if ((new Date()).getTime() - clicked < 100) return;
    const original = e.relatedTarget as HTMLElement;
    const current = e.target as HTMLElement;
    if (!original) return;
    if (current instanceof HTMLInputElement && isDescendentOf(original, current)) return; // for select buttons
    node.classList.remove("clicked");
  }

  switch (settings.type) {
    case "bar":
      console.log("mount")
      mount(BarFocusIndicator, { target: node });
      break;
    case "underline":
      mount(UnderlineFocusIndicator, { target: node });
      break;
    case "custom":
      break;
  }


  node.addEventListener("mousedown", mouseDown);
  node.addEventListener("focusout", focusOut);

  return {
    destroy() {
      node.removeEventListener("mousedown", mouseDown);
      node.addEventListener("focusout", focusOut);
    }
  }
}
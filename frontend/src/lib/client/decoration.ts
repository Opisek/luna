import { mount } from "svelte";

import BarFocusIndicator from "../../components/decoration/focus/BarFocusIndicator.svelte";
import Ripple from "../../components/decoration/Ripple.svelte"
import UnderlineFocusIndicator from "../../components/decoration/focus/UnderlineFocusIndicator.svelte";
import { isChildOfModal, isDescendentOf } from "$lib/common/misc";

export const addRipple = (e: MouseEvent, addToParent: boolean = true) => {
  if (!e.target) return;

  const parent = e.target as HTMLElement;

  mount(Ripple, { target: addToParent ? parent : e.target as HTMLElement, props: { event: e, parent: parent } });
}

export const focusIndicator = (node: HTMLElement, settings: FocusIndicatorSettings = { type: "bar" }) => {
  const mouseDown = () => {
    node.classList.add("clicked");
  }

  const focusOut = (e: FocusEvent) => {
    if (!e.relatedTarget) return;
    if (settings.ignoreParent && e.relatedTarget && node.parentElement?.contains(e.relatedTarget as HTMLElement)) return;
    if (e.target instanceof HTMLInputElement && e.relatedTarget && isDescendentOf(e.relatedTarget as HTMLElement, e.target.parentElement as HTMLElement)) return; // for select buttons
    if (!isChildOfModal(e.target as HTMLElement) && isChildOfModal(e.relatedTarget as HTMLElement)) return; // buttons that open modals
    node.classList.remove("clicked");
  }

  switch (settings.type) {
    case "bar":
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
      node.removeEventListener("focusout", focusOut);
    }
  }
}
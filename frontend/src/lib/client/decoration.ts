import { mount } from "svelte";

import BarFocusIndicator from "../../components/decoration/focus/BarFocusIndicator.svelte";
import Ripple from "../../components/decoration/Ripple.svelte"
import UnderlineFocusIndicator from "../../components/decoration/focus/UnderlineFocusIndicator.svelte";
import { isChildOfModal, isDescendentOf, parentModal } from "$lib/common/misc";

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
    const original = e.relatedTarget as HTMLElement;
    const current = e.target as HTMLElement;

    if (!original) return;
    if (node.parentElement?.contains(original as HTMLElement)) return; // Focus within the same element
    if (current instanceof HTMLInputElement && isDescendentOf(original, current)) return; // for select buttons
    if (isChildOfModal(original as HTMLElement) && parentModal(current as HTMLElement) != parentModal(original as HTMLElement) && (original instanceof HTMLButtonElement || original.ariaRoleDescription != null)) return; // Buttons that open modals
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
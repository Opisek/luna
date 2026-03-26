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
  let outTimeout: ReturnType<typeof setTimeout> | null = null;

  const click = () => {
    console.log("click")
    node.classList.add("clicked");
    clicked = (new Date()).getTime();
  }

  const focusIn = () => {
    console.log("in")
    if (!outTimeout) return;
    clearTimeout(outTimeout);
    outTimeout = null;
  }

  const focusOut = () => {
    console.log("out")
    if ((new Date()).getTime() - clicked < 100) return;
    outTimeout = setTimeout(() => {
      node.classList.remove("clicked");
    }, 100);
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


  node.addEventListener("mousedown", click);
  node.addEventListener("click", click);
  node.addEventListener("focusout", focusOut);
  node.addEventListener("focusin", focusIn);

  return {
    destroy() {
      node.addEventListener("mousedown", click);
      node.removeEventListener("click", click);
      node.addEventListener("focusout", focusOut);
      node.addEventListener("focusin", focusIn);
    }
  }
}
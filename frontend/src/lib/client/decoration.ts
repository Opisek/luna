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
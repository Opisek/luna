import { getContext } from "svelte";
import { cubicOut } from "svelte/easing";

export function svelteFlyInHorizontal(node: Node, { duration }: { duration: number }) {
  const context = getContext("flyDirection");
  const direction = typeof context === "function" ? context() : "left";
  return {
    duration: duration,
    easing: cubicOut,
    css: (t: number) => `transform: translateX(${(100 - 100 * t) * (direction === "left" ? 1 : -1)}%);`
  }
}

export function svelteFlyOutHorizontal(node: Node, { duration }: { duration: number }) {
  const context = getContext("flyDirection");
  const direction = typeof context === "function" ? context() : "left";
  return {
    duration: duration,
    easing: cubicOut,
    css: (t: number) => `transform: translateX(${(100 - 100 * t) * (direction === "left" ? -1 : 1)}%);`
  }
}
import { getContext } from "svelte";
import { cubicOut } from "svelte/easing";

export function svelteFlyInHorizontal(node: Node, { duration, flyDirection }: { duration: number, flyDirection: () => string }) {
  const direction = flyDirection();
  return {
    duration: duration,
    easing: cubicOut,
    css: (t: number) => `transform: translateX(${(100 - 100 * t) * (direction === "left" ? 1 : -1)}%);`
  }
}

export function svelteFlyOutHorizontal(node: Node, { duration, flyDirection }: { duration: number, flyDirection: () => string }) {
  const direction = flyDirection();
  return {
    duration: duration,
    easing: cubicOut,
    css: (t: number) => `transform: translateX(${(100 - 100 * t) * (direction === "left" ? -1 : 1)}%);`
  }
}
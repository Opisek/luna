import { browser } from "$app/environment";
import { cubicOut } from "svelte/easing";

export function svelteFlyInHorizontal(node: Node, { duration, flyDirection }: { duration: number, flyDirection: () => string }) {
  const direction = flyDirection();
  return {
    duration: prefersReducedMotion() ? 0 : duration,
    easing: cubicOut,
    css: (t: number) => `transform: translateX(${(100 - 100 * t) * (direction === "left" ? 1 : -1)}%);`
  }
}

export function svelteFlyOutHorizontal(node: Node, { duration, flyDirection }: { duration: number, flyDirection: () => string }) {
  const direction = flyDirection();
  return {
    duration: prefersReducedMotion() ? 0 : duration,
    easing: cubicOut,
    css: (t: number) => `transform: translateX(${(100 - 100 * t) * (direction === "left" ? -1 : 1)}%);`
  }
}

// https://alvin.codes/snippets/sveltekit-reduced-motion/
const preferReducedMotionMediaQuery = "(prefers-reduced-motion: reduce)";
export const prefersReducedMotion = () => browser && window.matchMedia(preferReducedMotionMediaQuery).matches;
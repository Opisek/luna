import type { ColorKeys } from "./colors";

export type Option<T> = {
  value: T;
  name: string;
  icon?: any;
  color?: ColorKeys;
}
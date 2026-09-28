import type { ColorKeys } from "./colors";

export type NotificationModel = {
  created: Date;
  message: string;
  details: string;
  count: number;
  color: ColorKeys; 
  disappear: boolean;
  remove: () => void;
};
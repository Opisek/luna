type NotificationModel = {
  created: Date;
  message: string;
  details: string;
  count: number;
  type: "info" | "success" | "failure";
  disappear: boolean;
  remove: () => void;
};
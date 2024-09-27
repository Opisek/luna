type NotificationModel = {
  created: Date;
  message: string;
  type: "info" | "success" | "failure";
  disappear: boolean;
  remove: () => void;
};
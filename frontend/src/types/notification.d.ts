type NotificationModel = {
  created: Date;
  message: string;
  details: string;
  type: "info" | "success" | "failure";
  disappear: boolean;
  remove: () => void;
};
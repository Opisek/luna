import { writable } from "svelte/store";

export const notificationExpireTime = 5000;

export const notifications = writable([] as NotificationModel[]);
export const notificationCount = writable(0);

export const queueNotification = (type: "info" | "success" | "failure", message: string) => {
  const notification = {
    created: new Date(),
    message: message,
    type: type,
    disappear: false,
    remove: () => {
      if (notification.disappear) return;

      notification.disappear = true;
      notifications.update((notifications) => notifications.map((n) => n.created === notification.created ? notification : n));
      setTimeout(() => {
        notificationCount.update((count) => count - 1);
        notifications.update((notifications) => notifications.filter((n) => n !== notification));
      }, 250);
    }
  }

  notifications.update((notifications) => [...notifications, notification]);
  setTimeout(() => {
    notificationCount.update((count) => count + 1);
  })

  setTimeout(() => {
    notification.remove();
  }, notificationExpireTime);
};
import { writable } from "svelte/store";

export const notificationExpireTime = 5000;

export const notifications = writable([] as NotificationModel[]);
export const notificationCount = writable(0);

let queue = [] as { type: "info" | "success" | "failure", message: string }[];

export const queueNotification = (type: "info" | "success" | "failure", message: string) => {
  queue.push({ type, message });
  if (queue.length === 1) {
    setTimeout(showNotification, 10);
  }
};

function showNotification() {
  const nextNotification = queue.shift();
  if (!nextNotification) return;

  const notification = {
    created: new Date(),
    message: nextNotification.message,
    type: nextNotification.type,
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
  }, 10)

  setTimeout(() => {
    notification.remove();
  }, notificationExpireTime);

  if (queue.length != 0) {
    setTimeout(() => {
      showNotification();
    }, 500);
  }
}

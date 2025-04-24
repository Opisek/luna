import { writable } from "svelte/store";
import type { ColorKeys } from "../../types/colors";
import type { NotificationModel } from "../../types/notification";

export const notificationExpireTime = 5000;

export const notifications = writable([] as NotificationModel[]);

let queue = [] as { color: ColorKeys, message: string, details: string }[];

export const queueNotification = (color: ColorKeys, message: string, details = "") => {
  queue.push({ color, message, details });
  if (queue.length === 1) {
    setTimeout(showNotification, 10);
  }
};

function showNotification() {
  const nextNotification = queue.shift();
  if (!nextNotification) return;

  notifications.update((currentNotifications) => {
    const duplicate = currentNotifications.filter(n =>
      n.message == nextNotification.message &&
      n.details == nextNotification.details &&
      n.color == nextNotification.color &&
      n.disappear == false
    );

    // Instead of showing the same notification multiple times, we increment a
    // counter if a notification with the same content is found.
    if (duplicate.length > 0) {
      duplicate[0].count++;
      return currentNotifications;
    } else {
      const notification = {
        created: new Date(),
        message: nextNotification.message,
        details: nextNotification.details,
        count: 1,
        color: nextNotification.color,
        disappear: false,
        remove: () => {
          if (notification.disappear) return;

          notification.disappear = true;
          notifications.update((notifications) => notifications.map((n) => n.created === notification.created ? notification : n));
          setTimeout(() => {
            notifications.update((notifications) => notifications.filter((n) => n !== notification));
          }, 250); // wait for the disappear animation to finish
        }
      }

      return [...currentNotifications, notification];
    }
  });

  if (queue.length != 0) {
    setTimeout(() => {
      showNotification();
    }, 500);
  }
}

export function redrawNotifications() {
  notifications.update((notifications) => notifications);
}

export function resetNotifications() {
  notifications.update(() => []);
}
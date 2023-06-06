import { isMobileSafariIos } from './helpers';

export function arePushNotificationSupported() {
  return 'Notification' in window && 'serviceWorker' in navigator && 'PushManager' in window;
}

export function canPushNotificationsBeSupported() {
  // Mobile safari will return false for supporting push notifications, but
  // it does support them. So we need to check for mobile safari and return
  // true if it is.
  if (isMobileSafariIos()) {
    return true;
  }
  return 'Notification' in window && 'serviceWorker' in navigator && 'PushManager' in window;
}

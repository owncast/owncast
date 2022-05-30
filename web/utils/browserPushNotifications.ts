export default function isPushNotificationSupported() {
  return 'serviceWorker' in navigator && 'PushManager' in window;
}

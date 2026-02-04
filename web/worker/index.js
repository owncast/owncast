/* eslint-disable no-restricted-globals */
/* eslint-disable no-underscore-dangle */
/* eslint-disable no-undef */

// Custom service worker code that next-pwa will bundle into sw.js.
// This handles push notifications and provides a way to pause precaching
// when notification registration needs priority.

// Handle messages from main thread
self.addEventListener('message', event => {
  if (event.data?.type === 'PAUSE_PRECACHING') {
    console.log('Pausing precaching for notification registration');

    // Signal to Workbox to stop precaching by clearing the queue
    // This works because Workbox checks this before each fetch
    if (self.__WB_MANIFEST) {
      // Clear the manifest to prevent further precaching
      self.__WB_MANIFEST.length = 0;
    }

    // Respond that we're ready
    if (event.ports[0]) {
      event.ports[0].postMessage({ paused: true });
    }
  }
});

// Push notification handler
self.addEventListener('push', event => {
  const data = JSON.parse(event.data.text());
  const { title, body, icon, tag } = data;
  const options = {
    title: title || 'Live!',
    body: body || 'This live stream has started.',
    icon: icon || '/logo/external',
    tag,
  };

  event.waitUntil(self.registration.showNotification(options.title, options));
});

// Handle notification click
self.addEventListener('notificationclick', event => {
  event.notification.close();
  event.waitUntil(clients.openWindow('/'));
});

self.addEventListener('activate', (event) => {
  console.log('Owncast service worker activated', event);
});

self.addEventListener('install', (event) => {
  console.log('installing Owncast service worker...', event);
});

self.addEventListener('push', (event) => {
  console.log('Received push event', event);

  const data = JSON.parse(event.data.text());
  const { title, body, icon, tag } = data;
  const options = {
    title: title || 'Live!',
    body: body || 'This live stream has started.',
    icon: icon || '/logo/external',
    tag: tag,
  };

  event.waitUntil(self.registration.showNotification(options.title, options));
});

self.addEventListener('notificationclick', (event) => {
  let notification = event.notification;
  console.log(notification);
  clients.openWindow('/');
});

/* eslint-disable no-restricted-globals */
// Dev-only self-destroying service worker. `self` is the service worker global
// scope here, not the window, so the no-restricted-globals rule is disabled.
//
// In development the dev server is configured (see next.config.js) to serve
// THIS file at /sw.js via a beforeFiles rewrite, instead of the production
// service worker that lives at public/sw.js. next-pwa is disabled in dev so it
// never registers a worker, but a browser may still hold a STALE production
// worker registered against a dev origin (e.g. localhost:3001) from an earlier
// run. That stale worker intercepts /_next/static/chunks/* and serves cached
// assets whose hashes no longer match the running dev build, leaving the page
// blank. Serving this self-destroying worker makes the browser's normal update
// check replace the stale worker with one that unregisters itself, clears all
// caches, and reloads open tabs -- auto-healing the dev origin with no manual
// "clear site data" step. It has no effect on production builds, which never
// rewrite /sw.js and continue to use the real public/sw.js.

self.addEventListener('install', () => {
  self.skipWaiting();
});

self.addEventListener('activate', event => {
  event.waitUntil(
    (async () => {
      try {
        const keys = await caches.keys();
        await Promise.all(keys.map(key => caches.delete(key)));
      } catch {
        // Best effort; unregister regardless.
      }
      await self.registration.unregister();
      const clients = await self.clients.matchAll({ type: 'window' });
      clients.forEach(client => client.navigate(client.url));
    })(),
  );
});

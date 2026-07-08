/* eslint-disable no-console */
// Generates the production service worker over the static export in out/,
// run by `npm run build` after `next build`. This replaces next-pwa, which
// was unmaintained and webpack-only (the build now uses Turbopack).
//
// The worker precaches the viewer-facing assets for fast repeat loads and
// PWA installability. Excluded, matching the old next-pwa configuration:
// admin page chunks (authenticated, per-session), platform logos and admin
// styles. There is no runtime caching, same as before.

const fs = require('fs');
const path = require('path');
const { generateSW } = require('workbox-build');

const outDir = path.join(__dirname, '..', 'out');

async function run() {
  if (!fs.existsSync(path.join(outDir, 'index.html'))) {
    throw new Error(`[generate-sw] no static export found at ${outDir}, run next build first`);
  }

  const { count, size, warnings } = await generateSW({
    globDirectory: outDir,
    globPatterns: [
      '_next/static/**/*',
      'fonts/**/*',
      'img/**/*',
      'styles/**/*',
      'manifest.json',
      'serviceWorker.js',
      '*.png',
    ],
    globIgnores: [
      '**/*.map',
      '_next/static/chunks/pages/admin/**',
      '_next/static/chunks/pages/admin*',
      'img/platformlogos/**',
      'styles/admin/**',
    ],
    swDest: path.join(outDir, 'sw.js'),
    skipWaiting: true,
    clientsClaim: true,
    // Manifest URLs are generated relative to out/, make them absolute.
    modifyURLPrefix: { '': '/' },
    // _next/static files are already content-hashed, no revision busting.
    dontCacheBustURLsMatching: /^\/_next\/static\//,
  });

  warnings.forEach(w => console.warn('[generate-sw]', w));

  // Fail the build loudly rather than ship a broken or leaky manifest.
  const sw = fs.readFileSync(path.join(outDir, 'sw.js'), 'utf8');
  const manifest = sw.match(/precacheAndRoute\(\[.*?\]/s)?.[0] ?? '';
  if (count < 50) {
    throw new Error(`[generate-sw] suspiciously small precache manifest (${count} entries)`);
  }
  if (/admin/i.test(manifest)) {
    throw new Error('[generate-sw] admin assets leaked into the precache manifest');
  }
  console.log(`[generate-sw] precached ${count} files, ${(size / 1024 / 1024).toFixed(1)} MB`);
}

run().catch(err => {
  console.error(err);
  process.exit(1);
});

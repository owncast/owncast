const withBundleAnalyzer = require('@next/bundle-analyzer');
const { PHASE_DEVELOPMENT_SERVER } = require('next/constants');

// The dev server proxies API/asset requests to a running Owncast backend to
// work around CORS. The target defaults to localhost:8080 but can be pointed
// at any backend via OWNCAST_DEV_BACKEND, which lets you run a second dev
// server against a second instance (e.g. localhost:8081) to test federation
// between two instances without rebuilding the embedded web bundle.
const BACKEND = process.env.OWNCAST_DEV_BACKEND || 'http://localhost:8080';

async function rewrites() {
  return {
    // beforeFiles runs before the filesystem, so this takes precedence over
    // the static public/sw.js. In dev we serve a self-destroying worker at
    // /sw.js to neutralize any stale production service worker a browser may
    // have registered on a dev port (which otherwise serves cached chunks
    // that don't match the running dev build and blanks the page). Dev-only:
    // these rewrites are only attached in PHASE_DEVELOPMENT_SERVER below.
    beforeFiles: [
      {
        source: '/sw.js',
        destination: '/dev-sw.js',
      },
    ],
    afterFiles: [
      {
        source: '/api/:path*',
        destination: `${BACKEND}/api/:path*`,
      },
      // Plugin admin iframes proxied so they're same-origin to the admin UI
      // in dev. Two patterns: the first matches slash-terminated URLs and
      // preserves the trailing slash through the rewrite; the second handles
      // everything else. Without the slash-preserving variant the backend
      // 301-redirects /plugins/<name>/admin to /plugins/<name>/admin/ to
      // canonicalize the directory, the proxy strips the slash again, and
      // the iframe runs into an infinite redirect loop.
      {
        source: '/plugins/:path*/',
        destination: `${BACKEND}/plugins/:path*/`,
      },
      {
        source: '/plugins/:path*',
        destination: `${BACKEND}/plugins/:path*`,
      },
      {
        source: '/hls/:path*',
        destination: `${BACKEND}/hls/:path*`,
      },
      {
        source: '/img/:path*',
        destination: `${BACKEND}/img/:path*`,
      },
      {
        source: '/logo',
        destination: `${BACKEND}/logo`,
      },
      {
        source: '/thumbnail.jpg',
        destination: `${BACKEND}/thumbnail.jpg`,
      },
      {
        source: '/customjavascript',
        destination: `${BACKEND}/customjavascript`,
      },
      {
        source: '/favicon.ico',
        destination: `${BACKEND}/favicon.ico`,
      },
    ],
  };
}

module.exports = async phase => {
  /**
   * @type {import('next').NextConfig}
   */
  let nextConfig = {
    productionBrowserSourceMaps: process.env.SOURCE_MAPS === 'true',
    distDir: '.next',
    trailingSlash: true,
    reactStrictMode: true,
    devIndicators: false,
    images: {
      unoptimized: true,
    },
    transpilePackages: ['antd', '@rc-component', '@ant-design'],
    // Component .svg imports compile to React components via svgr, as the
    // old custom webpack rule did before Turbopack became the bundler.
    turbopack: {
      // The repo root has its own package-lock (commitlint), pin the app
      // root so Turbopack doesn't infer the wrong workspace directory.
      root: __dirname,
      rules: {
        '*.svg': {
          loaders: ['@svgr/webpack'],
          as: '*.js',
        },
      },
    },
    pageExtensions: ['tsx'],
  };

  // The bundle analyzer works by injecting a webpack() config key, which
  // fails Turbopack builds. Attach it only when analyzing, and run those
  // builds through webpack: ANALYZE=true npx next build --webpack
  if (process.env.ANALYZE === 'true') {
    nextConfig = withBundleAnalyzer({ enabled: true })(nextConfig);
  }

  if (phase === PHASE_DEVELOPMENT_SERVER) {
    nextConfig = {
      ...nextConfig,
      rewrites,
      // Isolate the dev build cache so a second dev server (pointed at
      // another backend via OWNCAST_DEV_BACKEND) doesn't fight the first
      // over .next. Dev-phase only: production builds always emit to
      // .next, which the css inliner in pages/_document.tsx relies on.
      distDir: process.env.OWNCAST_DEV_DISTDIR || '.next',
    };
  } else {
    // The production build is a fully static export served by the Go
    // backend. The service worker is generated over it afterwards by
    // build-scripts/generate-sw.js (chained in the package.json build
    // script), replacing the webpack-only next-pwa plugin.
    nextConfig = {
      ...nextConfig,
      output: 'export',
    };
  }
  return nextConfig;
};

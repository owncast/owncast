const withLess = require('next-with-less');
const withBundleAnalyzer = require('@next/bundle-analyzer')({
  enabled: process.env.ANALYZE === 'true',
});
const { PHASE_DEVELOPMENT_SERVER } = require('next/constants');

const withPWA = require('next-pwa')({
  dest: 'public',
  runtimeCaching: [],
  register: true,
  skipWaiting: true,
  disableDevLogs: true,
  publicExcludes: ['!img/platformlogos/**/*', '!styles/admin/**/*'],
  buildExcludes: [/chunks\/pages\/admin.*/, '!**/admin/**/*'],
  sourcemap: process.env.NODE_ENV === 'development',
  disable: process.env.NODE_ENV === 'development',
});

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
  let nextConfig = withPWA(
    withBundleAnalyzer(
      withLess({
        productionBrowserSourceMaps: process.env.SOURCE_MAPS === 'true',
        // Isolate the build cache so a second dev server (pointed at another
        // backend via OWNCAST_DEV_BACKEND) doesn't fight the first over .next.
        // Defaults to .next, so normal dev and production builds are unchanged.
        distDir: process.env.OWNCAST_DEV_DISTDIR || '.next',
        trailingSlash: true,
        reactStrictMode: true,
        eslint: {
          ignoreDuringBuilds: true,
        },
        images: {
          unoptimized: true,
        },
        swcMinify: true,
        transpilePackages: [
          'antd',
          '@ant-design',
          'rc-util',
          'rc-pagination',
          'rc-picker',
          'rc-notification',
          'rc-tooltip',
          'rc-tree',
          'rc-table',
        ],
        webpack(config) {
          config.module.rules.push({
            test: /\.svg$/i,
            issuer: /\.[jt]sx?$/,
            use: ['@svgr/webpack'],
          });

          return config;
        },
        pageExtensions: ['tsx'],
      }),
    ),
  );

  if (phase === PHASE_DEVELOPMENT_SERVER) {
    nextConfig = {
      ...nextConfig,
      rewrites,
    };
  } else {
    nextConfig = {
      ...nextConfig,
      output: 'export',
    };
  }
  return nextConfig;
};

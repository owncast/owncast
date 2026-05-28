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

async function rewrites() {
  return [
    {
      source: '/api/:path*',
      destination: 'http://localhost:8080/api/:path*', // Proxy to Backend to work around CORS.
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
      destination: 'http://localhost:8080/plugins/:path*/',
    },
    {
      source: '/plugins/:path*',
      destination: 'http://localhost:8080/plugins/:path*',
    },
    {
      source: '/hls/:path*',
      destination: 'http://localhost:8080/hls/:path*', // Proxy to Backend to work around CORS.
    },
    {
      source: '/img/:path*',
      destination: 'http://localhost:8080/img/:path*', // Proxy to Backend to work around CORS.
    },
    {
      source: '/logo',
      destination: 'http://localhost:8080/logo', // Proxy to Backend to work around CORS.
    },
    {
      source: '/thumbnail.jpg',
      destination: 'http://localhost:8080/thumbnail.jpg', // Proxy to Backend to work around CORS.
    },
    {
      source: '/customjavascript',
      destination: 'http://localhost:8080/customjavascript', // Proxy to Backend to work around CORS.
    },
    {
      source: '/favicon.ico',
      destination: 'http://localhost:8080/favicon.ico', // Proxy to Backend to work around CORS.
    },
  ];
}

module.exports = async phase => {
  /**
   * @type {import('next').NextConfig}
   */
  let nextConfig = withPWA(
    withBundleAnalyzer(
      withLess({
        productionBrowserSourceMaps: process.env.SOURCE_MAPS === 'true',
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

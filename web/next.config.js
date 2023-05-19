const withLess = require('next-with-less');
const withBundleAnalyzer = require('@next/bundle-analyzer')({
  enabled: process.env.ANALYZE === 'true',
});
const withPWA = require('next-pwa')({
  dest: 'public',
  register: true,
  skipWaiting: true,
  publicExcludes: ['!img/platformlogos/**/*', '!styles/admin/**/*'],
  buildExcludes: [/chunks\/pages\/admin.*/, '!**/admin/**/*'],
  sourcemap: process.env.NODE_ENV === 'development',
});

module.exports = withPWA(
  withBundleAnalyzer(
    withLess({
      productionBrowserSourceMaps: process.env.SOURCE_MAPS === 'true',
      trailingSlash: true,
      reactStrictMode: true,
      images: {
        unoptimized: true,
      },
      swcMinify: true,
      webpack(config) {
        config.module.rules.push({
          test: /\.svg$/i,
          issuer: /\.[jt]sx?$/,
          use: ['@svgr/webpack'],
        });

        return config;
      },
      async rewrites() {
        return [
          {
            source: '/api/:path*',
            destination: 'http://localhost:8080/api/:path*', // Proxy to Backend to work around CORS.
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
        ];
      },
      pageExtensions: ['tsx'],
    }),
  ),
);

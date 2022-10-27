const withLess = require('next-with-less');

module.exports = withLess({
  trailingSlash: true,
  reactStrictMode: true,
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
    ];
  },
  pageExtensions: ['tsx'],
});

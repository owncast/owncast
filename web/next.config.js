const withLess = require('next-with-less');

module.exports = withLess({
  trailingSlash: true,
  async rewrites() {
    return [
      {
        source: '/api/:path*',
        destination: 'http://localhost:8080/api/:path*', // Proxy to Backend to work around CORS.
      },
    ];
  },
});

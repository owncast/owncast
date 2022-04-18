const withLess = require('next-with-less');

module.exports = withLess({
  basePath: '/admin',
  trailingSlash: true,
});

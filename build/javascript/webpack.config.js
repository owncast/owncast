// vendor-bundles.webpack.config.js
var webpack = require('webpack');

const packageJSON = require('./package.json');
var packages = [];

for (var key in packageJSON.dependencies) {
  packages.push(key);
}

module.exports = {
  mode: 'production',
  entry: {
    // create two library bundles, one with jQuery and
    // another with Angular and related libraries
    vendor: packages,
  },

  output: {
    filename: '[name].bundle.js',
    path: __dirname + '/dist',

    // The name of the global variable which the library's
    // require() function will be assigned to
    library: '[name]_lib',
  },

  plugins: [
    new webpack.DllPlugin({
      // The path to the manifest file which maps between
      // modules included in a bundle and the internal IDs
      // within that bundle
      path: 'dist/[name]-manifest.json',
      // The name of the global variable which the library's
      // require function has been assigned to. This must match the
      // output.library option above
      name: '[name]_lib',
    }),
  ],
};

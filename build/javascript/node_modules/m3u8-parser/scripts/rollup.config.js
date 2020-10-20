const generate = require('videojs-generate-rollup-config');

// see https://github.com/videojs/videojs-generate-rollup-config
// for options
const options = {
  input: 'src/index.js',
  externals(defaults) {
    defaults.module.push('@videojs/vhs-utils');

    return defaults;
  }
};
const config = generate(options);

if (config.builds.test) {
  config.builds.test.output[0].format = 'umd';
}

// Add additonal builds/customization here!

// export the builds to rollup
export default Object.values(config.builds);

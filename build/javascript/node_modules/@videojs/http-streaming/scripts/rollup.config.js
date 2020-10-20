const generate = require('videojs-generate-rollup-config');
const worker = require('@gkatsev/rollup-plugin-bundle-worker');
const {terser} = require('rollup-plugin-terser');
const createTestData = require('./create-test-data.js');

// see https://github.com/videojs/videojs-generate-rollup-config
// for options
const options = {
  input: 'src/videojs-http-streaming.js',
  distName: 'videojs-http-streaming',
  globals(defaults) {
    defaults.browser.xmldom = 'window';
    defaults.test.xmldom = 'window';
    return defaults;
  },
  externals(defaults) {
    return Object.assign(defaults, {
      module: defaults.module.concat([
        'aes-decrypter',
        'm3u8-parser',
        'mpd-parser',
        'mux.js',
        '@videojs/vhs-utils'
      ])
    });
  },
  plugins(defaults) {
    defaults.module.splice(2, 0, 'worker');
    defaults.browser.splice(2, 0, 'worker');
    defaults.test.splice(3, 0, 'worker');
    defaults.test.splice(0, 0, 'createTestData');

    // istanbul is only in the list for regular builds and not watch
    if (defaults.test.indexOf('istanbul') !== -1) {
      defaults.test.splice(defaults.test.indexOf('istanbul'), 1);
    }

    return defaults;
  },
  primedPlugins(defaults) {
    return Object.assign(defaults, {
      worker: worker(),
      uglify: terser({
        output: {comments: 'some'},
        compress: {passes: 2},
        include: [/^.+\.min\.js$/]
      }),
      createTestData: createTestData()
    });
  },
  babel(defaults) {
    const presetEnvSettings = defaults.presets[0][1];

    presetEnvSettings.exclude = presetEnvSettings.exclude || [];
    presetEnvSettings.exclude.push('@babel/plugin-transform-typeof-symbol');

    return defaults;
  }
};
const config = generate(options);

// Add additonal builds/customization here!

// export the builds to rollup
export default [
  config.makeBuild('browser', {
    input: 'src/decrypter-worker.js',
    output: {
      format: 'iife',
      name: 'decrypterWorker',
      file: 'src/decrypter-worker.worker.js'
    },
    external: []
  }),

  config.makeBuild('browser', {
    input: 'src/transmuxer-worker.js',
    output: {
      format: 'iife',
      name: 'transmuxerWorker',
      file: 'src/transmuxer-worker.worker.js'
    },
    external: []
  })
].concat(Object.values(config.builds));

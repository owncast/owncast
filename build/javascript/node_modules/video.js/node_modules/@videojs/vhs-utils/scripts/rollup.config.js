const generate = require('videojs-generate-rollup-config');
const fs = require('fs');
const path = require('path');

const BASE_DIR = path.join(__dirname, '..');
const SRC_DIR = path.join(BASE_DIR, 'src');

const files = fs.readdirSync(SRC_DIR);

const shared = {
  externals(defaults) {
    defaults.module.push('url-toolkit');
    return defaults;
  }
};
const builds = [];

files.forEach(function(file, i) {
  const config = generate(Object.assign({}, shared, {
    input: path.relative(BASE_DIR, path.join(SRC_DIR, file)),
    distName: path.basename(file, path.extname(file))
  }));

  // gaurd against test only builds
  if (config.builds.module) {
    const module = config.builds.module;

    module.output = module.output.filter((o) => o.format === 'cjs');
    module.output[0].file = module.output[0].file.replace('.cjs.js', '.js');
    builds.push(module);
  }

  // gaurd against production only builds
  // only add the last test bundle we generate as they are all the same
  if (i === (files.length - 1) && config.builds.test) {
    builds.push(config.builds.test);
  }
});

// export the builds to rollup

export default builds;

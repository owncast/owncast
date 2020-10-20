'use strict';
/* eslint no-console: 0 */

const fs = require('fs');
const path = require('path');

const basePath = path.resolve(__dirname, '..');
const testDataDir = path.join(basePath, 'test');
const manifestDir = path.join(basePath, 'test', 'fixtures', 'm3u8');
const manifestFilepath = path.join(testDataDir, 'dist', 'test-manifests.js');
const expectedFilepath = path.join(testDataDir, 'dist', 'test-expected.js');

const build = function() {
  let manifests = 'export default {\n';
  let expected = 'export default {\n';

  const files = fs.readdirSync(manifestDir);

  while (files.length > 0) {
    const file = path.resolve(manifestDir, files.shift());
    const extname = path.extname(file);

    if (extname === '.m3u8') {
      // translate this manifest
      manifests += '  \'' + path.basename(file, '.m3u8') + '\': ';
      manifests += fs.readFileSync(file, 'utf8')
        .split(/\r\n|\n/)
        // quote and concatenate
        .map(function(line) {
          return '    \'' + line + '\\n\' +\n';
        }).join('')
        // strip leading spaces and the trailing '+'
        .slice(4, -3);
      manifests += ',\n';
    } else if (extname === '.json') {
      // append the expected parse
      expected += '  "' + path.basename(file, '.json') + '": ';
      expected += fs.readFileSync(file, 'utf8');
      expected += ',\n';
    } else {
      console.log('Unknown file ' + file + ' found in manifest dir ' + manifestDir);
    }

  }

  // clean up and close the objects
  manifests = manifests.slice(0, -2);
  manifests += '\n};\n';
  expected = expected.slice(0, -2);
  expected += '\n};\n';

  fs.writeFileSync(manifestFilepath, manifests);
  fs.writeFileSync(expectedFilepath, expected);
  console.log('Wrote test data file ' + manifestFilepath);
  console.log('Wrote test data file ' + expectedFilepath);
};

const watch = function() {
  build();
  fs.watch(manifestDir, function(event, filename) {
    console.log('files in manifest dir were changed rebuilding manifest data');
    build();
  });
};

const clean = function() {
  try {
    fs.unlinkSync(manifestFilepath);
  } catch (e) {
    console.log(e);
  }
  try {
    fs.unlinkSync(expectedFilepath);
  } catch (e) {
    console.log(e);
  }
};

module.exports = {
  build,
  watch,
  clean
};

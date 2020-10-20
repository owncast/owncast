const { dirname, resolve } = require("path");
const { existsSync } = require("fs");

module.exports = function resolveId(file, url, opts) {
  if (opts.paths && opts.paths.length) {
    let absolutePath;

    for (let path of opts.paths) {
      absolutePath = resolve(path, url);

      if (existsSync(absolutePath)) {
        return absolutePath;
      }
    }

    return absolutePath;
  }

  if (file) {
    return resolve(dirname(file), url);
  }

  return resolve(url);
};

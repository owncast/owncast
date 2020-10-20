const execSync = require('child_process').execSync;
const path = require('path');
const semver = require('semver');
const pkg = require('../package.json');

if (!semver.prerelease(pkg.version)) {
  process.chdir(path.resolve(__dirname, '..'));
  execSync('conventional-changelog -p videojs -i CHANGELOG.md -s');
  execSync('git add CHANGELOG.md');
}

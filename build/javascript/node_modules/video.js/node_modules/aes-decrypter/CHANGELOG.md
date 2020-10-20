<a name="3.0.0"></a>
# [3.0.0](https://github.com/videojs/aes-decrypter/compare/v2.0.0...v3.0.0) (2017-07-24)

### Features

* Use Rollup for packaging ([bda57ab](https://github.com/videojs/aes-decrypter/commit/bda57ab))

### Chores

* prepare CHANGELOG for new process ([1a5175c](https://github.com/videojs/aes-decrypter/commit/1a5175c))


### BREAKING CHANGES

* revert to 1.x and stop using web crypto.

## 2.0.0 (2016-11-15)
* Use webcrypto for aes-cbc segment decryption when supported (#4)
* Lock the linter to a specific version

## 1.1.1 (2016-11-17)
* version to revert 1.1.0

## 1.0.3 (2016-06-16)
* dont do browserify-shim globally since we only use it in tests (#1)

## 1.0.2 (2016-06-16)
* specify browserify transform globally

## 1.0.1 (2016-06-16)
* fixing the build pipeline

## 1.0.0 (2016-06-16)
* initial


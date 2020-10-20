<a name="4.4.0"></a>
# [4.4.0](https://github.com/videojs/m3u8-parser/compare/v4.3.0...v4.4.0) (2019-06-25)

### Features

* parse key attributes for Widevine HLS ([#88](https://github.com/videojs/m3u8-parser/issues/88)) ([d835fa8](https://github.com/videojs/m3u8-parser/commit/d835fa8))

### Chores

* **package:** update all dev dependencies ([#89](https://github.com/videojs/m3u8-parser/issues/89)) ([e991447](https://github.com/videojs/m3u8-parser/commit/e991447))

<a name="4.3.0"></a>
# [4.3.0](https://github.com/videojs/m3u8-parser/compare/v4.2.0...v4.3.0) (2019-01-10)

### Features

* custom tag mapping ([#73](https://github.com/videojs/m3u8-parser/issues/73)) ([0ef040a](https://github.com/videojs/m3u8-parser/commit/0ef040a))

### Chores

* Update to plugin generator 7 standards ([#53](https://github.com/videojs/m3u8-parser/issues/53)) ([35ff471](https://github.com/videojs/m3u8-parser/commit/35ff471))
* **package:** update rollup to version 0.66.0 ([#55](https://github.com/videojs/m3u8-parser/issues/55)) ([2407466](https://github.com/videojs/m3u8-parser/commit/2407466))
* Update videojs-generate-karma-config to the latest version 🚀 ([#59](https://github.com/videojs/m3u8-parser/issues/59)) ([023c6c9](https://github.com/videojs/m3u8-parser/commit/023c6c9))
* Update videojs-generate-karma-config to the latest version 🚀 ([#60](https://github.com/videojs/m3u8-parser/issues/60)) ([2773819](https://github.com/videojs/m3u8-parser/commit/2773819))
* Update videojs-generate-rollup-config to the latest version 🚀 ([#58](https://github.com/videojs/m3u8-parser/issues/58)) ([8c28a8b](https://github.com/videojs/m3u8-parser/commit/8c28a8b))

<a name="4.2.0"></a>
# [4.2.0](https://github.com/videojs/m3u8-parser/compare/v4.1.0...v4.2.0) (2018-02-23)

### Features

* add program-date-time tag info to parsed segments ([#27](https://github.com/videojs/m3u8-parser/issues/27)) ([44fc6f8](https://github.com/videojs/m3u8-parser/commit/44fc6f8))

<a name="4.1.0"></a>
# [4.1.0](https://github.com/videojs/m3u8-parser/compare/v4.0.0...v4.1.0) (2018-01-24)

<a name="4.0.0"></a>
# [4.0.0](https://github.com/videojs/m3u8-parser/compare/v3.0.0...v4.0.0) (2017-11-21)

### Features

* added ability to parse EXT-X-START tags [#31](https://github.com/videojs/m3u8-parser/pull/31)

### BREAKING CHANGES

* camel case module name in rollup config to work with latest rollup [#32](https://github.com/videojs/m3u8-parser/pull/32)

<a name="3.0.0"></a>
# 3.0.0 (2017-06-09)

### Features

* Rollup ([#24](https://github.com/videojs/m3u8-parser/issues/24)) ([47ef11f](https://github.com/videojs/m3u8-parser/commit/47ef11f))


### BREAKING CHANGES

* drop bower support.

## 2.1.0 (2017-02-23)
* parse FORCED attribute of media-groups [#15](https://github.com/videojs/m3u8-parser/pull/15)
* Pass any CHARACTERISTICS value of a track with the track object [#14](https://github.com/videojs/m3u8-parser/pull/14)

## 2.0.1 (2017-01-20)
* Fix: Include the babel ES3 tranform to support IE8 [#13](https://github.com/videojs/m3u8-parser/pull/13)

## 2.0.0 (2017-01-13)
* Manifest object is now initialized with an empty segments arrays
* moved to latest videojs-standard version, brought code into
compliance with the latest eslint rules.

## 1.0.2 (2016-06-07)
* fix the build pipeline
* removed video.js css/js inclusion during tests

## 1.0.1 (2016-06-07)
* remove dependence on video.js
* added contributors to package.json

## 1.0.0 (2016-06-03)
Initial Release


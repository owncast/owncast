<a name="0.10.0"></a>
# [0.10.0](https://github.com/videojs/mpd-parser/compare/v0.9.0...v0.10.0) (2020-02-04)

### Features

* expose suggestPresentationDelay if the type is dynamic ([#82](https://github.com/videojs/mpd-parser/issues/82)) ([cd27003](https://github.com/videojs/mpd-parser/commit/cd27003))

<a name="0.9.0"></a>
# [0.9.0](https://github.com/videojs/mpd-parser/compare/v0.8.2...v0.9.0) (2019-08-30)

### Features

* node support ([#75](https://github.com/videojs/mpd-parser/issues/75)) ([58b43b0](https://github.com/videojs/mpd-parser/commit/58b43b0))

<a name="0.8.2"></a>
## [0.8.2](https://github.com/videojs/mpd-parser/compare/v0.8.1...v0.8.2) (2019-08-22)

### Chores

* update generator and use [@videojs](https://github.com/videojs)/vhs-utils ([#76](https://github.com/videojs/mpd-parser/issues/76)) ([1238749](https://github.com/videojs/mpd-parser/commit/1238749))

<a name="0.8.1"></a>
## [0.8.1](https://github.com/videojs/mpd-parser/compare/v0.8.0...v0.8.1) (2019-05-01)

### Bug Fixes

* skip playlists without sidx ([#73](https://github.com/videojs/mpd-parser/issues/73)) ([67d2bad](https://github.com/videojs/mpd-parser/commit/67d2bad)), closes [videojs/video.js#5289](https://github.com/videojs/video.js/issues/5289)

<a name="0.8.0"></a>
# [0.8.0](https://github.com/videojs/mpd-parser/compare/v0.7.0...v0.8.0) (2019-04-11)

### Features

* add sidx information to segment base playlists ([#41](https://github.com/videojs/mpd-parser/issues/41)) ([1176109](https://github.com/videojs/mpd-parser/commit/1176109))

### Bug Fixes

* make byteRange.length inclusive ([#43](https://github.com/videojs/mpd-parser/issues/43)) ([28d217a](https://github.com/videojs/mpd-parser/commit/28d217a))

### Chores

* add netlify for testing ([#45](https://github.com/videojs/mpd-parser/issues/45)) ([a78a7be](https://github.com/videojs/mpd-parser/commit/a78a7be))
* Update videojs-generate-karma-config to the latest version 🚀 ([#37](https://github.com/videojs/mpd-parser/issues/37)) ([a18c660](https://github.com/videojs/mpd-parser/commit/a18c660))
* Update videojs-generate-karma-config to the latest version 🚀 ([#38](https://github.com/videojs/mpd-parser/issues/38)) ([3aaabac](https://github.com/videojs/mpd-parser/commit/3aaabac))
* Update videojs-generate-rollup-config to the latest version 🚀 ([#36](https://github.com/videojs/mpd-parser/issues/36)) ([3f6ccbd](https://github.com/videojs/mpd-parser/commit/3f6ccbd))
* **package:** update videojs-generate-karma-config to 5.0.2 ([#54](https://github.com/videojs/mpd-parser/issues/54)) ([fcbabc3](https://github.com/videojs/mpd-parser/commit/fcbabc3))
* **package:** videojs-generate-karma-config[@4](https://github.com/4).0.0 does not exist ([#44](https://github.com/videojs/mpd-parser/issues/44)) ([bc361b5](https://github.com/videojs/mpd-parser/commit/bc361b5))

<a name="0.7.0"></a>
# [0.7.0](https://github.com/videojs/mpd-parser/compare/v0.6.1...v0.7.0) (2018-10-24)

### Features

* limited multiperiod support ([#35](https://github.com/videojs/mpd-parser/issues/35)) ([aee87a0](https://github.com/videojs/mpd-parser/commit/aee87a0))

### Bug Fixes

* fixed segment timeline parsing when duration is present ([#34](https://github.com/videojs/mpd-parser/issues/34)) ([90feb2d](https://github.com/videojs/mpd-parser/commit/90feb2d))
* Remove the postinstall script to prevent install issues ([#29](https://github.com/videojs/mpd-parser/issues/29)) ([ae458f4](https://github.com/videojs/mpd-parser/commit/ae458f4))

### Chores

* Update to generator-videojs-plugin[@7](https://github.com/7).2.0 ([#28](https://github.com/videojs/mpd-parser/issues/28)) ([909cf08](https://github.com/videojs/mpd-parser/commit/909cf08))
* **package:** Update dependencies to enable Greenkeeper 🌴 ([#30](https://github.com/videojs/mpd-parser/issues/30)) ([0593c2c](https://github.com/videojs/mpd-parser/commit/0593c2c))

<a name="0.6.1"></a>
## [0.6.1](https://github.com/videojs/mpd-parser/compare/v0.6.0...v0.6.1) (2018-05-17)

### Bug Fixes

* babel es module ([#25](https://github.com/videojs/mpd-parser/issues/25)) ([9a84461](https://github.com/videojs/mpd-parser/commit/9a84461))

<a name="0.6.0"></a>
# [0.6.0](https://github.com/videojs/mpd-parser/compare/v0.5.0...v0.6.0) (2018-03-30)

### Features

* support in-manifest DRM data ([#23](https://github.com/videojs/mpd-parser/issues/23)) ([7ce9aca](https://github.com/videojs/mpd-parser/commit/7ce9aca))

<a name="0.5.0"></a>
# [0.5.0](https://github.com/videojs/mpd-parser/compare/v0.4.0...v0.5.0) (2018-03-15)

### Features

* live support with SegmentTemplate[@duration](https://github.com/duration) and more ([#22](https://github.com/videojs/mpd-parser/issues/22)) ([f1cee87](https://github.com/videojs/mpd-parser/commit/f1cee87))

<a name="0.4.0"></a>
# [0.4.0](https://github.com/videojs/mpd-parser/compare/v0.3.0...v0.4.0) (2018-02-26)

### Features

* Adding support for segments in Period and Representation. ([#19](https://github.com/videojs/mpd-parser/issues/19)) ([8e59b38](https://github.com/videojs/mpd-parser/commit/8e59b38))

<a name="0.3.0"></a>
# [0.3.0](https://github.com/videojs/mpd-parser/compare/v0.2.1...v0.3.0) (2018-02-06)

### Features

* Parse <SegmentList> and <SegmentBase> ([#18](https://github.com/videojs/mpd-parser/issues/18)) ([71b8976](https://github.com/videojs/mpd-parser/commit/71b8976))
* Support for inheriting BaseURL and alternate BaseURLs ([#17](https://github.com/videojs/mpd-parser/issues/17)) ([7dad5d5](https://github.com/videojs/mpd-parser/commit/7dad5d5))
* add support for SegmentTemplate padding format string and SegmentTimeline ([#16](https://github.com/videojs/mpd-parser/issues/16)) ([87933f6](https://github.com/videojs/mpd-parser/commit/87933f6))

<a name="0.2.1"></a>
## [0.2.1](https://github.com/videojs/mpd-parser/compare/v0.2.0...v0.2.1) (2017-12-15)

### Bug Fixes

* access HTMLCollections with IE11 compatibility ([#15](https://github.com/videojs/mpd-parser/issues/15)) ([5612984](https://github.com/videojs/mpd-parser/commit/5612984))

<a name="0.2.0"></a>
# [0.2.0](https://github.com/videojs/mpd-parser/compare/v0.1.1...v0.2.0) (2017-12-12)

### Features

* Support for vtt ([#13](https://github.com/videojs/mpd-parser/issues/13)) ([96fc406](https://github.com/videojs/mpd-parser/commit/96fc406))

### Tests

* add more tests for vtt ([#14](https://github.com/videojs/mpd-parser/issues/14)) ([4068790](https://github.com/videojs/mpd-parser/commit/4068790))

<a name="0.1.1"></a>
## [0.1.1](https://github.com/videojs/mpd-parser/compare/v0.1.0...v0.1.1) (2017-12-07)

### Bug Fixes

* avoid using Array.prototype.fill for IE support ([#11](https://github.com/videojs/mpd-parser/issues/11)) ([5c444de](https://github.com/videojs/mpd-parser/commit/5c444de))

<a name="0.1.0"></a>
# 0.1.0 (2017-11-29)

### Bug Fixes

* switch off in-manifest caption support ([#8](https://github.com/videojs/mpd-parser/issues/8)) ([15712c6](https://github.com/videojs/mpd-parser/commit/15712c6))

CHANGELOG
=========

## HEAD (Unreleased)
_(none)_

--------------------


<a name="2.2.0"></a>
# [2.2.0](https://github.com/videojs/http-streaming/compare/v2.1.0...v2.2.0) (2020-09-25)

### Features

* default handleManfiestRedirect to true ([#927](https://github.com/videojs/http-streaming/issues/927)) ([556321f](https://github.com/videojs/http-streaming/commit/556321f))
* support MPD.Location ([#926](https://github.com/videojs/http-streaming/issues/926)) ([c4a43d7](https://github.com/videojs/http-streaming/commit/c4a43d7))
* Update minimumUpdatePeriod handling ([#942](https://github.com/videojs/http-streaming/issues/942)) ([8648e76](https://github.com/videojs/http-streaming/commit/8648e76))

### Bug Fixes

* audio groups with the same uri as media do not count ([#952](https://github.com/videojs/http-streaming/issues/952)) ([3927c0c](https://github.com/videojs/http-streaming/commit/3927c0c))
* dash manifest not refreshed if only some playlists are updated ([#949](https://github.com/videojs/http-streaming/issues/949)) ([31d3441](https://github.com/videojs/http-streaming/commit/31d3441))
* detect demuxed video underflow gaps ([#948](https://github.com/videojs/http-streaming/issues/948)) ([d0ef298](https://github.com/videojs/http-streaming/commit/d0ef298))
* MPD not refreshed if minimumUpdatePeriod is 0 ([#954](https://github.com/videojs/http-streaming/issues/954)) ([3a0682f](https://github.com/videojs/http-streaming/commit/3a0682f)), closes [#942](https://github.com/videojs/http-streaming/issues/942)
* noop vtt segment loader handle data ([#959](https://github.com/videojs/http-streaming/issues/959)) ([d1dcd7b](https://github.com/videojs/http-streaming/commit/d1dcd7b))
* report the correct buffered regardless of playlist change ([#950](https://github.com/videojs/http-streaming/issues/950)) ([043ccc6](https://github.com/videojs/http-streaming/commit/043ccc6))
* Throw a player error when trying to play DRM content without eme ([#938](https://github.com/videojs/http-streaming/issues/938)) ([ce4d6fd](https://github.com/videojs/http-streaming/commit/ce4d6fd))
* use playlist NAME when available as its ID ([#929](https://github.com/videojs/http-streaming/issues/929)) ([2269464](https://github.com/videojs/http-streaming/commit/2269464))
* use TIME_FUDGE_FACTOR rather than rounding by decimal digits ([#881](https://github.com/videojs/http-streaming/issues/881)) ([7eb112d](https://github.com/videojs/http-streaming/commit/7eb112d))

### Chores

* **package:** remove engine check in pkcs7 ([#947](https://github.com/videojs/http-streaming/issues/947)) ([89392fa](https://github.com/videojs/http-streaming/commit/89392fa))
* mark angel one dash subs as broken ([#956](https://github.com/videojs/http-streaming/issues/956)) ([56a0970](https://github.com/videojs/http-streaming/commit/56a0970))
* mediaConfig_ -> staringMediaInfo_, startingMedia_ -> currentMediaInfo_ ([#953](https://github.com/videojs/http-streaming/issues/953)) ([8801d1c](https://github.com/videojs/http-streaming/commit/8801d1c))
* playlist selector logging ([#921](https://github.com/videojs/http-streaming/issues/921)) ([ccdbaef](https://github.com/videojs/http-streaming/commit/ccdbaef))
* update m3u8-parser to v4.4.3 ([#928](https://github.com/videojs/http-streaming/issues/928)) ([af5b4ee](https://github.com/videojs/http-streaming/commit/af5b4ee))

### Reverts

* fix: use playlist NAME when available as its ID ([#929](https://github.com/videojs/http-streaming/issues/929)) ([#957](https://github.com/videojs/http-streaming/issues/957)) ([fe8376b](https://github.com/videojs/http-streaming/commit/fe8376b))

<a name="2.1.0"></a>
# [2.1.0](https://github.com/videojs/http-streaming/compare/v2.0.0...v2.1.0) (2020-07-28)

### Features

* Easier manual playlist switching, add codecs to renditions ([#850](https://github.com/videojs/http-streaming/issues/850)) ([f60fa1f](https://github.com/videojs/http-streaming/commit/f60fa1f))
* exclude all incompatable browser/muxer codecs ([#903](https://github.com/videojs/http-streaming/issues/903)) ([2d0f0d7](https://github.com/videojs/http-streaming/commit/2d0f0d7))
* expose canChangeType on the VHS property ([#911](https://github.com/videojs/http-streaming/issues/911)) ([a4ab285](https://github.com/videojs/http-streaming/commit/a4ab285))
* let back buffer be configurable ([8c96e6c](https://github.com/videojs/http-streaming/commit/8c96e6c))
* Support codecs switching when possible via sourceBuffer.changeType  ([#841](https://github.com/videojs/http-streaming/issues/841)) ([267cc34](https://github.com/videojs/http-streaming/commit/267cc34))

### Bug Fixes

* always append init segment after trackinfo change ([#913](https://github.com/videojs/http-streaming/issues/913)) ([ea3650a](https://github.com/videojs/http-streaming/commit/ea3650a))
* cleanup mediasource listeners on dispose ([#871](https://github.com/videojs/http-streaming/issues/871)) ([e50f4c9](https://github.com/videojs/http-streaming/commit/e50f4c9))
* do not try to use unsupported audio ([#896](https://github.com/videojs/http-streaming/issues/896)) ([7711b26](https://github.com/videojs/http-streaming/commit/7711b26))
* do not use remove source buffer on ie 11 ([#904](https://github.com/videojs/http-streaming/issues/904)) ([1ab0f07](https://github.com/videojs/http-streaming/commit/1ab0f07))
* do not wait for audio appends for muxed segments ([#894](https://github.com/videojs/http-streaming/issues/894)) ([406cbcd](https://github.com/videojs/http-streaming/commit/406cbcd))
* Fixed issue with MPEG-Dash MPD Playlist Finalisation during Live Play. ([#874](https://github.com/videojs/http-streaming/issues/874)) ([c807930](https://github.com/videojs/http-streaming/commit/c807930))
* handle null return value from CaptionParser.parse ([#890](https://github.com/videojs/http-streaming/issues/890)) ([7b8fff2](https://github.com/videojs/http-streaming/commit/7b8fff2)), closes [#863](https://github.com/videojs/http-streaming/issues/863)
* have reloadSourceOnError get src from player ([#893](https://github.com/videojs/http-streaming/issues/893)) ([1e50bc5](https://github.com/videojs/http-streaming/commit/1e50bc5)), closes [videojs/video.js#6744](https://github.com/videojs/video.js/issues/6744)
* initialize EME for all playlists and PSSH values ([#872](https://github.com/videojs/http-streaming/issues/872)) ([e0e497f](https://github.com/videojs/http-streaming/commit/e0e497f))
* more conservative stalled download check, better logging ([#884](https://github.com/videojs/http-streaming/issues/884)) ([615e77f](https://github.com/videojs/http-streaming/commit/615e77f))
* pause/abort loaders before an exclude, preventing bad appends ([#902](https://github.com/videojs/http-streaming/issues/902)) ([c9126e1](https://github.com/videojs/http-streaming/commit/c9126e1))
* stop alt loaders on main mediachanging to prevent append race ([#895](https://github.com/videojs/http-streaming/issues/895)) ([8690c78](https://github.com/videojs/http-streaming/commit/8690c78))
* Support aac data with or without id3 tags by using mux.js[@5](https://github.com/5).6.6 ([#899](https://github.com/videojs/http-streaming/issues/899)) ([9c742ce](https://github.com/videojs/http-streaming/commit/9c742ce))
* Use revokeObjectURL dispose for created MSE blob urls ([#849](https://github.com/videojs/http-streaming/issues/849)) ([ca73cac](https://github.com/videojs/http-streaming/commit/ca73cac))
* Wait for sourceBuffer creation so drm setup uses valid codecs ([#878](https://github.com/videojs/http-streaming/issues/878)) ([f879563](https://github.com/videojs/http-streaming/commit/f879563))

### Chores

* Add vhs & mpc (vhs.masterPlaylistController_) to window of index.html ([#875](https://github.com/videojs/http-streaming/issues/875)) ([bab61d6](https://github.com/videojs/http-streaming/commit/bab61d6))
* **demo:** add a representations selector to the demo page ([#901](https://github.com/videojs/http-streaming/issues/901)) ([0a54ae2](https://github.com/videojs/http-streaming/commit/0a54ae2))
* fix tears of steal playready on the demo page ([#915](https://github.com/videojs/http-streaming/issues/915)) ([29a10d0](https://github.com/videojs/http-streaming/commit/29a10d0))
* keep window vhs/mpc up to date on source switch ([#883](https://github.com/videojs/http-streaming/issues/883)) ([3ba85fd](https://github.com/videojs/http-streaming/commit/3ba85fd))
* update DASH stream urls ([#918](https://github.com/videojs/http-streaming/issues/918)) ([902c2a5](https://github.com/videojs/http-streaming/commit/902c2a5))
* update local video.js ([#876](https://github.com/videojs/http-streaming/issues/876)) ([c2cc9aa](https://github.com/videojs/http-streaming/commit/c2cc9aa))
* use playready license server ([#916](https://github.com/videojs/http-streaming/issues/916)) ([6728837](https://github.com/videojs/http-streaming/commit/6728837))

### Code Refactoring

* remove duplicate bufferIntersection code in util/buffer.js ([#880](https://github.com/videojs/http-streaming/issues/880)) ([0ca43bd](https://github.com/videojs/http-streaming/commit/0ca43bd))
* simplify setupEmeOptions and add tests ([#869](https://github.com/videojs/http-streaming/issues/869)) ([e3921ed](https://github.com/videojs/http-streaming/commit/e3921ed))

<a name="2.0.0"></a>
# [2.0.0](https://github.com/videojs/http-streaming/compare/v2.0.0-rc.2...v2.0.0) (2020-06-16)

### Features

* add external vhs properties and deprecate hls and dash references ([#859](https://github.com/videojs/http-streaming/issues/859)) ([22af0b2](https://github.com/videojs/http-streaming/commit/22af0b2))
* Use VHS playback on any non-Safari browser ([#843](https://github.com/videojs/http-streaming/issues/843)) ([225d127](https://github.com/videojs/http-streaming/commit/225d127))

### Chores

* fix demo page on firefox, always use vhs on safari ([#851](https://github.com/videojs/http-streaming/issues/851)) ([d567b7d](https://github.com/videojs/http-streaming/commit/d567b7d))
* **stats:** update vhs usage in the stats page ([#867](https://github.com/videojs/http-streaming/issues/867)) ([4dda42a](https://github.com/videojs/http-streaming/commit/4dda42a))

### Code Refactoring

* Move caption parser to webworker, saving 5732b offloading work ([#863](https://github.com/videojs/http-streaming/issues/863)) ([491d194](https://github.com/videojs/http-streaming/commit/491d194))
* remove aes-decrypter objects from Hls saving 1415gz bytes ([#860](https://github.com/videojs/http-streaming/issues/860)) ([a4f8302](https://github.com/videojs/http-streaming/commit/a4f8302))

### Documentation

* add supported features doc ([#848](https://github.com/videojs/http-streaming/issues/848)) ([38f5860](https://github.com/videojs/http-streaming/commit/38f5860))

### Reverts

* "fix: Use middleware and a wrapped function for seeking instead of relying on unreliable 'seeking' events ([#161](https://github.com/videojs/http-streaming/issues/161))"([#856](https://github.com/videojs/http-streaming/issues/856)) ([1165f8e](https://github.com/videojs/http-streaming/commit/1165f8e))


### BREAKING CHANGES

* The Hls object which was exposed on videojs no longer has Decrypter, AsyncStream, and decrypt from aes-decrypter.

<a name="1.10.2"></a>
## [1.10.2](https://github.com/videojs/http-streaming/compare/v1.10.1...v1.10.2) (2019-05-13)

### Bug Fixes

* clear the blacklist for other playlists if final rendition errors ([#479](https://github.com/videojs/http-streaming/issues/479)) ([fe3b378](https://github.com/videojs/http-streaming/commit/fe3b378)), closes [#396](https://github.com/videojs/http-streaming/issues/396) [#471](https://github.com/videojs/http-streaming/issues/471)
* **development:** rollup watch, via `npm run watch`, should work for es/cjs ([#484](https://github.com/videojs/http-streaming/issues/484)) ([ad6f292](https://github.com/videojs/http-streaming/commit/ad6f292))
* **HLSe:** slice keys properly on IE11 ([#506](https://github.com/videojs/http-streaming/issues/506)) ([681cd6f](https://github.com/videojs/http-streaming/commit/681cd6f))
* **package:** update mpd-parser to version 0.8.1 🚀 ([#490](https://github.com/videojs/http-streaming/issues/490)) ([a49ad3a](https://github.com/videojs/http-streaming/commit/a49ad3a))
* **package:** update mux.js to version 5.1.2 🚀 ([#477](https://github.com/videojs/http-streaming/issues/477)) ([57a38e9](https://github.com/videojs/http-streaming/commit/57a38e9)), closes [#503](https://github.com/videojs/http-streaming/issues/503) [#504](https://github.com/videojs/http-streaming/issues/504)
* **source-updater:** run callbacks after setting timestampOffset ([#480](https://github.com/videojs/http-streaming/issues/480)) ([6ecf859](https://github.com/videojs/http-streaming/commit/6ecf859))
* livestream timeout issues ([#469](https://github.com/videojs/http-streaming/issues/469)) ([cf3fafc](https://github.com/videojs/http-streaming/commit/cf3fafc)), closes [segment#16](https://github.com/segment/issues/16) [segment#15](https://github.com/segment/issues/15) [segment#16](https://github.com/segment/issues/16) [segment#15](https://github.com/segment/issues/15) [segment#16](https://github.com/segment/issues/16)
* remove both vttjs listeners to prevent leaking one of them ([#495](https://github.com/videojs/http-streaming/issues/495)) ([1db1e72](https://github.com/videojs/http-streaming/commit/1db1e72))

### Performance Improvements

* don't enable captionParser for audio or subtitle loaders ([#487](https://github.com/videojs/http-streaming/issues/487)) ([358877f](https://github.com/videojs/http-streaming/commit/358877f))

<a name="1.10.1"></a>
## [1.10.1](https://github.com/videojs/http-streaming/compare/v1.10.0...v1.10.1) (2019-04-16)

### Bug Fixes

* **dash-playlist-loader:** clear out timers on dispose ([#472](https://github.com/videojs/http-streaming/issues/472)) ([2f1c222](https://github.com/videojs/http-streaming/commit/2f1c222))

### Reverts

* "fix: clear the blacklist for other playlists if final rendition errors ([#396](https://github.com/videojs/http-streaming/issues/396))" ([#471](https://github.com/videojs/http-streaming/issues/471)) ([dd55028](https://github.com/videojs/http-streaming/commit/dd55028))

<a name="1.10.0"></a>
# [1.10.0](https://github.com/videojs/http-streaming/compare/v1.9.3...v1.10.0) (2019-04-12)

### Features

* add option to cache encrpytion keys in the player ([#446](https://github.com/videojs/http-streaming/issues/446)) ([599b94d](https://github.com/videojs/http-streaming/commit/599b94d)), closes [#140](https://github.com/videojs/http-streaming/issues/140)
* add support for dash manifests describing sidx boxes ([#455](https://github.com/videojs/http-streaming/issues/455)) ([80dde16](https://github.com/videojs/http-streaming/commit/80dde16))

### Bug Fixes

* clear the blacklist for other playlists if final rendition errors ([#396](https://github.com/videojs/http-streaming/issues/396)) ([6e6c8c2](https://github.com/videojs/http-streaming/commit/6e6c8c2))
* on dispose, don't call abort on SourceBuffer until after remove() has finished ([3806750](https://github.com/videojs/http-streaming/commit/3806750))

### Documentation

* **README:** update broken link to full docs ([#440](https://github.com/videojs/http-streaming/issues/440)) ([fbd615c](https://github.com/videojs/http-streaming/commit/fbd615c))

<a name="1.9.3"></a>
## [1.9.3](https://github.com/videojs/http-streaming/compare/v1.9.2...v1.9.3) (2019-03-21)

### Bug Fixes

* **id3:** ignore unsupported id3 frames ([#437](https://github.com/videojs/http-streaming/issues/437)) ([7040b7d](https://github.com/videojs/http-streaming/commit/7040b7d)), closes [videojs/video.js#5823](https://github.com/videojs/video.js/issues/5823)

### Documentation

* add diagrams for playlist loaders ([#426](https://github.com/videojs/http-streaming/issues/426)) ([52201f9](https://github.com/videojs/http-streaming/commit/52201f9))

<a name="1.9.2"></a>
## [1.9.2](https://github.com/videojs/http-streaming/compare/v1.9.1...v1.9.2) (2019-03-14)

### Bug Fixes

* expose `custom` segment property in the segment metadata track ([#429](https://github.com/videojs/http-streaming/issues/429)) ([17510da](https://github.com/videojs/http-streaming/commit/17510da))

<a name="1.9.1"></a>
## [1.9.1](https://github.com/videojs/http-streaming/compare/v1.9.0...v1.9.1) (2019-03-05)

### Bug Fixes

* fix for streams that would occasionally never fire an `ended` event ([fc09926](https://github.com/videojs/http-streaming/commit/fc09926))
* Fix video playback freezes caused by not using absolute current time ([#401](https://github.com/videojs/http-streaming/issues/401)) ([957ecfd](https://github.com/videojs/http-streaming/commit/957ecfd))
* only fire seekablechange when values of seekable ranges actually change ([#415](https://github.com/videojs/http-streaming/issues/415)) ([a4c056e](https://github.com/videojs/http-streaming/commit/a4c056e))
* Prevent infinite buffering at the start of looped video on edge ([#392](https://github.com/videojs/http-streaming/issues/392)) ([b6d1b97](https://github.com/videojs/http-streaming/commit/b6d1b97))

### Code Refactoring

* align DashPlaylistLoader closer to PlaylistLoader states ([#386](https://github.com/videojs/http-streaming/issues/386)) ([5d80fe7](https://github.com/videojs/http-streaming/commit/5d80fe7))

<a name="1.9.0"></a>
# [1.9.0](https://github.com/videojs/http-streaming/compare/v1.8.0...v1.9.0) (2019-02-07)

### Features

* Use exposed transmuxer time modifications for more accurate conversion between program and player times ([#371](https://github.com/videojs/http-streaming/issues/371)) ([41df5c0](https://github.com/videojs/http-streaming/commit/41df5c0))

### Bug Fixes

* m3u8 playlist is not updating when only endList changes ([#373](https://github.com/videojs/http-streaming/issues/373)) ([c7d1306](https://github.com/videojs/http-streaming/commit/c7d1306))
* Prevent exceptions from being thrown by the MediaSource ([#389](https://github.com/videojs/http-streaming/issues/389)) ([8c06366](https://github.com/videojs/http-streaming/commit/8c06366))

### Chores

* Update mux.js to the latest version 🚀 ([#397](https://github.com/videojs/http-streaming/issues/397)) ([38ec2a5](https://github.com/videojs/http-streaming/commit/38ec2a5))

### Tests

* added test for playlist not updating when only endList changes ([#394](https://github.com/videojs/http-streaming/issues/394)) ([39d0be2](https://github.com/videojs/http-streaming/commit/39d0be2))

<a name="1.8.0"></a>
# [1.8.0](https://github.com/videojs/http-streaming/compare/v1.7.0...v1.8.0) (2019-01-10)

### Features

* expose custom M3U8 mapper API ([#325](https://github.com/videojs/http-streaming/issues/325)) ([609beb3](https://github.com/videojs/http-streaming/commit/609beb3))

### Bug Fixes

* **id3:** cuechange event not being triggered on audio-only HLS streams ([#334](https://github.com/videojs/http-streaming/issues/334)) ([bab70fd](https://github.com/videojs/http-streaming/commit/bab70fd)), closes [#130](https://github.com/videojs/http-streaming/issues/130)

<a name="1.7.0"></a>
# [1.7.0](https://github.com/videojs/http-streaming/compare/v1.6.0...v1.7.0) (2019-01-04)

### Features

* expose custom M3U8 parser API ([#331](https://github.com/videojs/http-streaming/issues/331)) ([b0643a4](https://github.com/videojs/http-streaming/commit/b0643a4))

<a name="1.6.0"></a>
# [1.6.0](https://github.com/videojs/http-streaming/compare/v1.5.1...v1.6.0) (2018-12-21)

### Features

* Add allowSeeksWithinUnsafeLiveWindow property ([#320](https://github.com/videojs/http-streaming/issues/320)) ([74b28e8](https://github.com/videojs/http-streaming/commit/74b28e8))

### Chores

* add clock.ticks to now async operations in tests ([#315](https://github.com/videojs/http-streaming/issues/315)) ([895c86a](https://github.com/videojs/http-streaming/commit/895c86a))

### Documentation

* Add README entry on DRM and videojs-contrib-eme ([#307](https://github.com/videojs/http-streaming/issues/307)) ([93b6167](https://github.com/videojs/http-streaming/commit/93b6167))

<a name="1.5.1"></a>
## [1.5.1](https://github.com/videojs/http-streaming/compare/v1.5.0...v1.5.1) (2018-12-06)

### Bug Fixes

* added missing manifest information on to segments (EXT-X-PROGRAM-DATE-TIME) ([#236](https://github.com/videojs/http-streaming/issues/236)) ([a35dd09](https://github.com/videojs/http-streaming/commit/a35dd09))
* remove player props on dispose to stop middleware  ([#229](https://github.com/videojs/http-streaming/issues/229)) ([cd13f9f](https://github.com/videojs/http-streaming/commit/cd13f9f))

### Documentation

* add dash to package.json description ([#267](https://github.com/videojs/http-streaming/issues/267)) ([3296c68](https://github.com/videojs/http-streaming/commit/3296c68))
* add documentation for reloadSourceOnError ([#266](https://github.com/videojs/http-streaming/issues/266)) ([7448b37](https://github.com/videojs/http-streaming/commit/7448b37))

<a name="1.5.0"></a>
# [1.5.0](https://github.com/videojs/http-streaming/compare/v1.4.2...v1.5.0) (2018-11-13)

### Features

* Add useBandwidthFromLocalStorage option ([#275](https://github.com/videojs/http-streaming/issues/275)) ([60c88ae](https://github.com/videojs/http-streaming/commit/60c88ae))

### Bug Fixes

* don't wait for requests to finish when encountering an error in media-segment-request ([#286](https://github.com/videojs/http-streaming/issues/286)) ([970e3ce](https://github.com/videojs/http-streaming/commit/970e3ce))
* throttle final playlist reloads when using DASH ([#277](https://github.com/videojs/http-streaming/issues/277)) ([1c2887a](https://github.com/videojs/http-streaming/commit/1c2887a))

<a name="1.4.2"></a>
## [1.4.2](https://github.com/videojs/http-streaming/compare/v1.4.1...v1.4.2) (2018-11-01)

### Chores

* pin to node 8 for now ([#279](https://github.com/videojs/http-streaming/issues/279)) ([f900dc4](https://github.com/videojs/http-streaming/commit/f900dc4))
* update mux.js to 5.0.1 ([#282](https://github.com/videojs/http-streaming/issues/282)) ([af6ee4f](https://github.com/videojs/http-streaming/commit/af6ee4f))

<a name="1.4.1"></a>
## [1.4.1](https://github.com/videojs/http-streaming/compare/v1.4.0...v1.4.1) (2018-10-25)

### Bug Fixes

* **subtitles:** set default property if default and autoselect are both enabled ([#239](https://github.com/videojs/http-streaming/issues/239)) ([ee594e5](https://github.com/videojs/http-streaming/commit/ee594e5))

<a name="1.4.0"></a>
# [1.4.0](https://github.com/videojs/http-streaming/compare/v1.3.1...v1.4.0) (2018-10-24)

### Features

* limited experimental DASH multiperiod support ([#268](https://github.com/videojs/http-streaming/issues/268)) ([a213807](https://github.com/videojs/http-streaming/commit/a213807))
* smoothQualityChange flag ([#235](https://github.com/videojs/http-streaming/issues/235)) ([0e4fdf9](https://github.com/videojs/http-streaming/commit/0e4fdf9))

### Bug Fixes

* immediately setup EME if available ([#263](https://github.com/videojs/http-streaming/issues/263)) ([7577e90](https://github.com/videojs/http-streaming/commit/7577e90))

<a name="1.3.1"></a>
## [1.3.1](https://github.com/videojs/http-streaming/compare/v1.3.0...v1.3.1) (2018-10-15)

### Bug Fixes

* ensure content loops ([#259](https://github.com/videojs/http-streaming/issues/259)) ([26300df](https://github.com/videojs/http-streaming/commit/26300df))

<a name="1.3.0"></a>
# [1.3.0](https://github.com/videojs/http-streaming/compare/v1.2.6...v1.3.0) (2018-10-05)

### Features

* add an option to ignore player size in selection logic ([#238](https://github.com/videojs/http-streaming/issues/238)) ([7ae42b1](https://github.com/videojs/http-streaming/commit/7ae42b1))

### Documentation

* Update CONTRIBUTING.md ([#242](https://github.com/videojs/http-streaming/issues/242)) ([9d83e9d](https://github.com/videojs/http-streaming/commit/9d83e9d))

<a name="1.2.6"></a>
## [1.2.6](https://github.com/videojs/http-streaming/compare/v1.2.5...v1.2.6) (2018-09-21)

### Bug Fixes

* stutter after fast quality change in IE/Edge ([#213](https://github.com/videojs/http-streaming/issues/213)) ([2c0d9b2](https://github.com/videojs/http-streaming/commit/2c0d9b2))

### Documentation

* update issue template to link to the troubleshooting guide ([#215](https://github.com/videojs/http-streaming/issues/215)) ([413f0e8](https://github.com/videojs/http-streaming/commit/413f0e8))
* update README notes for video.js 7 ([#200](https://github.com/videojs/http-streaming/issues/200)) ([d68ce0c](https://github.com/videojs/http-streaming/commit/d68ce0c))
* update troubleshooting guide for Edge/mobile Chrome ([#216](https://github.com/videojs/http-streaming/issues/216)) ([21e5335](https://github.com/videojs/http-streaming/commit/21e5335))

<a name="1.2.5"></a>
## [1.2.5](https://github.com/videojs/http-streaming/compare/v1.2.4...v1.2.5) (2018-08-24)

### Bug Fixes

* fix replay functionality ([#204](https://github.com/videojs/http-streaming/issues/204)) ([fd6be83](https://github.com/videojs/http-streaming/commit/fd6be83))

<a name="1.2.4"></a>
## [1.2.4](https://github.com/videojs/http-streaming/compare/v1.2.3...v1.2.4) (2018-08-13)

### Bug Fixes

* Remove buffered data on fast quality switches ([#113](https://github.com/videojs/http-streaming/issues/113)) ([bc94fbb](https://github.com/videojs/http-streaming/commit/bc94fbb))

<a name="1.2.3"></a>
## [1.2.3](https://github.com/videojs/http-streaming/compare/v1.2.2...v1.2.3) (2018-08-09)

### Chores

* link to minified example in main page ([#189](https://github.com/videojs/http-streaming/issues/189)) ([15a7f92](https://github.com/videojs/http-streaming/commit/15a7f92))
* use netlify for easier testing ([#188](https://github.com/videojs/http-streaming/issues/188)) ([d2e0d35](https://github.com/videojs/http-streaming/commit/d2e0d35))

<a name="1.2.2"></a>
## [1.2.2](https://github.com/videojs/http-streaming/compare/v1.2.1...v1.2.2) (2018-08-07)

### Bug Fixes

* typeof minification ([#182](https://github.com/videojs/http-streaming/issues/182)) ([7c68335](https://github.com/videojs/http-streaming/commit/7c68335))
* Use middleware and a wrapped function for seeking instead of relying on unreliable 'seeking' events ([#161](https://github.com/videojs/http-streaming/issues/161)) ([6c68761](https://github.com/videojs/http-streaming/commit/6c68761))

### Chores

* add logo ([#184](https://github.com/videojs/http-streaming/issues/184)) ([a55626c](https://github.com/videojs/http-streaming/commit/a55626c))

### Documentation

* add note for Safari captions error ([#174](https://github.com/videojs/http-streaming/issues/174)) ([7b03530](https://github.com/videojs/http-streaming/commit/7b03530))

### Tests

* add support for real segments in tests ([#178](https://github.com/videojs/http-streaming/issues/178)) ([2b07fca](https://github.com/videojs/http-streaming/commit/2b07fca))

<a name="1.2.1"></a>
## [1.2.1](https://github.com/videojs/http-streaming/compare/v1.2.0...v1.2.1) (2018-07-17)

### Bug Fixes

* convert non-latin characters in IE ([#157](https://github.com/videojs/http-streaming/issues/157)) ([17678fb](https://github.com/videojs/http-streaming/commit/17678fb))

<a name="1.2.0"></a>
# [1.2.0](https://github.com/videojs/http-streaming/compare/v1.1.0...v1.2.0) (2018-07-16)

### Features

* **captions:** write in-band captions from DASH fmp4 segments to the textTrack API ([#108](https://github.com/videojs/http-streaming/issues/108)) ([7c11911](https://github.com/videojs/http-streaming/commit/7c11911))

### Chores

* add welcome bot config from video.js ([#150](https://github.com/videojs/http-streaming/issues/150)) ([922cfee](https://github.com/videojs/http-streaming/commit/922cfee))

<a name="1.1.0"></a>
# [1.1.0](https://github.com/videojs/http-streaming/compare/v1.0.2...v1.1.0) (2018-06-06)

### Features

* Utilize option to override native on tech ([#76](https://github.com/videojs/http-streaming/issues/76)) ([5c7ab4c](https://github.com/videojs/http-streaming/commit/5c7ab4c))

### Chores

* update tests and pages for video.js 7 ([#102](https://github.com/videojs/http-streaming/issues/102)) ([d6f5005](https://github.com/videojs/http-streaming/commit/d6f5005))

<a name="1.0.2"></a>
## [1.0.2](https://github.com/videojs/http-streaming/compare/v1.0.1...v1.0.2) (2018-05-17)

### Bug Fixes

* make project Video.js 7 ready ([#92](https://github.com/videojs/http-streaming/issues/92)) ([decad87](https://github.com/videojs/http-streaming/commit/decad87))
* make sure that es build is babelified ([#97](https://github.com/videojs/http-streaming/issues/97)) ([5f0428d](https://github.com/videojs/http-streaming/commit/5f0428d))

### Documentation

* update documentation with a glossary and intro page, added DASH background ([#94](https://github.com/videojs/http-streaming/issues/94)) ([4b0fde9](https://github.com/videojs/http-streaming/commit/4b0fde9))

<a name="1.0.1"></a>
## [1.0.1](https://github.com/videojs/http-streaming/compare/v1.0.0...v1.0.1) (2018-04-12)

### Bug Fixes

* minified build ([#84](https://github.com/videojs/http-streaming/issues/84)) ([2402ac6](https://github.com/videojs/http-streaming/commit/2402ac6))

<a name="1.0.0"></a>
# [1.0.0](https://github.com/videojs/http-streaming/compare/v0.9.0...v1.0.0) (2018-04-10)

### Chores

* sync videojs-contrib-hls updates ([#75](https://github.com/videojs/http-streaming/issues/75)) ([9223588](https://github.com/videojs/http-streaming/commit/9223588))
* update the aes-decrypter ([#71](https://github.com/videojs/http-streaming/issues/71)) ([27ed914](https://github.com/videojs/http-streaming/commit/27ed914))

### Documentation

* update docs for overrideNative ([#77](https://github.com/videojs/http-streaming/issues/77)) ([98ca6d3](https://github.com/videojs/http-streaming/commit/98ca6d3))
* update known issues for fmp4 captions ([#79](https://github.com/videojs/http-streaming/issues/79)) ([c418301](https://github.com/videojs/http-streaming/commit/c418301))

<a name="0.9.0"></a>
# [0.9.0](https://github.com/videojs/http-streaming/compare/v0.8.0...v0.9.0) (2018-03-30)

### Features

* support in-manifest DRM data ([#60](https://github.com/videojs/http-streaming/issues/60)) ([a1cad82](https://github.com/videojs/http-streaming/commit/a1cad82))

<a name="0.8.0"></a>
# [0.8.0](https://github.com/videojs/http-streaming/compare/v0.7.2...v0.8.0) (2018-03-30)

### Code Refactoring

* export corrections ([#68](https://github.com/videojs/http-streaming/issues/68)) ([aab3b90](https://github.com/videojs/http-streaming/commit/aab3b90))
* use rollup for build ([#69](https://github.com/videojs/http-streaming/issues/69)) ([c28c25c](https://github.com/videojs/http-streaming/commit/c28c25c))

# 0.7.0
* feat: Live support for DASH

# 0.6.1
* use webwackify for webworkers to support webpack bundle ([#50](https://github.com/videojs/http-streaming/pull/45))

# 0.5.3
* fix: program date time handling ([#45](https://github.com/videojs/http-streaming/pull/45))
  * update m3u8-parser to v4.2.0
  * use segment program date time info
* feat: Adding support for segments in Period and Representation ([#47](https://github.com/videojs/http-streaming/pull/47))
* wait for both main and audio loaders for endOfStream if main starting media unknown ([#44](https://github.com/videojs/http-streaming/pull/44))

# 0.5.2
* add debug logging statement for seekable updates ([#40](https://github.com/videojs/http-streaming/pull/40))

# 0.5.1
* Fix audio only streams with EXT-X-MEDIA tags ([#34](https://github.com/videojs/http-streaming/pull/34))
* Merge videojs-contrib-hls master into http-streaming master ([#35](https://github.com/videojs/http-streaming/pull/35))
  * Update sinon to 1.10.3=
  * Update videojs-contrib-quality-levels to ^2.0.4
  * Fix test for event handler cleanup on dispose by calling event handling methods
* fix: Don't reset eme options ([#32](https://github.com/videojs/http-streaming/pull/32))

# 0.5.0
* update mpd-parser to support more segment list types ([#27](https://github.com/videojs/http-streaming/issues/27))

# 0.4.0
* Removed Flash support ([#15](https://github.com/videojs/http-streaming/issues/15))
* Blacklist playlists not supported by browser media source before initial selection ([#17](https://github.com/videojs/http-streaming/issues/17))

# 0.3.1
* Skip flash-based source handler with DASH sources ([#14](https://github.com/videojs/http-streaming/issues/14))

# 0.3.0
* Added additional properties to the stats object ([#10](https://github.com/videojs/http-streaming/issues/10))

# 0.2.1
* Updated the mpd-parser to fix IE11 DASH support ([#12](https://github.com/videojs/http-streaming/issues/12))

# 0.2.0
* Initial DASH Support ([#8](https://github.com/videojs/http-streaming/issues/8))

# 0.1.0
* Initial release, based on [videojs-contrib-hls 5.12.2](https://github.com/videojs/videojs-contrib-hls)


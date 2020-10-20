<a name="1.13.2"></a>
## [1.13.2](https://github.com/videojs/http-streaming/compare/v1.13.1...v1.13.2) (2020-03-30)

### Bug Fixes

* dispose workers on dispose ([#788](https://github.com/videojs/http-streaming/issues/788)) ([03ddb4e](https://github.com/videojs/http-streaming/commit/03ddb4e))

<a name="1.13.1"></a>
## [1.13.1](https://github.com/videojs/http-streaming/compare/v1.13.0...v1.13.1) (2020-03-28)

### Bug Fixes

* put devicePixelRatio behind useDevicePixelRatio option ([#785](https://github.com/videojs/http-streaming/issues/785)) ([57335d9](https://github.com/videojs/http-streaming/commit/57335d9))

<a name="1.13.0"></a>
# [1.13.0](https://github.com/videojs/http-streaming/compare/v1.12.3...v1.13.0) (2020-03-26)

### Bug Fixes

* null check return value of selectPlaylist ([#779](https://github.com/videojs/http-streaming/issues/779)) ([90a0215](https://github.com/videojs/http-streaming/commit/90a0215))
* take devicePixelRatio into account during ABR  ([#784](https://github.com/videojs/http-streaming/issues/784)) ([bd63e57](https://github.com/videojs/http-streaming/commit/bd63e57)), closes [#744](https://github.com/videojs/http-streaming/issues/744)

### Chores

* port index page from master ([#774](https://github.com/videojs/http-streaming/issues/774)) ([808c3b1](https://github.com/videojs/http-streaming/commit/808c3b1))

### Reverts

* "fix: Use middleware and a wrapped function for seeking instead of relying on unreliable 'seeking' events ([#161](https://github.com/videojs/http-streaming/issues/161))" ([#777](https://github.com/videojs/http-streaming/issues/777)) ([1a4fc1e](https://github.com/videojs/http-streaming/commit/1a4fc1e)), closes [#378](https://github.com/videojs/http-streaming/issues/378) [videojs/video.js#6444](https://github.com/videojs/video.js/issues/6444)

<a name="1.12.3"></a>
## [1.12.3](https://github.com/videojs/http-streaming/compare/v1.12.2...v1.12.3) (2020-03-16)

### Bug Fixes

* **segment-loader:** resetEverything should remove through Infinity ([#754](https://github.com/videojs/http-streaming/issues/754)) ([#758](https://github.com/videojs/http-streaming/issues/758)) ([6ba1800](https://github.com/videojs/http-streaming/commit/6ba1800))
* add native cues when using native text tracks ([#769](https://github.com/videojs/http-streaming/issues/769)) ([10d25d1](https://github.com/videojs/http-streaming/commit/10d25d1)), closes [videojs/video.js#6410](https://github.com/videojs/video.js/issues/6410)

### Tests

* skip flaky test ([#771](https://github.com/videojs/http-streaming/issues/771)) ([99bf807](https://github.com/videojs/http-streaming/commit/99bf807))

<a name="1.12.2"></a>
## [1.12.2](https://github.com/videojs/http-streaming/compare/v1.12.1...v1.12.2) (2020-02-18)

### Bug Fixes

* add dispose functions and fix memory leaks [#643](https://github.com/videojs/http-streaming/issues/643) for 1.x ([#734](https://github.com/videojs/http-streaming/issues/734)) ([89ab859](https://github.com/videojs/http-streaming/commit/89ab859))
* resume live DASH playlist refreshes after pausing and loading DASH playlist loader ([#736](https://github.com/videojs/http-streaming/issues/736)) ([e966e9c](https://github.com/videojs/http-streaming/commit/e966e9c))
* trim 30s back from playhead even for VOD and LIVE DVR content ([#740](https://github.com/videojs/http-streaming/issues/740)) ([886f592](https://github.com/videojs/http-streaming/commit/886f592))

<a name="1.12.1"></a>
## [1.12.1](https://github.com/videojs/http-streaming/compare/v1.12.0...v1.12.1) (2020-02-11)

### Bug Fixes

* update to mux.js 5.5.1 ([#733](https://github.com/videojs/http-streaming/issues/733)) ([aa22f31](https://github.com/videojs/http-streaming/commit/aa22f31))

### Code Refactoring

* add semicolon and move variable closer to usage ([#729](https://github.com/videojs/http-streaming/issues/729)) ([675b1e9](https://github.com/videojs/http-streaming/commit/675b1e9))

<a name="1.12.0"></a>
# [1.12.0](https://github.com/videojs/http-streaming/compare/v1.11.3...v1.12.0) (2020-02-04)

### Features

* support suggestedPresentationDelay in DASH manifests ([#698](https://github.com/videojs/http-streaming/issues/698)) ([c14fb43](https://github.com/videojs/http-streaming/commit/c14fb43))

<a name="1.11.3"></a>
## [1.11.3](https://github.com/videojs/http-streaming/compare/v1.11.2...v1.11.3) (2020-01-17)

### Bug Fixes

* consider `hidden` tracks as active ([#564](https://github.com/videojs/http-streaming/issues/564)) ([6acdd20](https://github.com/videojs/http-streaming/commit/6acdd20))
* live startup failures when play happens before playlist is downloaded ([#700](https://github.com/videojs/http-streaming/issues/700)) ([92c93a7](https://github.com/videojs/http-streaming/commit/92c93a7)), closes [#464](https://github.com/videojs/http-streaming/issues/464) [#496](https://github.com/videojs/http-streaming/issues/496) [#500](https://github.com/videojs/http-streaming/issues/500)
* race condition preventing qualityLevels from being populating ([#707](https://github.com/videojs/http-streaming/issues/707)) ([8c4a11f](https://github.com/videojs/http-streaming/commit/8c4a11f)), closes [#677](https://github.com/videojs/http-streaming/issues/677)
* support multiple stream-inf with same URI ([#672](https://github.com/videojs/http-streaming/issues/672)) ([095515c](https://github.com/videojs/http-streaming/commit/095515c))

### Chores

* **stats:** add liveui option and various fixes ([#695](https://github.com/videojs/http-streaming/issues/695)) ([a0f6c8b](https://github.com/videojs/http-streaming/commit/a0f6c8b))

<a name="1.11.2"></a>
## [1.11.2](https://github.com/videojs/http-streaming/compare/v1.11.1...v1.11.2) (2019-11-12)

### Bug Fixes

* only do ondemand eme key initialization on non-ie11 browsers ([#682](https://github.com/videojs/http-streaming/issues/682)) ([fbfe68f](https://github.com/videojs/http-streaming/commit/fbfe68f))

<a name="1.11.1"></a>
## [1.11.1](https://github.com/videojs/http-streaming/compare/v1.11.0...v1.11.1) (2019-10-07)

### Bug Fixes

* improve timestampOffset calculation for fmp4s ([#666](https://github.com/videojs/http-streaming/issues/666)) ([bedc824](https://github.com/videojs/http-streaming/commit/bedc824))

<a name="1.11.0"></a>
# [1.11.0](https://github.com/videojs/http-streaming/compare/v1.11.0-0...v1.11.0) (2019-09-25)

<a name="1.10.6"></a>
## [1.10.6](https://github.com/videojs/http-streaming/compare/v1.10.5...v1.10.6) (2019-08-28)

### Bug Fixes

* fix anamorphic video with chrome ([#648](https://github.com/videojs/http-streaming/issues/648)) ([50b5992](https://github.com/videojs/http-streaming/commit/50b5992)), closes [#312](https://github.com/videojs/http-streaming/issues/312)

### Performance Improvements

* save 15430 gzipped bytes with better mux.js imports ([#628](https://github.com/videojs/http-streaming/issues/628)) ([4e69d21](https://github.com/videojs/http-streaming/commit/4e69d21))

<a name="1.10.5"></a>
## [1.10.5](https://github.com/videojs/http-streaming/compare/v1.10.5-1...v1.10.5) (2019-08-16)

### Chores

* **package:** update to mux.js 5.2.0 ([#636](https://github.com/videojs/http-streaming/issues/636)) ([0a793b0](https://github.com/videojs/http-streaming/commit/0a793b0))

### Tests

* **zero-length:** update segment and unskip test ([#635](https://github.com/videojs/http-streaming/issues/635)) ([665b526](https://github.com/videojs/http-streaming/commit/665b526))

<a name="1.10.4"></a>
## [1.10.4](https://github.com/videojs/http-streaming/compare/v1.10.3...v1.10.4) (2019-06-27)

### Bug Fixes

* verify sidx is present before doing anything related to it ([#561](https://github.com/videojs/http-streaming/issues/561)) ([15d753c](https://github.com/videojs/http-streaming/commit/15d753c)), closes [#457](https://github.com/videojs/http-streaming/issues/457) [#557](https://github.com/videojs/http-streaming/issues/557)

### Chores

* update package-lock ([1559930](https://github.com/videojs/http-streaming/commit/1559930))

<a name="1.10.3"></a>
## [1.10.3](https://github.com/videojs/http-streaming/compare/v1.10.2...v1.10.3) (2019-05-30)

### Bug Fixes

* only reset syncController if we passed a discontinuity ([#502](https://github.com/videojs/http-streaming/issues/502)) ([409634f](https://github.com/videojs/http-streaming/commit/409634f))

### Chores

* **package:** update mux.js to version 5.1.3 🚀 ([#510](https://github.com/videojs/http-streaming/issues/510)) ([55d764a](https://github.com/videojs/http-streaming/commit/55d764a))

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


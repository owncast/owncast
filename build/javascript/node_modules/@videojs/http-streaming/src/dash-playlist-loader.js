import videojs from 'video.js';
import {
  parse as parseMpd,
  parseUTCTiming
} from 'mpd-parser';
import {
  refreshDelay,
  updateMaster as updatePlaylist
} from './playlist-loader';
import { resolveUrl, resolveManifestRedirect } from './resolve-url';
import parseSidx from 'mux.js/lib/tools/parse-sidx';
import { segmentXhrHeaders } from './xhr';
import window from 'global/window';
import {
  forEachMediaGroup,
  addPropertiesToMaster
} from './manifest';
import containerRequest from './util/container-request.js';
import {toUint8} from '@videojs/vhs-utils/dist/byte-helpers';

const { EventTarget, mergeOptions } = videojs;

/**
 * Parses the master XML string and updates playlist URI references.
 *
 * @param {Object} config
 *        Object of arguments
 * @param {string} config.masterXml
 *        The mpd XML
 * @param {string} config.srcUrl
 *        The mpd URL
 * @param {Date} config.clientOffset
 *         A time difference between server and client
 * @param {Object} config.sidxMapping
 *        SIDX mappings for moof/mdat URIs and byte ranges
 * @return {Object}
 *         The parsed mpd manifest object
 */
export const parseMasterXml = ({ masterXml, srcUrl, clientOffset, sidxMapping }) => {
  const master = parseMpd(masterXml, {
    manifestUri: srcUrl,
    clientOffset,
    sidxMapping
  });

  addPropertiesToMaster(master, srcUrl);

  return master;
};

/**
 * Returns a new master manifest that is the result of merging an updated master manifest
 * into the original version.
 *
 * @param {Object} oldMaster
 *        The old parsed mpd object
 * @param {Object} newMaster
 *        The updated parsed mpd object
 * @return {Object}
 *         A new object representing the original master manifest with the updated media
 *         playlists merged in
 */
export const updateMaster = (oldMaster, newMaster) => {
  let noChanges = true;
  let update = mergeOptions(oldMaster, {
    // These are top level properties that can be updated
    duration: newMaster.duration,
    minimumUpdatePeriod: newMaster.minimumUpdatePeriod
  });

  // First update the playlists in playlist list
  for (let i = 0; i < newMaster.playlists.length; i++) {
    const playlistUpdate = updatePlaylist(update, newMaster.playlists[i]);

    if (playlistUpdate) {
      update = playlistUpdate;
      noChanges = false;
    }
  }

  // Then update media group playlists
  forEachMediaGroup(newMaster, (properties, type, group, label) => {
    if (properties.playlists && properties.playlists.length) {
      const id = properties.playlists[0].id;
      const playlistUpdate = updatePlaylist(update, properties.playlists[0]);

      if (playlistUpdate) {
        update = playlistUpdate;
        // update the playlist reference within media groups
        update.mediaGroups[type][group][label].playlists[0] = update.playlists[id];
        noChanges = false;
      }
    }
  });

  if (newMaster.minimumUpdatePeriod !== oldMaster.minimumUpdatePeriod) {
    noChanges = false;
  }

  if (noChanges) {
    return null;
  }

  return update;
};

export const generateSidxKey = (sidxInfo) => {
  // should be non-inclusive
  const sidxByteRangeEnd =
    sidxInfo.byterange.offset +
    sidxInfo.byterange.length -
    1;

  return sidxInfo.uri + '-' +
    sidxInfo.byterange.offset + '-' +
    sidxByteRangeEnd;
};

// SIDX should be equivalent if the URI and byteranges of the SIDX match.
// If the SIDXs have maps, the two maps should match,
// both `a` and `b` missing SIDXs is considered matching.
// If `a` or `b` but not both have a map, they aren't matching.
const equivalentSidx = (a, b) => {
  const neitherMap = Boolean(!a.map && !b.map);

  const equivalentMap = neitherMap || Boolean(a.map && b.map &&
    a.map.byterange.offset === b.map.byterange.offset &&
    a.map.byterange.length === b.map.byterange.length);

  return equivalentMap &&
    a.uri === b.uri &&
    a.byterange.offset === b.byterange.offset &&
    a.byterange.length === b.byterange.length;
};

// exported for testing
export const compareSidxEntry = (playlists, oldSidxMapping) => {
  const newSidxMapping = {};

  for (const id in playlists) {
    const playlist = playlists[id];
    const currentSidxInfo = playlist.sidx;

    if (currentSidxInfo) {
      const key = generateSidxKey(currentSidxInfo);

      if (!oldSidxMapping[key]) {
        break;
      }

      const savedSidxInfo = oldSidxMapping[key].sidxInfo;

      if (equivalentSidx(savedSidxInfo, currentSidxInfo)) {
        newSidxMapping[key] = oldSidxMapping[key];
      }
    }
  }

  return newSidxMapping;
};

/**
 *  A function that filters out changed items as they need to be requested separately.
 *
 *  The method is exported for testing
 *
 *  @param {Object} masterXml the mpd XML
 *  @param {string} srcUrl the mpd url
 *  @param {Date} clientOffset a time difference between server and client (passed through and not used)
 *  @param {Object} oldSidxMapping the SIDX to compare against
 */
export const filterChangedSidxMappings = (masterXml, srcUrl, clientOffset, oldSidxMapping) => {
  // Don't pass current sidx mapping
  const master = parseMpd(masterXml, {
    manifestUri: srcUrl,
    clientOffset
  });

  const videoSidx = compareSidxEntry(master.playlists, oldSidxMapping);
  let mediaGroupSidx = videoSidx;

  forEachMediaGroup(master, (properties, mediaType, groupKey, labelKey) => {
    if (properties.playlists && properties.playlists.length) {
      const playlists = properties.playlists;

      mediaGroupSidx = mergeOptions(
        mediaGroupSidx,
        compareSidxEntry(playlists, oldSidxMapping)
      );
    }
  });

  return mediaGroupSidx;
};

// exported for testing
export const requestSidx_ = (loader, sidxRange, playlist, xhr, options, finishProcessingFn) => {
  const sidxInfo = {
    // resolve the segment URL relative to the playlist
    uri: resolveManifestRedirect(options.handleManifestRedirects, sidxRange.resolvedUri),
    // resolvedUri: sidxRange.resolvedUri,
    byterange: sidxRange.byterange,
    // the segment's playlist
    playlist
  };

  const sidxRequestOptions = videojs.mergeOptions(sidxInfo, {
    responseType: 'arraybuffer',
    headers: segmentXhrHeaders(sidxInfo)
  });

  return containerRequest(sidxInfo.uri, xhr, (err, request, container, bytes) => {
    if (err) {
      return finishProcessingFn(err, request);
    }

    if (!container || container !== 'mp4') {
      return finishProcessingFn({
        status: request.status,
        message: `Unsupported ${container || 'unknown'} container type for sidx segment at URL: ${sidxInfo.uri}`,
        // response is just bytes in this case
        // but we really don't want to return that.
        response: '',
        playlist,
        internal: true,
        blacklistDuration: Infinity,
        // MEDIA_ERR_NETWORK
        code: 2
      }, request);
    }

    // if we already downloaded the sidx bytes in the container request, use them
    const {offset, length} = sidxInfo.byterange;

    if (bytes.length >= (length + offset)) {
      return finishProcessingFn(err, {
        response: bytes.subarray(offset, offset + length),
        status: request.status,
        uri: request.uri
      });
    }

    // otherwise request sidx bytes
    loader.request = xhr(sidxRequestOptions, finishProcessingFn);
  });
};

export default class DashPlaylistLoader extends EventTarget {
  // DashPlaylistLoader must accept either a src url or a playlist because subsequent
  // playlist loader setups from media groups will expect to be able to pass a playlist
  // (since there aren't external URLs to media playlists with DASH)
  constructor(srcUrlOrPlaylist, vhs, options = { }, masterPlaylistLoader) {
    super();

    const { withCredentials = false, handleManifestRedirects = false } = options;

    this.vhs_ = vhs;
    this.withCredentials = withCredentials;
    this.handleManifestRedirects = handleManifestRedirects;

    if (!srcUrlOrPlaylist) {
      throw new Error('A non-empty playlist URL or object is required');
    }

    // event naming?
    this.on('minimumUpdatePeriod', () => {
      this.refreshXml_();
    });

    // live playlist staleness timeout
    this.on('mediaupdatetimeout', () => {
      this.refreshMedia_(this.media().id);
    });

    this.state = 'HAVE_NOTHING';
    this.loadedPlaylists_ = {};

    // initialize the loader state
    // The masterPlaylistLoader will be created with a string
    if (typeof srcUrlOrPlaylist === 'string') {
      this.srcUrl = srcUrlOrPlaylist;
      // TODO: reset sidxMapping between period changes
      // once multi-period is refactored
      this.sidxMapping_ = {};
      return;
    }

    this.setupChildLoader(masterPlaylistLoader, srcUrlOrPlaylist);
  }

  setupChildLoader(masterPlaylistLoader, playlist) {
    this.masterPlaylistLoader_ = masterPlaylistLoader;
    this.childPlaylist_ = playlist;
  }

  dispose() {
    this.trigger('dispose');
    this.stopRequest();
    this.loadedPlaylists_ = {};
    window.clearTimeout(this.minimumUpdatePeriodTimeout_);
    window.clearTimeout(this.mediaRequest_);
    window.clearTimeout(this.mediaUpdateTimeout);

    this.off();
  }

  hasPendingRequest() {
    return this.request || this.mediaRequest_;
  }

  stopRequest() {
    if (this.request) {
      const oldRequest = this.request;

      this.request = null;
      oldRequest.onreadystatechange = null;
      oldRequest.abort();
    }
  }

  sidxRequestFinished_(playlist, master, startingState, doneFn) {
    return (err, request) => {
      // disposed
      if (!this.request) {
        return;
      }

      // pending request is cleared
      this.request = null;

      if (err) {
        // use the provided error or create one
        // see requestSidx_ for the container request
        // that can cause this.
        this.error = typeof err === 'object' ? err : {
          status: request.status,
          message: 'DASH playlist request error at URL: ' + playlist.uri,
          response: request.response,
          // MEDIA_ERR_NETWORK
          code: 2
        };
        if (startingState) {
          this.state = startingState;
        }

        this.trigger('error');
        return;
      }

      const bytes = toUint8(request.response);
      const sidx = parseSidx(bytes.subarray(8));

      return doneFn(master, sidx);
    };
  }

  media(playlist) {
    // getter
    if (!playlist) {
      return this.media_;
    }

    // setter
    if (this.state === 'HAVE_NOTHING') {
      throw new Error('Cannot switch media playlist from ' + this.state);
    }

    const startingState = this.state;

    // find the playlist object if the target playlist has been specified by URI
    if (typeof playlist === 'string') {
      if (!this.master.playlists[playlist]) {
        throw new Error('Unknown playlist URI: ' + playlist);
      }
      playlist = this.master.playlists[playlist];
    }

    const mediaChange = !this.media_ || playlist.id !== this.media_.id;

    // switch to previously loaded playlists immediately
    if (mediaChange &&
      this.loadedPlaylists_[playlist.id] &&
      this.loadedPlaylists_[playlist.id].endList) {
      this.state = 'HAVE_METADATA';
      this.media_ = playlist;

      // trigger media change if the active media has been updated
      if (mediaChange) {
        this.trigger('mediachanging');
        this.trigger('mediachange');
      }
      return;
    }

    // switching to the active playlist is a no-op
    if (!mediaChange) {
      return;
    }

    // switching from an already loaded playlist
    if (this.media_) {
      this.trigger('mediachanging');
    }

    if (!playlist.sidx) {
      // Continue asynchronously if there is no sidx
      // wait one tick to allow haveMaster to run first on a child loader
      this.mediaRequest_ = window.setTimeout(
        this.haveMetadata.bind(this, { startingState, playlist }),
        0
      );

      // exit early and don't do sidx work
      return;
    }

    // we have sidx mappings
    let oldMaster;
    let sidxMapping;

    // sidxMapping is used when parsing the masterXml, so store
    // it on the masterPlaylistLoader
    if (this.masterPlaylistLoader_) {
      oldMaster = this.masterPlaylistLoader_.master;
      sidxMapping = this.masterPlaylistLoader_.sidxMapping_;
    } else {
      oldMaster = this.master;
      sidxMapping = this.sidxMapping_;
    }

    const sidxKey = generateSidxKey(playlist.sidx);

    sidxMapping[sidxKey] = {
      sidxInfo: playlist.sidx
    };

    this.request = requestSidx_(
      this,
      playlist.sidx,
      playlist,
      this.vhs_.xhr,
      { handleManifestRedirects: this.handleManifestRedirects },
      this.sidxRequestFinished_(playlist, oldMaster, startingState, (newMaster, sidx) => {
        if (!newMaster || !sidx) {
          throw new Error('failed to request sidx');
        }

        // update loader's sidxMapping with parsed sidx box
        sidxMapping[sidxKey].sidx = sidx;

        // everything is ready just continue to haveMetadata
        this.haveMetadata({
          startingState,
          playlist: newMaster.playlists[playlist.id]
        });
      })
    );
  }

  haveMetadata({startingState, playlist}) {
    this.state = 'HAVE_METADATA';
    this.loadedPlaylists_[playlist.id] = playlist;
    this.mediaRequest_ = null;

    // This will trigger loadedplaylist
    this.refreshMedia_(playlist.id);

    // fire loadedmetadata the first time a media playlist is loaded
    // to resolve setup of media groups
    if (startingState === 'HAVE_MASTER') {
      this.trigger('loadedmetadata');
    } else {
      // trigger media change if the active media has been updated
      this.trigger('mediachange');
    }
  }

  pause() {
    this.stopRequest();
    window.clearTimeout(this.mediaUpdateTimeout);
    window.clearTimeout(this.minimumUpdatePeriodTimeout_);
    if (this.state === 'HAVE_NOTHING') {
      // If we pause the loader before any data has been retrieved, its as if we never
      // started, so reset to an unstarted state.
      this.started = false;
    }
  }

  load(isFinalRendition) {
    window.clearTimeout(this.mediaUpdateTimeout);
    window.clearTimeout(this.minimumUpdatePeriodTimeout_);

    const media = this.media();

    if (isFinalRendition) {
      const delay = media ? (media.targetDuration / 2) * 1000 : 5 * 1000;

      this.mediaUpdateTimeout = window.setTimeout(() => this.load(), delay);
      return;
    }

    // because the playlists are internal to the manifest, load should either load the
    // main manifest, or do nothing but trigger an event
    if (!this.started) {
      this.start();
      return;
    }

    if (media && !media.endList) {
      this.trigger('mediaupdatetimeout');
    } else {
      this.trigger('loadedplaylist');
    }
  }

  start() {
    this.started = true;

    // We don't need to request the master manifest again
    // Call this asynchronously to match the xhr request behavior below
    if (this.masterPlaylistLoader_) {
      this.mediaRequest_ = window.setTimeout(
        this.haveMaster_.bind(this),
        0
      );
      return;
    }

    // request the specified URL
    this.request = this.vhs_.xhr({
      uri: this.srcUrl,
      withCredentials: this.withCredentials
    }, (error, req) => {
      // disposed
      if (!this.request) {
        return;
      }

      // clear the loader's request reference
      this.request = null;

      if (error) {
        this.error = {
          status: req.status,
          message: 'DASH playlist request error at URL: ' + this.srcUrl,
          responseText: req.responseText,
          // MEDIA_ERR_NETWORK
          code: 2
        };
        if (this.state === 'HAVE_NOTHING') {
          this.started = false;
        }
        return this.trigger('error');
      }

      this.masterXml_ = req.responseText;

      if (req.responseHeaders && req.responseHeaders.date) {
        this.masterLoaded_ = Date.parse(req.responseHeaders.date);
      } else {
        this.masterLoaded_ = Date.now();
      }

      this.srcUrl = resolveManifestRedirect(this.handleManifestRedirects, this.srcUrl, req);

      this.syncClientServerClock_(this.onClientServerClockSync_.bind(this));
    });
  }

  /**
   * Parses the master xml for UTCTiming node to sync the client clock to the server
   * clock. If the UTCTiming node requires a HEAD or GET request, that request is made.
   *
   * @param {Function} done
   *        Function to call when clock sync has completed
   */
  syncClientServerClock_(done) {
    const utcTiming = parseUTCTiming(this.masterXml_);

    // No UTCTiming element found in the mpd. Use Date header from mpd request as the
    // server clock
    if (utcTiming === null) {
      this.clientOffset_ = this.masterLoaded_ - Date.now();
      return done();
    }

    if (utcTiming.method === 'DIRECT') {
      this.clientOffset_ = utcTiming.value - Date.now();
      return done();
    }

    this.request = this.vhs_.xhr({
      uri: resolveUrl(this.srcUrl, utcTiming.value),
      method: utcTiming.method,
      withCredentials: this.withCredentials
    }, (error, req) => {
      // disposed
      if (!this.request) {
        return;
      }

      if (error) {
        // sync request failed, fall back to using date header from mpd
        // TODO: log warning
        this.clientOffset_ = this.masterLoaded_ - Date.now();
        return done();
      }

      let serverTime;

      if (utcTiming.method === 'HEAD') {
        if (!req.responseHeaders || !req.responseHeaders.date) {
          // expected date header not preset, fall back to using date header from mpd
          // TODO: log warning
          serverTime = this.masterLoaded_;
        } else {
          serverTime = Date.parse(req.responseHeaders.date);
        }
      } else {
        serverTime = Date.parse(req.responseText);
      }

      this.clientOffset_ = serverTime - Date.now();

      done();
    });
  }

  haveMaster_() {
    this.state = 'HAVE_MASTER';
    // clear media request
    this.mediaRequest_ = null;

    if (!this.masterPlaylistLoader_) {
      this.updateMainManifest_(parseMasterXml({
        masterXml: this.masterXml_,
        srcUrl: this.srcUrl,
        clientOffset: this.clientOffset_,
        sidxMapping: this.sidxMapping_
      }));
      // We have the master playlist at this point, so
      // trigger this to allow MasterPlaylistController
      // to make an initial playlist selection
      this.trigger('loadedplaylist');
    } else if (!this.media_) {
      // no media playlist was specifically selected so select
      // the one the child playlist loader was created with
      this.media(this.childPlaylist_);
    }
  }

  updateMinimumUpdatePeriodTimeout_() {
    // Clear existing timeout
    window.clearTimeout(this.minimumUpdatePeriodTimeout_);

    const createMUPTimeout = (mup) => {
      this.minimumUpdatePeriodTimeout_ = window.setTimeout(() => {
        this.trigger('minimumUpdatePeriod');
      }, mup);
    };

    const minimumUpdatePeriod = this.master && this.master.minimumUpdatePeriod;

    if (minimumUpdatePeriod > 0) {
      createMUPTimeout(minimumUpdatePeriod);

    // If the minimumUpdatePeriod has a value of 0, that indicates that the current
    // MPD has no future validity, so a new one will need to be acquired when new
    // media segments are to be made available. Thus, we use the target duration
    // in this case
    } else if (minimumUpdatePeriod === 0) {
      // If we haven't yet selected a playlist, wait until then so we know the
      // target duration
      if (!this.media()) {
        this.one('loadedplaylist', () => {
          createMUPTimeout(this.media().targetDuration * 1000);
        });
      } else {
        createMUPTimeout(this.media().targetDuration * 1000);
      }
    }
  }

  /**
   * Handler for after client/server clock synchronization has happened. Sets up
   * xml refresh timer if specificed by the manifest.
   */
  onClientServerClockSync_() {
    this.haveMaster_();

    if (!this.hasPendingRequest() && !this.media_) {
      this.media(this.master.playlists[0]);
    }

    this.updateMinimumUpdatePeriodTimeout_();
  }

  /**
   * Given a new manifest, update our pointer to it and update the srcUrl based on the location elements of the manifest, if they exist.
   *
   * @param {Object} updatedManifest the manifest to update to
   */
  updateMainManifest_(updatedManifest) {
    this.master = updatedManifest;

    // if locations isn't set or is an empty array, exit early
    if (!this.master.locations || !this.master.locations.length) {
      return;
    }

    const location = this.master.locations[0];

    if (location !== this.srcUrl) {
      this.srcUrl = location;
    }
  }

  /**
   * Sends request to refresh the master xml and updates the parsed master manifest
   * TODO: Does the client offset need to be recalculated when the xml is refreshed?
   */
  refreshXml_() {
    // The srcUrl here *may* need to pass through handleManifestsRedirects when
    // sidx is implemented
    this.request = this.vhs_.xhr({
      uri: this.srcUrl,
      withCredentials: this.withCredentials
    }, (error, req) => {
      // disposed
      if (!this.request) {
        return;
      }

      // clear the loader's request reference
      this.request = null;

      if (error) {
        this.error = {
          status: req.status,
          message: 'DASH playlist request error at URL: ' + this.srcUrl,
          responseText: req.responseText,
          // MEDIA_ERR_NETWORK
          code: 2
        };
        if (this.state === 'HAVE_NOTHING') {
          this.started = false;
        }
        return this.trigger('error');
      }

      this.masterXml_ = req.responseText;

      // This will filter out updated sidx info from the mapping
      this.sidxMapping_ = filterChangedSidxMappings(
        this.masterXml_,
        this.srcUrl,
        this.clientOffset_,
        this.sidxMapping_
      );

      const master = parseMasterXml({
        masterXml: this.masterXml_,
        srcUrl: this.srcUrl,
        clientOffset: this.clientOffset_,
        sidxMapping: this.sidxMapping_
      });
      const updatedMaster = updateMaster(this.master, master);
      const currentSidxInfo = this.media().sidx;

      if (updatedMaster) {
        if (currentSidxInfo) {
          const sidxKey = generateSidxKey(currentSidxInfo);

          // the sidx was updated, so the previous mapping was removed
          if (!this.sidxMapping_[sidxKey]) {
            const playlist = this.media();

            this.request = requestSidx_(
              this,
              playlist.sidx,
              playlist,
              this.vhs_.xhr,
              { handleManifestRedirects: this.handleManifestRedirects },
              this.sidxRequestFinished_(playlist, master, this.state, (newMaster, sidx) => {
                if (!newMaster || !sidx) {
                  throw new Error('failed to request sidx on minimumUpdatePeriod');
                }

                // update loader's sidxMapping with parsed sidx box
                this.sidxMapping_[sidxKey].sidx = sidx;

                this.updateMinimumUpdatePeriodTimeout_();

                // TODO: do we need to reload the current playlist?
                this.refreshMedia_(this.media().id);

                return;
              })
            );
          }
        } else {
          this.updateMainManifest_(updatedMaster);
          if (this.media_) {
            this.media_ = this.master.playlists[this.media_.id];
          }
        }
      }

      this.updateMinimumUpdatePeriodTimeout_();
    });
  }

  /**
   * Refreshes the media playlist by re-parsing the master xml and updating playlist
   * references. If this is an alternate loader, the updated parsed manifest is retrieved
   * from the master loader.
   */
  refreshMedia_(mediaID) {
    if (!mediaID) {
      throw new Error('refreshMedia_ must take a media id');
    }

    let oldMaster;
    let newMaster;

    if (this.masterPlaylistLoader_) {
      oldMaster = this.masterPlaylistLoader_.master;
      newMaster = parseMasterXml({
        masterXml: this.masterPlaylistLoader_.masterXml_,
        srcUrl: this.masterPlaylistLoader_.srcUrl,
        clientOffset: this.masterPlaylistLoader_.clientOffset_,
        sidxMapping: this.masterPlaylistLoader_.sidxMapping_
      });
    } else {
      oldMaster = this.master;
      newMaster = parseMasterXml({
        masterXml: this.masterXml_,
        srcUrl: this.srcUrl,
        clientOffset: this.clientOffset_,
        sidxMapping: this.sidxMapping_
      });
    }

    const updatedMaster = updateMaster(oldMaster, newMaster);

    if (updatedMaster) {
      if (this.masterPlaylistLoader_) {
        this.masterPlaylistLoader_.master = updatedMaster;
      } else {
        this.master = updatedMaster;
      }
      this.media_ = updatedMaster.playlists[mediaID];
    } else {
      this.media_ = oldMaster.playlists[mediaID];
      this.trigger('playlistunchanged');
    }

    if (!this.media().endList) {
      this.mediaUpdateTimeout = window.setTimeout(() => {
        this.trigger('mediaupdatetimeout');
      }, refreshDelay(this.media(), !!updatedMaster));
    }

    this.trigger('loadedplaylist');
  }
}

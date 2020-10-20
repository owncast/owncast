import videojs from 'video.js';
import PlaylistLoader from './playlist-loader';
import DashPlaylistLoader from './dash-playlist-loader';
import noop from './util/noop';

/**
 * Convert the properties of an HLS track into an audioTrackKind.
 *
 * @private
 */
const audioTrackKind_ = (properties) => {
  let kind = properties.default ? 'main' : 'alternative';

  if (properties.characteristics &&
      properties.characteristics.indexOf('public.accessibility.describes-video') >= 0) {
    kind = 'main-desc';
  }

  return kind;
};

/**
 * Pause provided segment loader and playlist loader if active
 *
 * @param {SegmentLoader} segmentLoader
 *        SegmentLoader to pause
 * @param {Object} mediaType
 *        Active media type
 * @function stopLoaders
 */
export const stopLoaders = (segmentLoader, mediaType) => {
  segmentLoader.abort();
  segmentLoader.pause();

  if (mediaType && mediaType.activePlaylistLoader) {
    mediaType.activePlaylistLoader.pause();
    mediaType.activePlaylistLoader = null;
  }
};

/**
 * Start loading provided segment loader and playlist loader
 *
 * @param {PlaylistLoader} playlistLoader
 *        PlaylistLoader to start loading
 * @param {Object} mediaType
 *        Active media type
 * @function startLoaders
 */
export const startLoaders = (playlistLoader, mediaType) => {
  // Segment loader will be started after `loadedmetadata` or `loadedplaylist` from the
  // playlist loader
  mediaType.activePlaylistLoader = playlistLoader;
  playlistLoader.load();
};

/**
 * Returns a function to be called when the media group changes. It performs a
 * non-destructive (preserve the buffer) resync of the SegmentLoader. This is because a
 * change of group is merely a rendition switch of the same content at another encoding,
 * rather than a change of content, such as switching audio from English to Spanish.
 *
 * @param {string} type
 *        MediaGroup type
 * @param {Object} settings
 *        Object containing required information for media groups
 * @return {Function}
 *         Handler for a non-destructive resync of SegmentLoader when the active media
 *         group changes.
 * @function onGroupChanged
 */
export const onGroupChanged = (type, settings) => () => {
  const {
    segmentLoaders: {
      [type]: segmentLoader,
      main: mainSegmentLoader
    },
    mediaTypes: { [type]: mediaType }
  } = settings;
  const activeTrack = mediaType.activeTrack();
  const activeGroup = mediaType.activeGroup(activeTrack);
  const previousActiveLoader = mediaType.activePlaylistLoader;

  stopLoaders(segmentLoader, mediaType);

  if (!activeGroup) {
    // there is no group active
    return;
  }

  if (!activeGroup.playlistLoader) {
    if (previousActiveLoader) {
      // The previous group had a playlist loader but the new active group does not
      // this means we are switching from demuxed to muxed audio. In this case we want to
      // do a destructive reset of the main segment loader and not restart the audio
      // loaders.
      mainSegmentLoader.resetEverything();
    }
    return;
  }

  // Non-destructive resync
  segmentLoader.resyncLoader();

  startLoaders(activeGroup.playlistLoader, mediaType);
};

export const onGroupChanging = (type, settings) => () => {
  const {
    segmentLoaders: {
      [type]: segmentLoader
    }
  } = settings;

  segmentLoader.abort();
  segmentLoader.pause();
};

/**
 * Returns a function to be called when the media track changes. It performs a
 * destructive reset of the SegmentLoader to ensure we start loading as close to
 * currentTime as possible.
 *
 * @param {string} type
 *        MediaGroup type
 * @param {Object} settings
 *        Object containing required information for media groups
 * @return {Function}
 *         Handler for a destructive reset of SegmentLoader when the active media
 *         track changes.
 * @function onTrackChanged
 */
export const onTrackChanged = (type, settings) => () => {
  const {
    segmentLoaders: {
      [type]: segmentLoader,
      main: mainSegmentLoader
    },
    mediaTypes: { [type]: mediaType }
  } = settings;
  const activeTrack = mediaType.activeTrack();
  const activeGroup = mediaType.activeGroup(activeTrack);
  const previousActiveLoader = mediaType.activePlaylistLoader;

  stopLoaders(segmentLoader, mediaType);

  if (!activeGroup) {
    // there is no group active so we do not want to restart loaders
    return;
  }

  if (type === 'AUDIO') {
    if (!activeGroup.playlistLoader) {
      // when switching from demuxed audio/video to muxed audio/video (noted by no
      // playlist loader for the audio group), we want to do a destructive reset of the
      // main segment loader and not restart the audio loaders
      mainSegmentLoader.setAudio(true);
      // don't have to worry about disabling the audio of the audio segment loader since
      // it should be stopped
      mainSegmentLoader.resetEverything();
      return;
    }

    // although the segment loader is an audio segment loader, call the setAudio
    // function to ensure it is prepared to re-append the init segment (or handle other
    // config changes)
    segmentLoader.setAudio(true);
    mainSegmentLoader.setAudio(false);
  }

  if (previousActiveLoader === activeGroup.playlistLoader) {
    // Nothing has actually changed. This can happen because track change events can fire
    // multiple times for a "single" change. One for enabling the new active track, and
    // one for disabling the track that was active
    startLoaders(activeGroup.playlistLoader, mediaType);
    return;
  }

  if (segmentLoader.track) {
    // For WebVTT, set the new text track in the segmentloader
    segmentLoader.track(activeTrack);
  }

  // destructive reset
  segmentLoader.resetEverything();

  startLoaders(activeGroup.playlistLoader, mediaType);
};

export const onError = {
  /**
   * Returns a function to be called when a SegmentLoader or PlaylistLoader encounters
   * an error.
   *
   * @param {string} type
   *        MediaGroup type
   * @param {Object} settings
   *        Object containing required information for media groups
   * @return {Function}
   *         Error handler. Logs warning (or error if the playlist is blacklisted) to
   *         console and switches back to default audio track.
   * @function onError.AUDIO
   */
  AUDIO: (type, settings) => () => {
    const {
      segmentLoaders: { [type]: segmentLoader},
      mediaTypes: { [type]: mediaType },
      blacklistCurrentPlaylist
    } = settings;

    stopLoaders(segmentLoader, mediaType);

    // switch back to default audio track
    const activeTrack = mediaType.activeTrack();
    const activeGroup = mediaType.activeGroup();
    const id = (activeGroup.filter(group => group.default)[0] || activeGroup[0]).id;
    const defaultTrack = mediaType.tracks[id];

    if (activeTrack === defaultTrack) {
      // Default track encountered an error. All we can do now is blacklist the current
      // rendition and hope another will switch audio groups
      blacklistCurrentPlaylist({
        message: 'Problem encountered loading the default audio track.'
      });
      return;
    }

    videojs.log.warn('Problem encountered loading the alternate audio track.' +
                       'Switching back to default.');

    for (const trackId in mediaType.tracks) {
      mediaType.tracks[trackId].enabled = mediaType.tracks[trackId] === defaultTrack;
    }

    mediaType.onTrackChanged();
  },
  /**
   * Returns a function to be called when a SegmentLoader or PlaylistLoader encounters
   * an error.
   *
   * @param {string} type
   *        MediaGroup type
   * @param {Object} settings
   *        Object containing required information for media groups
   * @return {Function}
   *         Error handler. Logs warning to console and disables the active subtitle track
   * @function onError.SUBTITLES
   */
  SUBTITLES: (type, settings) => () => {
    const {
      segmentLoaders: { [type]: segmentLoader},
      mediaTypes: { [type]: mediaType }
    } = settings;

    videojs.log.warn('Problem encountered loading the subtitle track.' +
                     'Disabling subtitle track.');

    stopLoaders(segmentLoader, mediaType);

    const track = mediaType.activeTrack();

    if (track) {
      track.mode = 'disabled';
    }

    mediaType.onTrackChanged();
  }
};

export const setupListeners = {
  /**
   * Setup event listeners for audio playlist loader
   *
   * @param {string} type
   *        MediaGroup type
   * @param {PlaylistLoader|null} playlistLoader
   *        PlaylistLoader to register listeners on
   * @param {Object} settings
   *        Object containing required information for media groups
   * @function setupListeners.AUDIO
   */
  AUDIO: (type, playlistLoader, settings) => {
    if (!playlistLoader) {
      // no playlist loader means audio will be muxed with the video
      return;
    }

    const {
      tech,
      requestOptions,
      segmentLoaders: { [type]: segmentLoader }
    } = settings;

    playlistLoader.on('loadedmetadata', () => {
      const media = playlistLoader.media();

      segmentLoader.playlist(media, requestOptions);

      // if the video is already playing, or if this isn't a live video and preload
      // permits, start downloading segments
      if (!tech.paused() || (media.endList && tech.preload() !== 'none')) {
        segmentLoader.load();
      }
    });

    playlistLoader.on('loadedplaylist', () => {
      segmentLoader.playlist(playlistLoader.media(), requestOptions);

      // If the player isn't paused, ensure that the segment loader is running
      if (!tech.paused()) {
        segmentLoader.load();
      }
    });

    playlistLoader.on('error', onError[type](type, settings));
  },
  /**
   * Setup event listeners for subtitle playlist loader
   *
   * @param {string} type
   *        MediaGroup type
   * @param {PlaylistLoader|null} playlistLoader
   *        PlaylistLoader to register listeners on
   * @param {Object} settings
   *        Object containing required information for media groups
   * @function setupListeners.SUBTITLES
   */
  SUBTITLES: (type, playlistLoader, settings) => {
    const {
      tech,
      requestOptions,
      segmentLoaders: { [type]: segmentLoader },
      mediaTypes: { [type]: mediaType }
    } = settings;

    playlistLoader.on('loadedmetadata', () => {
      const media = playlistLoader.media();

      segmentLoader.playlist(media, requestOptions);
      segmentLoader.track(mediaType.activeTrack());

      // if the video is already playing, or if this isn't a live video and preload
      // permits, start downloading segments
      if (!tech.paused() || (media.endList && tech.preload() !== 'none')) {
        segmentLoader.load();
      }
    });

    playlistLoader.on('loadedplaylist', () => {
      segmentLoader.playlist(playlistLoader.media(), requestOptions);

      // If the player isn't paused, ensure that the segment loader is running
      if (!tech.paused()) {
        segmentLoader.load();
      }
    });

    playlistLoader.on('error', onError[type](type, settings));
  }
};

export const initialize = {
  /**
   * Setup PlaylistLoaders and AudioTracks for the audio groups
   *
   * @param {string} type
   *        MediaGroup type
   * @param {Object} settings
   *        Object containing required information for media groups
   * @function initialize.AUDIO
   */
  'AUDIO': (type, settings) => {
    const {
      vhs,
      sourceType,
      segmentLoaders: { [type]: segmentLoader },
      requestOptions,
      master: { mediaGroups, playlists },
      mediaTypes: {
        [type]: {
          groups,
          tracks
        }
      },
      masterPlaylistLoader
    } = settings;

    // force a default if we have none
    if (!mediaGroups[type] ||
        Object.keys(mediaGroups[type]).length === 0) {
      mediaGroups[type] = { main: { default: { default: true } } };
    }

    for (const groupId in mediaGroups[type]) {
      if (!groups[groupId]) {
        groups[groupId] = [];
      }

      // List of playlists that have an AUDIO attribute value matching the current
      // group ID
      const groupPlaylists = playlists.filter(playlist => {
        return playlist.attributes[type] === groupId;
      });

      for (const variantLabel in mediaGroups[type][groupId]) {
        let properties = mediaGroups[type][groupId][variantLabel];

        // List of playlists for the current group ID that have a matching uri with
        // this alternate audio variant
        const matchingPlaylists = groupPlaylists.filter(playlist => {
          return playlist.resolvedUri === properties.resolvedUri;
        });

        if (matchingPlaylists.length) {
          // If there is a playlist that has the same uri as this audio variant, assume
          // that the playlist is audio only. We delete the resolvedUri property here
          // to prevent a playlist loader from being created so that we don't have
          // both the main and audio segment loaders loading the same audio segments
          // from the same playlist.
          delete properties.resolvedUri;
        }

        let playlistLoader;

        // if vhs-json was provided as the source, and the media playlist was resolved,
        // use the resolved media playlist object
        if (sourceType === 'vhs-json' && properties.playlists) {
          playlistLoader = new PlaylistLoader(
            properties.playlists[0],
            vhs,
            requestOptions
          );
        } else if (properties.resolvedUri) {
          playlistLoader = new PlaylistLoader(
            properties.resolvedUri,
            vhs,
            requestOptions
          );
        } else if (properties.playlists && sourceType === 'dash') {
          playlistLoader = new DashPlaylistLoader(
            properties.playlists[0],
            vhs,
            requestOptions,
            masterPlaylistLoader
          );
        } else {
          // no resolvedUri means the audio is muxed with the video when using this
          // audio track
          playlistLoader = null;
        }

        properties = videojs.mergeOptions(
          { id: variantLabel, playlistLoader },
          properties
        );

        setupListeners[type](type, properties.playlistLoader, settings);

        groups[groupId].push(properties);

        if (typeof tracks[variantLabel] === 'undefined') {
          const track = new videojs.AudioTrack({
            id: variantLabel,
            kind: audioTrackKind_(properties),
            enabled: false,
            language: properties.language,
            default: properties.default,
            label: variantLabel
          });

          tracks[variantLabel] = track;
        }
      }
    }

    // setup single error event handler for the segment loader
    segmentLoader.on('error', onError[type](type, settings));
  },
  /**
   * Setup PlaylistLoaders and TextTracks for the subtitle groups
   *
   * @param {string} type
   *        MediaGroup type
   * @param {Object} settings
   *        Object containing required information for media groups
   * @function initialize.SUBTITLES
   */
  'SUBTITLES': (type, settings) => {
    const {
      tech,
      vhs,
      sourceType,
      segmentLoaders: { [type]: segmentLoader },
      requestOptions,
      master: { mediaGroups },
      mediaTypes: {
        [type]: {
          groups,
          tracks
        }
      },
      masterPlaylistLoader
    } = settings;

    for (const groupId in mediaGroups[type]) {
      if (!groups[groupId]) {
        groups[groupId] = [];
      }

      for (const variantLabel in mediaGroups[type][groupId]) {
        if (mediaGroups[type][groupId][variantLabel].forced) {
          // Subtitle playlists with the forced attribute are not selectable in Safari.
          // According to Apple's HLS Authoring Specification:
          //   If content has forced subtitles and regular subtitles in a given language,
          //   the regular subtitles track in that language MUST contain both the forced
          //   subtitles and the regular subtitles for that language.
          // Because of this requirement and that Safari does not add forced subtitles,
          // forced subtitles are skipped here to maintain consistent experience across
          // all platforms
          continue;
        }

        let properties = mediaGroups[type][groupId][variantLabel];

        let playlistLoader;

        if (sourceType === 'hls') {
          playlistLoader =
            new PlaylistLoader(properties.resolvedUri, vhs, requestOptions);
        } else if (sourceType === 'dash') {
          playlistLoader = new DashPlaylistLoader(
            properties.playlists[0],
            vhs,
            requestOptions,
            masterPlaylistLoader
          );
        } else if (sourceType === 'vhs-json') {
          playlistLoader = new PlaylistLoader(
            // if the vhs-json object included the media playlist, use the media playlist
            // as provided, otherwise use the resolved URI to load the playlist
            properties.playlists ? properties.playlists[0] : properties.resolvedUri,
            vhs,
            requestOptions
          );
        }

        properties = videojs.mergeOptions({
          id: variantLabel,
          playlistLoader
        }, properties);

        setupListeners[type](type, properties.playlistLoader, settings);

        groups[groupId].push(properties);

        if (typeof tracks[variantLabel] === 'undefined') {
          const track = tech.addRemoteTextTrack({
            id: variantLabel,
            kind: 'subtitles',
            default: properties.default && properties.autoselect,
            language: properties.language,
            label: variantLabel
          }, false).track;

          tracks[variantLabel] = track;
        }
      }
    }

    // setup single error event handler for the segment loader
    segmentLoader.on('error', onError[type](type, settings));
  },
  /**
   * Setup TextTracks for the closed-caption groups
   *
   * @param {String} type
   *        MediaGroup type
   * @param {Object} settings
   *        Object containing required information for media groups
   * @function initialize['CLOSED-CAPTIONS']
   */
  'CLOSED-CAPTIONS': (type, settings) => {
    const {
      tech,
      master: { mediaGroups },
      mediaTypes: {
        [type]: {
          groups,
          tracks
        }
      }
    } = settings;

    for (const groupId in mediaGroups[type]) {
      if (!groups[groupId]) {
        groups[groupId] = [];
      }

      for (const variantLabel in mediaGroups[type][groupId]) {
        const properties = mediaGroups[type][groupId][variantLabel];

        // We only support CEA608 captions for now, so ignore anything that
        // doesn't use a CCx INSTREAM-ID
        if (!properties.instreamId.match(/CC\d/)) {
          continue;
        }

        // No PlaylistLoader is required for Closed-Captions because the captions are
        // embedded within the video stream
        groups[groupId].push(videojs.mergeOptions({ id: variantLabel }, properties));

        if (typeof tracks[variantLabel] === 'undefined') {
          const track = tech.addRemoteTextTrack({
            id: properties.instreamId,
            kind: 'captions',
            default: properties.default && properties.autoselect,
            language: properties.language,
            label: variantLabel
          }, false).track;

          tracks[variantLabel] = track;
        }
      }
    }
  }
};

/**
 * Returns a function used to get the active group of the provided type
 *
 * @param {string} type
 *        MediaGroup type
 * @param {Object} settings
 *        Object containing required information for media groups
 * @return {Function}
 *         Function that returns the active media group for the provided type. Takes an
 *         optional parameter {TextTrack} track. If no track is provided, a list of all
 *         variants in the group, otherwise the variant corresponding to the provided
 *         track is returned.
 * @function activeGroup
 */
export const activeGroup = (type, settings) => (track) => {
  const {
    masterPlaylistLoader,
    mediaTypes: { [type]: { groups } }
  } = settings;

  const media = masterPlaylistLoader.media();

  if (!media) {
    return null;
  }

  let variants = null;

  if (media.attributes[type]) {
    variants = groups[media.attributes[type]];
  }

  variants = variants || groups.main;

  if (typeof track === 'undefined') {
    return variants;
  }

  if (track === null) {
    // An active track was specified so a corresponding group is expected. track === null
    // means no track is currently active so there is no corresponding group
    return null;
  }

  return variants.filter((props) => props.id === track.id)[0] || null;
};

export const activeTrack = {
  /**
   * Returns a function used to get the active track of type provided
   *
   * @param {string} type
   *        MediaGroup type
   * @param {Object} settings
   *        Object containing required information for media groups
   * @return {Function}
   *         Function that returns the active media track for the provided type. Returns
   *         null if no track is active
   * @function activeTrack.AUDIO
   */
  AUDIO: (type, settings) => () => {
    const { mediaTypes: { [type]: { tracks } } } = settings;

    for (const id in tracks) {
      if (tracks[id].enabled) {
        return tracks[id];
      }
    }

    return null;
  },
  /**
   * Returns a function used to get the active track of type provided
   *
   * @param {string} type
   *        MediaGroup type
   * @param {Object} settings
   *        Object containing required information for media groups
   * @return {Function}
   *         Function that returns the active media track for the provided type. Returns
   *         null if no track is active
   * @function activeTrack.SUBTITLES
   */
  SUBTITLES: (type, settings) => () => {
    const { mediaTypes: { [type]: { tracks } } } = settings;

    for (const id in tracks) {
      if (tracks[id].mode === 'showing' || tracks[id].mode === 'hidden') {
        return tracks[id];
      }
    }

    return null;
  }
};

/**
 * Setup PlaylistLoaders and Tracks for media groups (Audio, Subtitles,
 * Closed-Captions) specified in the master manifest.
 *
 * @param {Object} settings
 *        Object containing required information for setting up the media groups
 * @param {Tech} settings.tech
 *        The tech of the player
 * @param {Object} settings.requestOptions
 *        XHR request options used by the segment loaders
 * @param {PlaylistLoader} settings.masterPlaylistLoader
 *        PlaylistLoader for the master source
 * @param {VhsHandler} settings.vhs
 *        VHS SourceHandler
 * @param {Object} settings.master
 *        The parsed master manifest
 * @param {Object} settings.mediaTypes
 *        Object to store the loaders, tracks, and utility methods for each media type
 * @param {Function} settings.blacklistCurrentPlaylist
 *        Blacklists the current rendition and forces a rendition switch.
 * @function setupMediaGroups
 */
export const setupMediaGroups = (settings) => {
  ['AUDIO', 'SUBTITLES', 'CLOSED-CAPTIONS'].forEach((type) => {
    initialize[type](type, settings);
  });

  const {
    mediaTypes,
    masterPlaylistLoader,
    tech,
    vhs
  } = settings;

  // setup active group and track getters and change event handlers
  ['AUDIO', 'SUBTITLES'].forEach((type) => {
    mediaTypes[type].activeGroup = activeGroup(type, settings);
    mediaTypes[type].activeTrack = activeTrack[type](type, settings);
    mediaTypes[type].onGroupChanged = onGroupChanged(type, settings);
    mediaTypes[type].onGroupChanging = onGroupChanging(type, settings);
    mediaTypes[type].onTrackChanged = onTrackChanged(type, settings);
  });

  // DO NOT enable the default subtitle or caption track.
  // DO enable the default audio track
  const audioGroup = mediaTypes.AUDIO.activeGroup();

  if (audioGroup) {
    const groupId = (audioGroup.filter(group => group.default)[0] || audioGroup[0]).id;

    mediaTypes.AUDIO.tracks[groupId].enabled = true;
    mediaTypes.AUDIO.onTrackChanged();
  }

  masterPlaylistLoader.on('mediachange', () => {
    ['AUDIO', 'SUBTITLES'].forEach(type => mediaTypes[type].onGroupChanged());
  });

  masterPlaylistLoader.on('mediachanging', () => {
    ['AUDIO', 'SUBTITLES'].forEach(type => mediaTypes[type].onGroupChanging());
  });

  // custom audio track change event handler for usage event
  const onAudioTrackChanged = () => {
    mediaTypes.AUDIO.onTrackChanged();
    tech.trigger({ type: 'usage', name: 'vhs-audio-change' });
    tech.trigger({ type: 'usage', name: 'hls-audio-change' });
  };

  tech.audioTracks().addEventListener('change', onAudioTrackChanged);
  tech.remoteTextTracks().addEventListener(
    'change',
    mediaTypes.SUBTITLES.onTrackChanged
  );

  vhs.on('dispose', () => {
    tech.audioTracks().removeEventListener('change', onAudioTrackChanged);
    tech.remoteTextTracks().removeEventListener(
      'change',
      mediaTypes.SUBTITLES.onTrackChanged
    );
  });

  // clear existing audio tracks and add the ones we just created
  tech.clearTracks('audio');

  for (const id in mediaTypes.AUDIO.tracks) {
    tech.audioTracks().addTrack(mediaTypes.AUDIO.tracks[id]);
  }
};

/**
 * Creates skeleton object used to store the loaders, tracks, and utility methods for each
 * media type
 *
 * @return {Object}
 *         Object to store the loaders, tracks, and utility methods for each media type
 * @function createMediaTypes
 */
export const createMediaTypes = () => {
  const mediaTypes = {};

  ['AUDIO', 'SUBTITLES', 'CLOSED-CAPTIONS'].forEach((type) => {
    mediaTypes[type] = {
      groups: {},
      tracks: {},
      activePlaylistLoader: null,
      activeGroup: noop,
      activeTrack: noop,
      onGroupChanged: noop,
      onTrackChanged: noop
    };
  });

  return mediaTypes;
};

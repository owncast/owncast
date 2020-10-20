import { values } from './utils/object';
import { findIndexes } from './utils/list';
import { addSegmentsToPlaylist } from './segment/segmentBase';
import { byteRangeToString } from './segment/urlType';

const mergeDiscontiguousPlaylists = playlists => {
  const mergedPlaylists = values(playlists.reduce((acc, playlist) => {
    // assuming playlist IDs are the same across periods
    // TODO: handle multiperiod where representation sets are not the same
    // across periods
    const name = playlist.attributes.id + (playlist.attributes.lang || '');

    // Periods after first
    if (acc[name]) {
      // first segment of subsequent periods signal a discontinuity
      if (playlist.segments[0]) {
        playlist.segments[0].discontinuity = true;
      }
      acc[name].segments.push(...playlist.segments);

      // bubble up contentProtection, this assumes all DRM content
      // has the same contentProtection
      if (playlist.attributes.contentProtection) {
        acc[name].attributes.contentProtection =
          playlist.attributes.contentProtection;
      }
    } else {
      // first Period
      acc[name] = playlist;
    }

    return acc;
  }, {}));

  return mergedPlaylists.map(playlist => {
    playlist.discontinuityStarts =
        findIndexes(playlist.segments, 'discontinuity');

    return playlist;
  });
};

const addSegmentInfoFromSidx = (playlists, sidxMapping = {}) => {
  if (!Object.keys(sidxMapping).length) {
    return playlists;
  }

  for (const i in playlists) {
    const playlist = playlists[i];

    if (!playlist.sidx) {
      continue;
    }

    const sidxKey = playlist.sidx.uri + '-' +
      byteRangeToString(playlist.sidx.byterange);
    const sidxMatch = sidxMapping[sidxKey] && sidxMapping[sidxKey].sidx;

    if (playlist.sidx && sidxMatch) {
      addSegmentsToPlaylist(playlist, sidxMatch, playlist.sidx.resolvedUri);
    }
  }

  return playlists;
};

export const formatAudioPlaylist = ({ attributes, segments, sidx }) => {
  const playlist = {
    attributes: {
      NAME: attributes.id,
      BANDWIDTH: attributes.bandwidth,
      CODECS: attributes.codecs,
      ['PROGRAM-ID']: 1
    },
    uri: '',
    endList: (attributes.type || 'static') === 'static',
    timeline: attributes.periodIndex,
    resolvedUri: '',
    targetDuration: attributes.duration,
    segments,
    mediaSequence: segments.length ? segments[0].number : 1
  };

  if (attributes.contentProtection) {
    playlist.contentProtection = attributes.contentProtection;
  }

  if (sidx) {
    playlist.sidx = sidx;
  }

  return playlist;
};

export const formatVttPlaylist = ({ attributes, segments }) => {
  if (typeof segments === 'undefined') {
    // vtt tracks may use single file in BaseURL
    segments = [{
      uri: attributes.baseUrl,
      timeline: attributes.periodIndex,
      resolvedUri: attributes.baseUrl || '',
      duration: attributes.sourceDuration,
      number: 0
    }];
    // targetDuration should be the same duration as the only segment
    attributes.duration = attributes.sourceDuration;
  }
  return {
    attributes: {
      NAME: attributes.id,
      BANDWIDTH: attributes.bandwidth,
      ['PROGRAM-ID']: 1
    },
    uri: '',
    endList: (attributes.type || 'static') === 'static',
    timeline: attributes.periodIndex,
    resolvedUri: attributes.baseUrl || '',
    targetDuration: attributes.duration,
    segments,
    mediaSequence: segments.length ? segments[0].number : 1
  };
};

export const organizeAudioPlaylists = (playlists, sidxMapping = {}) => {
  let mainPlaylist;

  const formattedPlaylists = playlists.reduce((a, playlist) => {
    const role = playlist.attributes.role &&
      playlist.attributes.role.value || '';
    const language = playlist.attributes.lang || '';

    let label = 'main';

    if (language) {
      const roleLabel = role ? ` (${role})` : '';

      label = `${playlist.attributes.lang}${roleLabel}`;
    }

    // skip if we already have the highest quality audio for a language
    if (a[label] &&
      a[label].playlists[0].attributes.BANDWIDTH >
      playlist.attributes.bandwidth) {
      return a;
    }

    a[label] = {
      language,
      autoselect: true,
      default: role === 'main',
      playlists: addSegmentInfoFromSidx(
        [formatAudioPlaylist(playlist)],
        sidxMapping
      ),
      uri: ''
    };

    if (typeof mainPlaylist === 'undefined' && role === 'main') {
      mainPlaylist = playlist;
      mainPlaylist.default = true;
    }

    return a;
  }, {});

  // if no playlists have role "main", mark the first as main
  if (!mainPlaylist) {
    const firstLabel = Object.keys(formattedPlaylists)[0];

    formattedPlaylists[firstLabel].default = true;
  }

  return formattedPlaylists;
};

export const organizeVttPlaylists = (playlists, sidxMapping = {}) => {
  return playlists.reduce((a, playlist) => {
    const label = playlist.attributes.lang || 'text';

    // skip if we already have subtitles
    if (a[label]) {
      return a;
    }

    a[label] = {
      language: label,
      default: false,
      autoselect: false,
      playlists: addSegmentInfoFromSidx(
        [formatVttPlaylist(playlist)],
        sidxMapping
      ),
      uri: ''
    };

    return a;
  }, {});
};

export const formatVideoPlaylist = ({ attributes, segments, sidx }) => {
  const playlist = {
    attributes: {
      NAME: attributes.id,
      AUDIO: 'audio',
      SUBTITLES: 'subs',
      RESOLUTION: {
        width: attributes.width,
        height: attributes.height
      },
      CODECS: attributes.codecs,
      BANDWIDTH: attributes.bandwidth,
      ['PROGRAM-ID']: 1
    },
    uri: '',
    endList: (attributes.type || 'static') === 'static',
    timeline: attributes.periodIndex,
    resolvedUri: '',
    targetDuration: attributes.duration,
    segments,
    mediaSequence: segments.length ? segments[0].number : 1
  };

  if (attributes.contentProtection) {
    playlist.contentProtection = attributes.contentProtection;
  }

  if (sidx) {
    playlist.sidx = sidx;
  }

  return playlist;
};

export const toM3u8 = (dashPlaylists, sidxMapping = {}) => {
  if (!dashPlaylists.length) {
    return {};
  }

  // grab all master attributes
  const {
    sourceDuration: duration,
    type = 'static',
    suggestedPresentationDelay,
    minimumUpdatePeriod = 0
  } = dashPlaylists[0].attributes;

  const videoOnly = ({ attributes }) =>
    attributes.mimeType === 'video/mp4' || attributes.contentType === 'video';
  const audioOnly = ({ attributes }) =>
    attributes.mimeType === 'audio/mp4' || attributes.contentType === 'audio';
  const vttOnly = ({ attributes }) =>
    attributes.mimeType === 'text/vtt' || attributes.contentType === 'text';

  const videoPlaylists = mergeDiscontiguousPlaylists(dashPlaylists.filter(videoOnly)).map(formatVideoPlaylist);
  const audioPlaylists = mergeDiscontiguousPlaylists(dashPlaylists.filter(audioOnly));
  const vttPlaylists = dashPlaylists.filter(vttOnly);

  const master = {
    allowCache: true,
    discontinuityStarts: [],
    segments: [],
    endList: true,
    mediaGroups: {
      AUDIO: {},
      VIDEO: {},
      ['CLOSED-CAPTIONS']: {},
      SUBTITLES: {}
    },
    uri: '',
    duration,
    playlists: addSegmentInfoFromSidx(videoPlaylists, sidxMapping),
    minimumUpdatePeriod: minimumUpdatePeriod * 1000
  };

  if (type === 'dynamic') {
    master.suggestedPresentationDelay = suggestedPresentationDelay;
  }

  if (audioPlaylists.length) {
    master.mediaGroups.AUDIO.audio = organizeAudioPlaylists(audioPlaylists, sidxMapping);
  }

  if (vttPlaylists.length) {
    master.mediaGroups.SUBTITLES.subs = organizeVttPlaylists(vttPlaylists, sidxMapping);
  }

  return master;
};

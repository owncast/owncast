import videojs from 'video.js';
import window from 'global/window';
import { Parser as M3u8Parser } from 'm3u8-parser';
import { resolveUrl } from './resolve-url';

const { log } = videojs;

export const createPlaylistID = (index, uri) => {
  return `${index}-${uri}`;
};

/**
 * Parses a given m3u8 playlist
 *
 * @param {string} manifestString
 *        The downloaded manifest string
 * @param {Object[]} [customTagParsers]
 *        An array of custom tag parsers for the m3u8-parser instance
 * @param {Object[]} [customTagMappers]
 *         An array of custom tag mappers for the m3u8-parser instance
 * @return {Object}
 *         The manifest object
 */
export const parseManifest = ({
  manifestString,
  customTagParsers = [],
  customTagMappers = []
}) => {
  const parser = new M3u8Parser();

  customTagParsers.forEach(customParser => parser.addParser(customParser));
  customTagMappers.forEach(mapper => parser.addTagMapper(mapper));

  parser.push(manifestString);
  parser.end();

  return parser.manifest;
};

/**
 * Loops through all supported media groups in master and calls the provided
 * callback for each group
 *
 * @param {Object} master
 *        The parsed master manifest object
 * @param {Function} callback
 *        Callback to call for each media group
 */
export const forEachMediaGroup = (master, callback) => {
  ['AUDIO', 'SUBTITLES'].forEach((mediaType) => {
    for (const groupKey in master.mediaGroups[mediaType]) {
      for (const labelKey in master.mediaGroups[mediaType][groupKey]) {
        const mediaProperties = master.mediaGroups[mediaType][groupKey][labelKey];

        callback(mediaProperties, mediaType, groupKey, labelKey);
      }
    }
  });
};

/**
 * Adds properties and attributes to the playlist to keep consistent functionality for
 * playlists throughout VHS.
 *
 * @param {Object} config
 *        Arguments object
 * @param {Object} config.playlist
 *        The media playlist
 * @param {string} [config.uri]
 *        The uri to the media playlist (if media playlist is not from within a master
 *        playlist)
 * @param {string} id
 *        ID to use for the playlist
 */
export const setupMediaPlaylist = ({ playlist, uri, id }) => {
  playlist.id = id;

  if (uri) {
    // For media playlists, m3u8-parser does not have access to a URI, as HLS media
    // playlists do not contain their own source URI, but one is needed for consistency in
    // VHS.
    playlist.uri = uri;
  }

  // For HLS master playlists, even though certain attributes MUST be defined, the
  // stream may still be played without them.
  // For HLS media playlists, m3u8-parser does not attach an attributes object to the
  // manifest.
  //
  // To avoid undefined reference errors through the project, and make the code easier
  // to write/read, add an empty attributes object for these cases.
  playlist.attributes = playlist.attributes || {};
};

/**
 * Adds ID, resolvedUri, and attributes properties to each playlist of the master, where
 * necessary. In addition, creates playlist IDs for each playlist and adds playlist ID to
 * playlist references to the playlists array.
 *
 * @param {Object} master
 *        The master playlist
 */
export const setupMediaPlaylists = (master) => {
  let i = master.playlists.length;

  while (i--) {
    const playlist = master.playlists[i];

    setupMediaPlaylist({
      playlist,
      id: createPlaylistID(i, playlist.uri)
    });
    playlist.resolvedUri = resolveUrl(master.uri, playlist.uri);
    master.playlists[playlist.id] = playlist;
    // URI reference added for backwards compatibility
    master.playlists[playlist.uri] = playlist;

    // Although the spec states an #EXT-X-STREAM-INF tag MUST have a BANDWIDTH attribute,
    // the stream can be played without it. Although an attributes property may have been
    // added to the playlist to prevent undefined references, issue a warning to fix the
    // manifest.
    if (!playlist.attributes.BANDWIDTH) {
      log.warn('Invalid playlist STREAM-INF detected. Missing BANDWIDTH attribute.');
    }
  }
};

/**
 * Adds resolvedUri properties to each media group.
 *
 * @param {Object} master
 *        The master playlist
 */
export const resolveMediaGroupUris = (master) => {
  forEachMediaGroup(master, (properties) => {
    if (properties.uri) {
      properties.resolvedUri = resolveUrl(master.uri, properties.uri);
    }
  });
};

/**
 * Creates a master playlist wrapper to insert a sole media playlist into.
 *
 * @param {Object} media
 *        Media playlist
 * @param {string} uri
 *        The media URI
 *
 * @return {Object}
 *         Master playlist
 */
export const masterForMedia = (media, uri) => {
  const id = createPlaylistID(0, uri);
  const master = {
    mediaGroups: {
      'AUDIO': {},
      'VIDEO': {},
      'CLOSED-CAPTIONS': {},
      'SUBTITLES': {}
    },
    uri: window.location.href,
    resolvedUri: window.location.href,
    playlists: [{
      uri,
      id,
      resolvedUri: uri,
      // m3u8-parser does not attach an attributes property to media playlists so make
      // sure that the property is attached to avoid undefined reference errors
      attributes: {}
    }]
  };

  // set up ID reference
  master.playlists[id] = master.playlists[0];
  // URI reference added for backwards compatibility
  master.playlists[uri] = master.playlists[0];

  return master;
};

/**
 * Does an in-place update of the master manifest to add updated playlist URI references
 * as well as other properties needed by VHS that aren't included by the parser.
 *
 * @param {Object} master
 *        Master manifest object
 * @param {string} uri
 *        The source URI
 */
export const addPropertiesToMaster = (master, uri) => {
  master.uri = uri;

  for (let i = 0; i < master.playlists.length; i++) {
    if (!master.playlists[i].uri) {
      // Set up phony URIs for the playlists since playlists are referenced by their URIs
      // throughout VHS, but some formats (e.g., DASH) don't have external URIs
      // TODO: consider adding dummy URIs in mpd-parser
      const phonyUri = `placeholder-uri-${i}`;

      master.playlists[i].uri = phonyUri;
    }
  }

  forEachMediaGroup(master, (properties, mediaType, groupKey, labelKey) => {
    if (!properties.playlists ||
        !properties.playlists.length ||
        properties.playlists[0].uri) {
      return;
    }

    // Set up phony URIs for the media group playlists since playlists are referenced by
    // their URIs throughout VHS, but some formats (e.g., DASH) don't have external URIs
    const phonyUri = `placeholder-uri-${mediaType}-${groupKey}-${labelKey}`;
    const id = createPlaylistID(0, phonyUri);

    properties.playlists[0].uri = phonyUri;
    properties.playlists[0].id = id;
    // setup ID and URI references (URI for backwards compatibility)
    master.playlists[id] = properties.playlists[0];
    master.playlists[phonyUri] = properties.playlists[0];
  });

  setupMediaPlaylists(master);
  resolveMediaGroupUris(master);
};

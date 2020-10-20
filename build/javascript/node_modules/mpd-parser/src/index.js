import { version } from '../package.json';
import { toM3u8 } from './toM3u8';
import { toPlaylists } from './toPlaylists';
import { inheritAttributes } from './inheritAttributes';
import { stringToMpdXml } from './stringToMpdXml';
import { parseUTCTimingScheme } from './parseUTCTimingScheme';

const VERSION = version;

const parse = (manifestString, options = {}) => {
  const parsedManifestInfo = inheritAttributes(stringToMpdXml(manifestString), options);
  const playlists = toPlaylists(parsedManifestInfo.representationInfo);

  return toM3u8(playlists, parsedManifestInfo.locations, options.sidxMapping);
};

/**
 * Parses the manifest for a UTCTiming node, returning the nodes attributes if found
 *
 * @param {string} manifestString
 *        XML string of the MPD manifest
 * @return {Object|null}
 *         Attributes of UTCTiming node specified in the manifest. Null if none found
 */
const parseUTCTiming = (manifestString) =>
  parseUTCTimingScheme(stringToMpdXml(manifestString));

export {
  VERSION,
  parse,
  parseUTCTiming,
  stringToMpdXml,
  inheritAttributes,
  toPlaylists,
  toM3u8
};

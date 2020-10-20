import {DOMParser} from 'xmldom';
import errors from './errors';

export const stringToMpdXml = (manifestString) => {
  if (manifestString === '') {
    throw new Error(errors.DASH_EMPTY_MANIFEST);
  }

  const parser = new DOMParser();
  const xml = parser.parseFromString(manifestString, 'application/xml');
  const mpd = xml && xml.documentElement.tagName === 'MPD' ?
    xml.documentElement : null;

  if (!mpd || mpd &&
      mpd.getElementsByTagName('parsererror').length > 0) {
    throw new Error(errors.DASH_INVALID_XML);
  }

  return mpd;
};

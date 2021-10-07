const m3u8Parser = require('m3u8-parser');
const fetch = require('node-fetch');
const url = require('url');
const { test } = require('@jest/globals');

const HLS_SUBDIRECTORY = '/hls/';
const PLAYLIST_NAME = 'stream.m3u8';
const TEST_OWNCAST_INSTANCE = 'http://localhost:8080';
const HLS_FETCH_ITERATIONS = 5;

jest.setTimeout(40000);

async function getPlaylist(urlString) {
  const response = await fetch(urlString);
  expect(response.status).toBe(200);
  const body = await response.text();

  var parser = new m3u8Parser.Parser();

  parser.push(body);
  parser.end();

  return parser.manifest;
}

function normalizeUrl(urlString, baseUrl) {
  let parsedString = url.parse(urlString);
  if (!parsedString.host) {
    const testInstanceRoot = url.parse(baseUrl);
    parsedString.protocol = testInstanceRoot.protocol;
    parsedString.host = testInstanceRoot.host;

    const filename = baseUrl.substring(baseUrl.lastIndexOf('/') + 1);
    parsedString.pathname =
      testInstanceRoot.pathname.replace(filename, '') + urlString;
  }
  return url.format(parsedString).toString();
}

// Iterate over an array of video segments and make sure they return back
// valid status.
async function validateSegments(segments) {
  for (let segment of segments) {
    const res = await fetch(segment);
    expect(res.status).toBe(200);
  }
}

describe('fetch and parse HLS', () => {
  const masterPlaylistUrl = `${TEST_OWNCAST_INSTANCE}${HLS_SUBDIRECTORY}${PLAYLIST_NAME}`;
  var masterPlaylist;
  var mediaPlaylistUrl;

  test('fetch master playlist', async (done) => {
    try {
      masterPlaylist = await getPlaylist(masterPlaylistUrl);
    } catch (e) {
      console.error('error fetching and parsing master playlist', e);
    }

    done();
  });

  test('verify there is a media playlist', () => {
    // Master playlist should have at least one media playlist.
    expect(masterPlaylist.playlists.length).toBe(1);

    try {
      mediaPlaylistUrl = normalizeUrl(
        masterPlaylist.playlists[0].uri,
        masterPlaylistUrl
      );
    } catch (e) {
      console.error('error fetching and parsing media playlist', e);
    }
  });

  test('verify there are segments', async (done) => {
    let playlist;
    try {
      playlist = await getPlaylist(mediaPlaylistUrl);
    } catch (e) {
      console.error('error verifying segments in media playlist', e);
    }

    const segments = playlist.segments;
    expect(segments.length).toBeGreaterThan(0);

    done();
  });

  // Iterate over segments and make sure they change.
  // Use the reported duration of the segment to wait to
  // fetch another just like a real HLS player would do.
  var lastSegmentUrl;
  for (let i = 0; i < HLS_FETCH_ITERATIONS; i++) {
    test('fetch and monitor media playlist segments ' + i, async (done) => {
      await new Promise((r) => setTimeout(r, 3000));

      try {
        var playlist = await getPlaylist(mediaPlaylistUrl);
      } catch (e) {
        console.error('error updating media playlist', mediaPlaylistUrl, e);
      }

      const segments = playlist.segments;
      const segment = segments[segments.length - 1];
      expect(segment.uri).not.toBe(lastSegmentUrl);

      try {
        var segmentUrl = normalizeUrl(segment.uri, mediaPlaylistUrl);
        await validateSegments([segmentUrl]);
      } catch (e) {
        console.error('unable to validate HLS segment', segmentUrl, e);
      }

      lastSegmentUrl = segment.uri;

      done();
    });
  }
});

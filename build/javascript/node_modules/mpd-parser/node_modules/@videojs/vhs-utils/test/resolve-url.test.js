import QUnit from 'qunit';
import window from 'global/window';
import resolveUrl from '../src/resolve-url';

// A modified subset of tests from https://github.com/tjenkinson/url-toolkit

QUnit.module('URL resolver');

QUnit.test('works with a selection of valid urls', function(assert) {
  let currentLocation = '';

  if (window.location && window.location.protocol) {
    currentLocation = window.location.protocol + '//' + window.location.host;
  }

  assert.equal(
    resolveUrl('http://a.com/b/cd/e.m3u8', 'https://example.com/z.ts'),
    'https://example.com/z.ts'
  );
  assert.equal(resolveUrl('http://a.com/b/cd/e.m3u8', 'z.ts'), 'http://a.com/b/cd/z.ts');
  assert.equal(resolveUrl('//a.com/b/cd/e.m3u8', 'z.ts'), '//a.com/b/cd/z.ts');
  assert.equal(
    resolveUrl('/a/b/cd/e.m3u8', 'https://example.com:8080/z.ts'),
    'https://example.com:8080/z.ts'
  );
  assert.equal(resolveUrl('/a/b/cd/e.m3u8', 'z.ts'), currentLocation + '/a/b/cd/z.ts');
  assert.equal(resolveUrl('/a/b/cd/e.m3u8', '../../../z.ts'), currentLocation + '/z.ts');
});

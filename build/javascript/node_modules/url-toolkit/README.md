[![npm version](https://badge.fury.io/js/url-toolkit.svg)](https://badge.fury.io/js/url-toolkit)
[![Build Status](https://travis-ci.org/tjenkinson/url-toolkit.svg?branch=master)](https://travis-ci.org/tjenkinson/url-toolkit)

# URL Toolkit

Lightweight library to build an absolute URL from a base URL and a relative URL, written from [the spec (RFC 1808)](https://tools.ietf.org/html/rfc1808). Initially part of [HLS.JS](https://github.com/dailymotion/hls.js).

## Methods

### `buildAbsoluteURL(baseURL, relativeURL, opts={})`

Build an absolute URL from a relative and base one.

```javascript
URLToolkit.buildAbsoluteURL('http://a.com/b/cd', 'e/f/../g'); // => http://a.com/b/e/g
```

If you want to ensure that the URL is treated as a relative one you should prefix it with `./`.

```javascript
URLToolkit.buildAbsoluteURL('http://a.com/b/cd', 'a:b'); // => a:b
URLToolkit.buildAbsoluteURL('http://a.com/b/cd', './a:b'); // => http://a.com/b/a:b
```

By default the paths will not be normalized unless necessary, according to the spec. However you can ensure paths are always normalized by setting the `opts.alwaysNormalize` option to `true`.

```javascript
URLToolkit.buildAbsoluteURL('http://a.com/b/cd', '/e/f/../g'); // => http://a.com/e/f/../g
URLToolkit.buildAbsoluteURL('http://a.com/b/cd', '/e/f/../g', {
  alwaysNormalize: true,
}); // => http://a.com/e/g
```

### `normalizePath(url)`

Normalizes a path.

```javascript
URLToolkit.normalizePath('a/b/../c'); // => a/c
```

### `parseURL(url)`

Parse a URL into its separate components.

```javascript
URLToolkit.parseURL('http://a/b/c/d;p?q#f'); // =>
/* {
	scheme: 'http:',
	netLoc: '//a',
	path: '/b/c/d',
	params: ';p',
	query: '?q',
	fragment: '#f'
} */
```

### `buildURLFromParts(parts)`

Puts all the parts from `parseURL()` back together into a string.

## Example

```javascript
var URLToolkit = require('url-toolkit');
var url = URLToolkit.buildAbsoluteURL(
  'https://a.com/b/cd/e.m3u8?test=1#something',
  '../z.ts?abc=1#test'
);
console.log(url); // 'https://a.com/b/z.ts?abc=1#test'
```

## Browser

This can also be used in the browser thanks to [jsDelivr](https://github.com/jsdelivr/jsdelivr):

```html
<head>
  <script
    type="text/javascript"
    src="https://cdn.jsdelivr.net/npm/url-toolkit@2"
  ></script>
  <script type="text/javascript">
    var url = URLToolkit.buildAbsoluteURL(
      'https://a.com/b/cd/e.m3u8?test=1#something',
      '../z.ts?abc=1#test'
    );
    console.log(url); // 'https://a.com/b/z.ts?abc=1#test'
  </script>
</head>
```

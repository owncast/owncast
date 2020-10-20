# decrypter



## Table of Contents

<!-- START doctoc -->
<!-- END doctoc -->
## Installation

```sh
npm install --save aes-decrypter
```

Also available to install globally:

```sh
npm install --global aes-decrypter
```

The npm installation is preferred, but Bower works, too.

```sh
bower install  --save aes-decrypter
```

## Usage

To include decrypter on your website or npm application, use any of the following methods.
```js
var Decrypter = require('aes-decrypter').Decrypter;
var fs = require('fs');
var keyContent = fs.readFileSync('something.key');
var encryptedBytes = fs.readFileSync('somithing.txt');

// keyContent is a string of the aes-keys content
var keyContent = fs.readFileSync(keyFile);

var view = new DataView(keyContent.buffer);
var key.bytes = new Uint32Array([
  view.getUint32(0),
  view.getUint32(4),
  view.getUint32(8),
  view.getUint32(12)
]);

key.iv = new Uint32Array([
  0, 0, 0, 0
]);

var d = new Decrypter(
  encryptedBytes,
  key.bytes,
  key.iv,
  function(err, decryptedBytes) {
    // err always null
});
```

## [License](LICENSE)

Apache-2.0. Copyright (c) Brightcove, Inc.


'use strict';
const m3u8 = require('./export-m3u8s.js');

const args = require('minimist')(process.argv.slice(2), {
  boolean: ['watch', 'clean', 'build'],
  default: {
    build: true
  },
  alias: {
    b: 'build',
    c: 'clean',
    w: 'watch'
  }
});

if (args.w) {
  m3u8.watch();
} else if (args.c) {
  m3u8.clean();
} else if (args.b) {
  m3u8.build();
}

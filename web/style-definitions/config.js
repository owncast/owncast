const yaml = require('yaml');

module.exports = {
  hooks: {
    parsers: {
      // A custom parser will only run against filenames that match the pattern.
      // This pattern will match any file with the .yaml extension, which
      // allows you to mix different types of files in your token source.
      'yaml-parser': {
        pattern: /\.yaml$/,
        parser: ({ contents }) => yaml.parse(contents),
      },
    },
    fileHeaders: {
      // Intentionally omit Style Dictionary's default "Generated on <timestamp>"
      // line. The timestamp made every regeneration produce a diff, so the
      // committed outputs could never be reproduced byte-for-byte from the
      // source.
      myCustomHeader: () => [
        `Do not edit directly`,
        `This file is generated from the token sources under style-definitions/.`,
        ``,
        `How to edit these values:`,
        `Edit the corresponding token file under the style-definitions directory`,
        `in the Owncast web project.`,
      ],
    },
  },
  parsers: ['yaml-parser'],
  source: [`tokens/**/*.yaml`],
  platforms: {
    css: {
      transformGroup: 'css',
      buildPath: 'build/',
      files: [
        {
          destination: 'variables.css',
          format: 'css/variables',
          options: {
            fileHeader: 'myCustomHeader',
          },
        },
        {
          // The :root token block for the served plugin stylesheet. build.sh
          // appends the hand-authored element baseline (plugin-elements.css)
          // to this to form public/styles/plugin.css. Generated so the token
          // defaults stay in sync; the element rules only reference tokens.
          destination: 'plugin.css',
          format: 'css/variables',
          options: {
            fileHeader: 'myCustomHeader',
          },
        },
      ],
    },
  },
};

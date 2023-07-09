const yaml = require('yaml');
const StyleDictionary = require('style-dictionary');

StyleDictionary.registerFileHeader({
  name: 'myCustomHeader',
  fileHeader: defaultMessage => [
    ...defaultMessage,
    ``,
    `How to edit these values:`,
    `Edit the corresponding token file under the style-definitions directory`,
    `in the Owncast web project.`,
  ],
});

module.exports = {
  parsers: [
    {
      // A custom parser will only run against filenames that match the pattern
      // This pattern will match any file with the .yaml extension.
      // This allows you to mix different types of files in your token source
      pattern: /\.yaml$/,
      // the parse function takes a single argument, which is an object with
      // 2 attributes: contents which is a string of the file contents, and
      // filePath which is the path of the file.
      // The function is expected to return a plain object.
      parse: ({ contents, filePath }) => yaml.parse(contents),
    },
  ],
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
      ],
    },
    less: {
      transformGroup: 'less',
      buildPath: 'build/',
      files: [
        {
          destination: 'variables.less',
          format: 'less/variables',
          options: {
            fileHeader: 'myCustomHeader',
          },
        },
      ],
    },
    'ios-swift': {
      transforms: [
        'attribute/cti',
        'name/ti/camel',
        'color/ColorSwiftUI',
        'content/swift/literal',
        'asset/swift/literal',
        'size/swift/remToCGFloat',
        'font/swift/literal',
      ],
      buildPath: 'build/',
      files: [
        {
          format: 'ios-swift/class.swift',
          className: 'PlatformColor',
          destination: 'Colors.swift',
          filter: {
            attributes: {
              category: 'color',
            },
          },
        },
      ],
    },
  },
};

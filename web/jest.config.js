module.exports = {
  transform: {
    '^.+\\.(js|jsx|ts|tsx)$': [
      'babel-jest',
      {
        presets: [
          [
            '@babel/preset-env',
            {
              targets: {
                node: 'current',
              },
            },
          ],
          '@babel/preset-react',
          '@babel/preset-typescript',
        ],
        plugins: ['dynamic-import-node'],
      },
    ],
  },
  // @ant-design/icons v6's CJS build requires the ESM path
  // @ant-design/colors/es/generate directly; webpack transpiles it via
  // transpilePackages, jest needs the same exception here.
  transformIgnorePatterns: ['/node_modules/(?!@ant-design/colors)'],
  moduleNameMapper: {
    '\\.(css|less|scss)$': '<rootDir>/tests/__mocks__/styleMock.js',
  },
  testEnvironment: 'jsdom',
  testRegex: '/tests/.*\\.(test|spec)?\\.(ts|tsx)$',
  moduleFileExtensions: ['ts', 'tsx', 'js', 'jsx', 'json', 'node'],
  setupFilesAfterEnv: ['<rootDir>/tests/setup.ts'],
};

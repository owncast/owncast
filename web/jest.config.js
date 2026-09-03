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
          // Automatic runtime to match tsconfig's react-jsx: components no
          // longer import React just for JSX.
          ['@babel/preset-react', { runtime: 'automatic' }],
          '@babel/preset-typescript',
        ],
        plugins: ['dynamic-import-node'],
      },
    ],
  },
  // @ant-design/icons v6's CJS build requires the ESM path
  // @ant-design/colors/es/generate directly; webpack transpiles it via
  // transpilePackages, jest needs the same exception here.
  // FullCalendar publishes modern CommonJS that still needs Babel in Jest.
  transformIgnorePatterns: ['/node_modules/(?!@ant-design/colors|@fullcalendar)'],
  moduleNameMapper: {
    '\\.(css|less|scss)$': '<rootDir>/tests/__mocks__/styleMock.js',
  },
  testEnvironment: 'jsdom',
  testRegex: '/tests/.*\\.(test|spec)?\\.(ts|tsx)$',
  moduleFileExtensions: ['ts', 'tsx', 'js', 'jsx', 'json', 'node'],
  setupFilesAfterEnv: ['<rootDir>/tests/setup.ts'],
};

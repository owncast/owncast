module.exports = {
  transform: { '^.+\\.(js|jsx|ts|tsx)$': ['babel-jest', { presets: ['next/babel'] }] },
  testEnvironment: 'node',
  testRegex: '/tests/.*\\.(test|spec)?\\.(ts|tsx)$',
  moduleFileExtensions: ['ts', 'tsx', 'js', 'jsx', 'json', 'node'],
};

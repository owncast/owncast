module.exports = {
  env: {
    browser: true,
    es2021: true,
    jest: true,
  },
  extends: [
    'plugin:react/recommended',
    'airbnb',
    'prettier',
    'plugin:@next/next/recommended',
    'plugin:storybook/recommended',
  ],
  parser: '@typescript-eslint/parser',
  parserOptions: {
    ecmaFeatures: {
      jsx: true,
    },
    ecmaVersion: 12,
    sourceType: 'module',
  },
  plugins: ['react', 'prettier', '@typescript-eslint', 'import'],
  ignorePatterns: ['!./storybook/**'],
  rules: {
    'prettier/prettier': 'error',
    'react/destructuring-assignment': 'off',
    'react/prop-types': 0,
    'react/react-in-jsx-scope': 'off',
    'react/require-default-props': 'off',
    'react/no-is-mounted': 'off',
    'react/jsx-filename-extension': [
      1,
      {
        extensions: ['.js', '.jsx', '.tsx'],
      },
    ],
    'react/jsx-props-no-spreading': 'off',
    'react/jsx-no-bind': 'off',
    'react/function-component-definition': [
      'warn',
      {
        namedComponents: 'arrow-function',
        unnamedComponents: 'arrow-function',
      },
    ],
    '@next/next/no-img-element': 'off',
    'no-unused-vars': 'off',
    '@typescript-eslint/no-unused-vars': 'error',
    'no-console': 'off',
    'no-use-before-define': [0],
    '@typescript-eslint/no-use-before-define': [0],
    'no-shadow': 'off',
    '@typescript-eslint/no-shadow': ['error'],
    'no-restricted-exports': 'off',
    'no-plusplus': 'off',
    'react/jsx-no-target-blank': [
      1,
      {
        allowReferrer: false,
        enforceDynamicLinks: 'always',
      },
    ],
    'import/extensions': [
      'error',
      'ignorePackages',
      {
        js: 'never',
        jsx: 'never',
        ts: 'never',
        tsx: 'never',
      },
    ],
    'import/no-extraneous-dependencies': [
      'error',
      {
        devDependencies: ['**/*.stories.*', '**/.storybook/**/*.*'],
        peerDependencies: true,
      },
    ],
    'import/no-unused-modules': [1, { unusedExports: true }],
  },
  settings: {
    'import/resolver': {
      node: {
        extensions: ['.js', '.jsx', '.ts', '.tsx'],
      },
    },
  },
};

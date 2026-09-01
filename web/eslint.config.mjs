import { fixupConfigRules } from '@eslint/compat';
import { defineConfig, globalIgnores } from 'eslint/config';
import next from 'eslint-config-next';
import nextTypescript from 'eslint-config-next/typescript';
import prettier from 'eslint-config-prettier/flat';
import storybook from 'eslint-plugin-storybook';

export default defineConfig([
  ...fixupConfigRules(next),
  ...fixupConfigRules(nextTypescript),
  ...storybook.configs['flat/recommended'],
  prettier,
  {
    linterOptions: {
      reportUnusedDisableDirectives: 'off',
    },
    rules: {
      'react/destructuring-assignment': 'off',
      'react/prop-types': 'off',
      'react/react-in-jsx-scope': 'off',
      'react/require-default-props': 'off',
      'react/no-is-mounted': 'off',
      'react/jsx-filename-extension': 'off',
      'react/jsx-props-no-spreading': 'off',
      'react/jsx-no-bind': 'off',
      'react/function-component-definition': 'off',
      '@next/next/no-img-element': 'off',
      '@next/next/no-location-assign-relative-destination': 'off',
      'react/display-name': 'off',
      'react/jsx-key': 'off',
      'react-hooks/exhaustive-deps': 'off',
      'react-hooks/globals': 'off',
      'react-hooks/immutability': 'off',
      'react-hooks/purity': 'off',
      'react-hooks/rules-of-hooks': 'off',
      'react-hooks/set-state-in-effect': 'off',
      'no-unused-vars': 'off',
      '@typescript-eslint/no-unused-vars': 'error',
      '@typescript-eslint/no-empty-object-type': 'off',
      '@typescript-eslint/no-explicit-any': 'off',
      '@typescript-eslint/no-unsafe-function-type': 'off',
      '@typescript-eslint/no-wrapper-object-types': 'off',
      'no-console': 'off',
      'no-use-before-define': 'off',
      '@typescript-eslint/no-use-before-define': 'off',
      'no-shadow': 'off',
      '@typescript-eslint/no-shadow': 'error',
      'no-restricted-exports': 'off',
      'no-plusplus': 'off',
      'react/jsx-no-target-blank': [
        'warn',
        {
          allowReferrer: false,
          enforceDynamicLinks: 'always',
        },
      ],
      'import/no-extraneous-dependencies': [
        'error',
        {
          devDependencies: [
            '**/*.stories.*',
            '**/.storybook/**/*.*',
            '**/style-definitions/**/*.*',
            '**/build-scripts/**/*.*',
          ],
          peerDependencies: true,
        },
      ],
    },
  },
  {
    files: ['eslint.config.mjs'],
    rules: {
      'import/no-extraneous-dependencies': 'off',
    },
  },
  {
    files: ['components/**/*.{js,jsx,ts,tsx}'],
    rules: {
      'react/function-component-definition': [
        'warn',
        {
          namedComponents: 'arrow-function',
          unnamedComponents: 'arrow-function',
        },
      ],
      '@typescript-eslint/no-use-before-define': 'error',
      'import/prefer-default-export': 'off',
    },
  },
  {
    files: ['pages/**/*.{js,jsx,ts,tsx}'],
    rules: {
      'react/function-component-definition': 'off',
      'import/prefer-default-export': 'error',
    },
  },
  globalIgnores([
    'next-env.d.ts',
    'node_modules/**',
    'out/**',
    '../static/web/**',
    'public/styles/admin/config-public-details.css',
    '**/*.stories.*',
  ]),
]);

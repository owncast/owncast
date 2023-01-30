// ESLint rules specific to writing react components.

module.exports = {
  rules: {
    // Prefer arrow functions when defining functional components
    // This enables the `export const Foo: FC<FooProps> = ({ bar }) => ...` style we prefer.
    'react/function-component-definition': [
      'warn',
      {
        namedComponents: 'arrow-function',
        unnamedComponents: 'arrow-function',
      },
    ],

    // In functional components, mostly ensures Props are defined above components.
    '@typescript-eslint/no-use-before-define': 'error',

    // React components tend to use named exports.
    // Additionally, the `export default function` syntax cannot be auto-fixed by eslint when using ts.
    'import/prefer-default-export': 'off',
  },
};

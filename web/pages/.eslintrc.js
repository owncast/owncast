// ESLint rules specific to writing NextJS Pages.

module.exports = {
  rules: {
    // We don't care which syntax is used for NextPage definitions.
    'react/function-component-definition': 'off',

    // The default export is used by NextJS when rendering pages.
    'import/prefer-default-export': 'error',
  },
};

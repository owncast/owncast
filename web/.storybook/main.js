module.exports = {
  features: {
    previewMdx2: true,
  },

  // Admin components build their API URLs from NEXT_PUBLIC_API_HOST
  // (utils/apis.ts). Next.js loads it from .env.development/.env.production,
  // but Storybook does not, leaving it undefined and producing fetches to
  // "/undefinedapi/admin/...". Default it to "/" so admin page stories fetch
  // "/api/admin/..." and per-story msw handlers can intercept them.
  env: config => ({
    ...config,
    NEXT_PUBLIC_API_HOST: config.NEXT_PUBLIC_API_HOST || '/',
  }),

  stories: [
    '../.storybook/stories-category-doc-pages/**/*.@(mdx|stories.js)',
    '../stories/**/*.stories.@(js|jsx|ts|tsx)',
    '../components/**/*.stories.@(js|jsx|ts|tsx)',
    '../pages/**/*.stories.@(js|jsx|ts|tsx)',
  ],

  addons: ['@storybook/addon-links', '@storybook/addon-a11y', '@storybook/addon-docs'],

  webpackFinal: async (config, { configType }) => {
    // @see https://github.com/storybookjs/storybook/issues/9070
    const fileLoaderRule = config.module.rules.find(rule => rule.test && rule.test.test('.svg'));
    fileLoaderRule.exclude = /\.svg$/;

    // https://www.npmjs.com/package/@svgr/webpack
    config.module.rules.push({
      test: /\.svg$/,
      use: ['@svgr/webpack'],
    });

    return config;
  },

  framework: {
    name: '@storybook/nextjs',
    options: {},
  },

  staticDirs: ['../public', '../../static', './story-assets'],
  docs: {},

  typescript: {
    reactDocgen: 'react-docgen-typescript',
  },
};

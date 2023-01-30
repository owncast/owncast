module.exports = {
  core: {
    builder: 'webpack5',
  },
  features: {
    previewMdx2: true,
  },
  stories: [
    '../.storybook/stories-category-doc-pages/**/*.stories.mdx',
    '../stories/**/*.stories.@(js|jsx|ts|tsx)',
    '../components/**/*.stories.@(js|jsx|ts|tsx)',
    '../pages/**/*.stories.@(js|jsx|ts|tsx)',
  ],
  addons: [
    '@storybook/addon-links',
    '@storybook/addon-essentials',
    '@storybook/preset-scss',
    '@storybook/addon-postcss',
    '@storybook/addon-a11y',
    'storybook-addon-designs',
    'storybook-addon-fetch-mock',
  ],
  webpackFinal: async (config, { configType }) => {
    // @see https://github.com/storybookjs/storybook/issues/9070
    const fileLoaderRule = config.module.rules.find(rule => rule.test && rule.test.test('.svg'));
    fileLoaderRule.exclude = /\.svg$/;

    // https://www.npmjs.com/package/@svgr/webpack
    config.module.rules.push({
      test: /\.svg$/,
      use: ['@svgr/webpack'],
    });

    config.module.rules.push({
      test: /\.less$/,
      use: [
        require.resolve('style-loader'),
        require.resolve('css-loader'),
        {
          loader: require.resolve('less-loader'),
          options: {
            lessOptions: { javascriptEnabled: true },
          },
        },
      ],
    });

    return config;
  },
  framework: '@storybook/react',
  staticDirs: ['../public', '../../static', './story-assets'],
};

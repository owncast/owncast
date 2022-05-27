module.exports = {
  core: {
    builder: 'webpack5',
  },
  stories: [
    '../stories/Introduction.stories.mdx',
    '../stories/**/*.stories.mdx',
    '../stories/**/*.stories.@(js|jsx|ts|tsx)',
  ],
  addons: [
    '@storybook/addon-links',
    '@storybook/addon-essentials',
    '@storybook/addon-interactions',
    '@storybook/preset-scss',
    '@storybook/addon-postcss',
    '@storybook/addon-a11y',
    'storybook-addon-designs',
    'storybook-dark-mode',
    'addon-screen-reader',
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
};

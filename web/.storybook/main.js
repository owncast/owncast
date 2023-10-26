module.exports = {
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
    '@storybook/addon-a11y',
    'storybook-addon-fetch-mock',
    '@storybook/addon-mdx-gfm',
    '@storybook/addon-styling-webpack',
    {
      name: '@storybook/addon-styling-webpack',

      options: {
        rules: [
          {
            test: /\.css$/,
            sideEffects: true,
            use: [
              require.resolve('style-loader'),
              {
                loader: require.resolve('css-loader'),
                options: {
                  // Want to add more CSS Modules options? Read more here: https://github.com/webpack-contrib/css-loader#modules
                  modules: {
                    auto: true,
                  },
                },
              },
            ],
          },
          {
            test: /\.s[ac]ss$/,
            sideEffects: true,
            use: [
              require.resolve('style-loader'),
              {
                loader: require.resolve('css-loader'),
                options: {
                  // Want to add more CSS Modules options? Read more here: https://github.com/webpack-contrib/css-loader#modules
                  modules: {
                    auto: true,
                  },
                  importLoaders: 2,
                },
              },
              require.resolve('resolve-url-loader'),
              {
                loader: require.resolve('sass-loader'),
                options: {
                  // Want to add more Sass options? Read more here: https://webpack.js.org/loaders/sass-loader/#options
                  implementation: require.resolve('sass'),
                  sourceMap: true,
                  sassOptions: {},
                },
              },
            ],
          },
          {
            test: /\.less$/,
            sideEffects: true,
            use: [
              require.resolve('style-loader'),
              {
                loader: require.resolve('css-loader'),
                options: {
                  // Want to add more CSS Modules options? Read more here: https://github.com/webpack-contrib/css-loader#modules
                  modules: {
                    auto: true,
                  },
                  importLoaders: 1,
                },
              },
              {
                loader: require.resolve('less-loader'),
                options: {
                  // Want to add more Less options? Read more here: https://webpack.js.org/loaders/less-loader/#options
                  implementation: require.resolve('less'),
                  sourceMap: true,
                  lessOptions: {
                    javascriptEnabled: true,
                  },
                },
              },
            ],
          },
        ],
      },
    },
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

    return config;
  },

  framework: {
    name: '@storybook/nextjs',
    options: {},
  },

  staticDirs: ['../public', '../../static', './story-assets'],

  docs: {
    autodocs: false,
  },
};

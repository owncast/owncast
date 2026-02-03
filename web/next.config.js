const withBundleAnalyzer = require('@next/bundle-analyzer')({
  enabled: process.env.ANALYZE === 'true',
});
const { PHASE_DEVELOPMENT_SERVER } = require('next/constants');
const withPWA = require('next-pwa')({
  dest: 'public',
  runtimeCaching: [],
  register: true,
  skipWaiting: true,
  disableDevLogs: true,
  publicExcludes: ['!img/platformlogos/**/*', '!styles/admin/**/*'],
  buildExcludes: [/chunks\/pages\/admin.*/, '!**/admin/**/*'],
  sourcemap: process.env.NODE_ENV === 'development',
  disable: process.env.NODE_ENV === 'development',
});

async function rewrites() {
  return [
    {
      source: '/api/:path*',
      destination: 'http://localhost:8080/api/:path*',
    },
    {
      source: '/hls/:path*',
      destination: 'http://localhost:8080/hls/:path*',
    },
    {
      source: '/img/:path*',
      destination: 'http://localhost:8080/img/:path*',
    },
    {
      source: '/logo',
      destination: 'http://localhost:8080/logo',
    },
    {
      source: '/thumbnail.jpg',
      destination: 'http://localhost:8080/thumbnail.jpg',
    },
    {
      source: '/customjavascript',
      destination: 'http://localhost:8080/customjavascript',
    },
    {
      source: '/favicon.ico',
      destination: 'http://localhost:8080/favicon.ico',
    },
  ];
}

module.exports = async phase => {
  /**
   * @type {import('next').NextConfig}
   */
  let nextConfig = withBundleAnalyzer({
    productionBrowserSourceMaps: process.env.SOURCE_MAPS === 'true',
    trailingSlash: true,
    reactStrictMode: true,
    images: {
      unoptimized: true,
    },
    swcMinify: true,
    transpilePackages: [
      'antd',
      '@ant-design/cssinjs',
      '@ant-design/cssinjs-utils',
      '@ant-design/icons',
      '@rc-component/color-picker',
      '@rc-component/mutate-observer',
      '@rc-component/tour',
      '@rc-component/portal',
      '@rc-component/trigger',
      'rc-cascader',
      'rc-checkbox',
      'rc-collapse',
      'rc-dialog',
      'rc-drawer',
      'rc-dropdown',
      'rc-field-form',
      'rc-image',
      'rc-input',
      'rc-input-number',
      'rc-mentions',
      'rc-menu',
      'rc-motion',
      'rc-notification',
      'rc-overflow',
      'rc-pagination',
      'rc-picker',
      'rc-progress',
      'rc-rate',
      'rc-resize-observer',
      'rc-segmented',
      'rc-select',
      'rc-slider',
      'rc-steps',
      'rc-switch',
      'rc-table',
      'rc-tabs',
      'rc-textarea',
      'rc-tooltip',
      'rc-tree',
      'rc-tree-select',
      'rc-upload',
      'rc-util',
      'rc-virtual-list',
    ],
    webpack(config) {
      config.module.rules.push({
        test: /\.svg$/i,
        issuer: /\.[jt]sx?$/,
        use: ['@svgr/webpack'],
      });

      return config;
    },
    pageExtensions: ['tsx'],
  });

  if (phase === PHASE_DEVELOPMENT_SERVER) {
    nextConfig = {
      ...nextConfig,
      rewrites,
    };
  } else {
    nextConfig = {
      ...nextConfig,
      output: 'export',
    };
  }
  return withPWA(nextConfig);
};

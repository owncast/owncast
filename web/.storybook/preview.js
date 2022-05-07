import '../styles/variables.css';
import '../styles/global.less';
import '../styles/theme.less';
import '../stories/preview.scss';
import { themes } from '@storybook/theming';

export const parameters = {
  actions: { argTypesRegex: '^on[A-Z].*' },
  controls: {
    matchers: {
      color: /(background|color)$/i,
      date: /Date$/,
    },
  },
  darkMode: {
    current: 'dark',
    // Override the default dark theme
    dark: {
      ...themes.dark,
      appBg: '#171523',
      brandImage: 'https://owncast.online/images/logo.svg',
      brandTitle: 'Owncast',
      brandUrl: 'https://owncast.online',
    },
    // Override the default light theme
    light: { ...themes.normal },
  },
};

import '../styles/variables.scss';
import '../styles/global.less';
import '../styles/theme.less';
import '../stories/preview.scss';

export const parameters = {
  actions: { argTypesRegex: '^on[A-Z].*' },
  controls: {
    matchers: {
      color: /(background|color)$/i,
      date: /Date$/,
    },
  },
};

import '../styles/variables.css';
import '../styles/global.less';
import '../styles/theme.less';
import './preview.scss';
import { themes } from '@storybook/theming';
import { DocsContainer } from './storybook-theme';
import { INITIAL_VIEWPORTS } from '@storybook/addon-viewport';
import _ from 'lodash';

/**
 * Takes an entry of a viewport (from Object.entries()) and converts it
 * into two entries, one for landscape and one for portrait.
 *
 * @template {string} Key
 *
 * @param {[Key, import('@storybook/addon-viewport/dist/ts3.9/models').Viewport]} entry
 * @returns {Array<[`${Key}${'Portrait' | 'Landscape'}`, import('@storybook/addon-viewport/dist/ts3.9/models').Viewport]>}
 */
const convertToLandscapeAndPortraitEntries = ([objectKey, viewport]) => {
  const pixelStringToNumber = str => parseInt(str.split('px')[0]);
  const dimensions = [viewport.styles.width, viewport.styles.height].map(pixelStringToNumber);
  const minDimension = Math.min(...dimensions);
  const maxDimension = Math.max(...dimensions);

  return [
    [
      `${objectKey}Portrait`,
      {
        ...viewport,
        name: viewport.name + ' (Portrait)',
        styles: {
          ...viewport.styles,
          height: maxDimension + 'px',
          width: minDimension + 'px',
        },
      },
    ],
    [
      `${objectKey}Landscape`,
      {
        ...viewport,
        name: viewport.name + ' (Landscape)',
        styles: {
          ...viewport.styles,
          height: minDimension + 'px',
          width: maxDimension + 'px',
        },
      },
    ],
  ];
};

/**
 * Takes an object and a function f and returns a new object.
 * f takes the original object's entries (key-value-pairs
 * from Object.entries) and returns a list of new entries
 * (also key-value-pairs). These new entries then form the
 * result.
 * @template {string | number} OriginalKey
 * @template {string | number} NewKey
 * @template OriginalValue
 * @template OriginalValue
 *
 * @param {Record<OriginalKey, OriginalValue>} obj
 * @param {(entry: [OriginalKey, OriginalValue], index: number, all: Array<[OriginalKey, OriginalValue]>) => Array<[NewKey, NewValue]>} f
 * @returns {Record<NewKey, NevValue>}
 */
const flatMapObject = (obj, f) => Object.fromEntries(Object.entries(obj).flatMap(f));

export const parameters = {
  fetchMock: {
    mocks: [],
  },
  actions: { argTypesRegex: '^on[A-Z].*' },
  docs: {
    container: DocsContainer,
  },
  actions: { argTypesRegex: '^on[A-Z].*' },
  viewMode: 'docs',
  controls: {
    matchers: {
      color: /(background|color)$/i,
      date: /Date$/,
    },
    viewMode: 'docs',
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
      appContentBg: '#171523',
    },
    // Override the default light theme
    light: { ...themes.normal },
  },
  viewport: {
    // Take a bunch of viewports from the storybook addon and convert them
    // to portrait + landscape. Keys are appended with 'Landscape' or 'Portrait'.
    viewports: flatMapObject(INITIAL_VIEWPORTS, convertToLandscapeAndPortraitEntries),
  },
};

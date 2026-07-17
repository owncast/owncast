import '../styles/variables.css';
// Pre-extracted Ant Design styles: required in zeroRuntime mode (see
// build-scripts/extract-antd-styles.js).
import '../styles/antd.css';
import './preview.scss';
import { themes } from 'storybook/theming';
import { DocsContainer } from './storybook-theme';
import { INITIAL_VIEWPORTS } from 'storybook/viewport';
import _ from 'lodash';
import { initialize, mswLoader } from 'msw-storybook-addon';
import React from 'react';
import { Provider, useSetAtom } from 'jotai';
import { AntdProvider } from '../components/theme/AntdProvider';
import { Theme } from '../components/theme/Theme';
import { clientConfigStateAtom } from '../components/stores/ClientConfigStore';
import { makeEmptyClientConfig } from '../interfaces/client-config.model';
import THEMES from '../stories/themePresets';

/**
 * Takes an entry of a viewport (from Object.entries()) and converts it
 * into two entries, one for landscape and one for portrait.
 *
 * @template {string} Key
 *
 * @param {[Key, import('storybook/viewport/dist/ts3.9/models').Viewport]} entry
 * @returns {Array<[`${Key}${'Portrait' | 'Landscape'}`, import('storybook/viewport/dist/ts3.9/models').Viewport]>}
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

// Initialize MSW
initialize();

export const parameters = {
  // actions: { argTypesRegex: '^on[A-Z].*' },
  docs: {
    container: DocsContainer,
  },
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
    // Keep the base viewport keys (stories reference 'tablet', 'mobile1', ...)
    // and add the rotated variants; replacing `options` wholesale would drop
    // the base names and every viewport-pinned story would silently render
    // at the default desktop width (Chromatic captures included).
    options: {
      ...INITIAL_VIEWPORTS,
      ...flatMapObject(INITIAL_VIEWPORTS, convertToLandscapeAndPortraitEntries),
    },
  },
};

export const loaders = [mswLoader];

// Toolbar switcher for the appearance-theme presets (stories/themePresets).
// Applies the selected preset through the two real theming paths — the
// clientConfig atom (AntdProvider maps it onto antd design tokens) and the
// Theme component (emits the --theme-* CSS variables) — so any story can be
// checked against custom themes, not just the Theme playground.
export const globalTypes = {
  owncastTheme: {
    description: 'Owncast appearance theme preset',
    toolbar: {
      title: 'Theme',
      icon: 'paintbrush',
      items: Object.entries(THEMES).map(([value, t]) => ({ value, title: t.label })),
      dynamicTitle: true,
    },
  },
};

export const initialGlobals = { owncastTheme: 'default' };

const ApplyStorybookTheme = ({ theme, children }) => {
  const setClientConfig = useSetAtom(clientConfigStateAtom);
  const preset = theme !== 'default' && THEMES[theme];
  React.useEffect(() => {
    if (preset) {
      setClientConfig({
        ...makeEmptyClientConfig(),
        appearanceVariables: preset.variables,
      });
    }
  }, [preset, setClientConfig]);
  if (!preset) {
    return children;
  }
  return (
    <>
      <Theme />
      {children}
    </>
  );
};

// Mirror the app-wide providers from pages/_app.tsx: antd components rely on
// AntdProvider for the Owncast theme tokens. Without this decorator those
// stories would render unthemed. AntdProvider reads jotai atoms, so a
// Provider wraps it (stories with their own Provider simply nest; that
// is supported). The Provider is keyed by the selected toolbar theme so
// switching themes starts from a fresh store instead of layering presets.
export const decorators = [
  (Story, context) => (
    <Provider key={context.globals.owncastTheme || 'default'}>
      <AntdProvider>
        <ApplyStorybookTheme theme={context.globals.owncastTheme}>
          <Story />
        </ApplyStorybookTheme>
      </AntdProvider>
    </Provider>
  ),
];

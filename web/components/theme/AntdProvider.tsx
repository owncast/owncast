import { FC, ReactNode, useEffect, useMemo } from 'react';
import { ConfigProvider } from 'antd';
import type { ThemeConfig } from 'antd';
import { useAtomValue } from 'jotai';
import { clientConfigStateAtom } from '../stores/ClientConfigStore';
import antdDefaultTheme from './antd-default-theme.json';

/*
Theme bridge for Ant Design.

There is exactly ONE instance of this provider, mounted in pages/_app.tsx
(and mirrored as a Storybook decorator in .storybook/preview.js). It maps
the Owncast theme onto antd design tokens: static defaults below mirror
web/style-definitions default theme values; runtime admin customizations
arrive through clientConfig.appearanceVariables (the same values Theme.tsx
sets as CSS variables) and override the defaults.

Component style RULES are pre-extracted at build time into styles/antd.css
(zeroRuntime mode); the runtime injects only the CSS-variable values for the
tokens, so theme changes apply live with no runtime style generation.
*/

// Defaults mirroring web/style-definitions/tokens (default theme). The JSON
// is shared with build-scripts/extract-antd-styles.js, which bakes these
// values into the pre-extracted stylesheet as the pre-hydration fallback.
const defaultTheme = antdDefaultTheme as ThemeConfig;

// Map the admin-customizable Owncast appearance variables (documented at
// https://owncast.online and set as CSS variables by Theme.tsx) onto v6
// design tokens. Only variables relevant to migrated components are mapped;
// extend this as more components move to v6.
type TokenOverrides = NonNullable<ThemeConfig['token']>;
type ComponentOverrides = NonNullable<ThemeConfig['components']>;
type ModalOverrides = NonNullable<ComponentOverrides['Modal']>;
type TabsOverrides = NonNullable<ComponentOverrides['Tabs']>;
type DropdownOverrides = NonNullable<ComponentOverrides['Dropdown']>;
type MenuOverrides = NonNullable<ComponentOverrides['Menu']>;

function buildTheme(appearanceVariables: Record<string, string>): ThemeConfig {
  const vars = appearanceVariables || {};
  const token: TokenOverrides = { ...defaultTheme.token };
  const modal: ModalOverrides = { ...defaultTheme.components?.Modal };
  const tabs: TabsOverrides = { ...defaultTheme.components?.Tabs };
  const dropdown: DropdownOverrides = { ...defaultTheme.components?.Dropdown };
  // The admin sider menu runs antd's dark theme; these tokens come straight
  // from the default theme file and are not admin-customizable (the admin
  // never receives appearanceVariables).
  const menu: MenuOverrides = { ...defaultTheme.components?.Menu };

  if (vars['theme-color-action']) {
    token.colorPrimary = vars['theme-color-action'];
    token.colorLink = vars['theme-color-action'];
    tabs.itemActiveColor = vars['theme-color-action'];
    tabs.itemSelectedColor = vars['theme-color-action'];
    tabs.inkBarColor = vars['theme-color-action'];
  }
  if (vars['theme-color-action-hover']) {
    token.colorLinkHover = vars['theme-color-action-hover'];
    token.colorLinkActive = vars['theme-color-action-hover'];
    tabs.itemHoverColor = vars['theme-color-action-hover'];
  }
  if (vars['theme-color-background-main']) {
    token.colorBgLayout = vars['theme-color-background-main'];
  }
  if (vars['theme-color-components-text-on-light']) {
    token.colorText = vars['theme-color-components-text-on-light'];
  }
  if (vars['theme-rounded-corners']) {
    const radius = parseInt(vars['theme-rounded-corners'], 10);
    if (!Number.isNaN(radius)) {
      token.borderRadius = radius;
    }
  }
  if (vars['theme-color-components-form-field-background']) {
    token.colorBgContainer = vars['theme-color-components-form-field-background'];
  }

  if (vars['theme-color-components-modal-header-background']) {
    modal.headerBg = vars['theme-color-components-modal-header-background'];
  }
  if (vars['theme-color-components-modal-header-text']) {
    modal.titleColor = vars['theme-color-components-modal-header-text'];
  }
  if (vars['theme-color-components-modal-content-background']) {
    modal.contentBg = vars['theme-color-components-modal-content-background'];
  }

  if (vars['theme-color-components-menu-background']) {
    dropdown.colorBgElevated = vars['theme-color-components-menu-background'];
  }
  if (vars['theme-color-components-menu-item-text']) {
    dropdown.colorText = vars['theme-color-components-menu-item-text'];
  }
  if (vars['theme-color-components-menu-item-hover-bg']) {
    dropdown.controlItemBgHover = vars['theme-color-components-menu-item-hover-bg'];
  }
  if (vars['theme-color-components-menu-item-focus-bg']) {
    // v4 themed the tab bar divider with this variable (ant-overrides.scss
    // .ant-tabs-nav::before); the v6 equivalent is the Tabs border token.
    tabs.colorBorderSecondary = vars['theme-color-components-menu-item-focus-bg'];
  }

  // Whether any admin appearance customization is in play. The pre-extracted
  // antd.css bakes the DEFAULT token values under the stable cssVar key so
  // the static export is fully styled at first paint. Those static value
  // blocks have the same specificity as the runtime-injected ones, so a
  // customized theme must live under a DIFFERENT scope key: otherwise the
  // static defaults can win the cascade and admin theming silently reverts
  // to default (a stuck-purple primary button).
  const hasCustomizations = Object.keys(vars).length > 0;

  return {
    token,
    components: { Modal: modal, Tabs: tabs, Dropdown: dropdown, Menu: menu },
    // Stable css-var scope key. Without an explicit key antd derives it from
    // React useId, which differs between the extraction script, the Next
    // build render, and the browser, so the class baked into the static
    // export never matches the selectors in the pre-extracted antd.css and
    // var()-sized components (like the header Avatar) render unconstrained
    // until hydration injects the runtime values.
    cssVar: hasCustomizations
      ? { ...defaultTheme.cssVar, key: `${defaultTheme.cssVar.key}-custom` }
      : defaultTheme.cssVar,
    // Component style RULES come from the pre-extracted styles/antd.css
    // (see build-scripts/extract-antd-styles.js); the runtime only injects
    // the CSS-variable VALUES derived from the tokens above, which is what
    // keeps dynamic theming working with no runtime rule generation and no
    // FOUC on statically exported pages.
    zeroRuntime: true,
    // Must match the extraction script: without the hash class, selectors in
    // the static CSS line up with the classes components render.
    hashed: false,
  };
}

// Static APIs (message, notification, Modal.confirm) render into their own
// detached React root, outside the provider tree above. Their holder is
// configured once at module scope and reads the current theme through this
// variable, which the provider keeps up to date. antd re-invokes
// holderRender on every static call (message and notification resync their
// holder per call, Modal.confirm renders fresh per call), so statics opened
// after an appearance change pick up the latest tokens.
let staticsTheme: ThemeConfig = buildTheme({});
ConfigProvider.config({
  holderRender: node => <ConfigProvider theme={staticsTheme}>{node}</ConfigProvider>,
});

export type AntdProviderProps = {
  children: ReactNode;
};

export const AntdProvider: FC<AntdProviderProps> = ({ children }) => {
  const clientConfig = useAtomValue(clientConfigStateAtom);
  const { appearanceVariables } = clientConfig;

  // Stringify so the memo reacts to content changes, not object identity.
  const appearanceKey = JSON.stringify(appearanceVariables);
  const theme = useMemo(() => buildTheme(appearanceVariables), [appearanceKey]);

  useEffect(() => {
    staticsTheme = theme;
  }, [theme]);

  return <ConfigProvider theme={theme}>{children}</ConfigProvider>;
};

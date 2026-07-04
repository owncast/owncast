import { FC, ReactNode, useMemo } from 'react';
import { ConfigProvider } from 'antd6';
import type { ThemeConfig } from 'antd6';
import { useRecoilValue } from 'recoil';
import { ClientConfig } from '../../interfaces/client-config.model';
import { clientConfigStateAtom } from '../stores/ClientConfigStore';
import antd6DefaultTheme from './antd6-default-theme.json';

/*
Theme bridge for Ant Design v6 components during the incremental v4 -> v6
migration.

antd v6 is installed under the npm alias "antd6" so it can coexist with v4.
There is exactly ONE instance of this provider, mounted in pages/_app.tsx
(and mirrored as a Storybook decorator in .storybook/preview.js); migrated
components just import from "antd6". It does two things:

1. Scopes v6 under the "ant6" class prefix so the global v4 Less styles
	 (.ant-*) cannot bleed into v6 components, and vice versa. When the
	 migration is complete and v4 is removed, drop the prefix.

2. Maps the Owncast theme onto v6 design tokens. Static defaults below
	 mirror web/style-definitions default theme values; runtime admin
	 customizations arrive through clientConfig.appearanceVariables (the same
	 values Theme.tsx sets as CSS variables) and override the defaults.

Component style RULES are pre-extracted at build time into styles/antd6.css
(zeroRuntime mode); the runtime injects only the CSS-variable values for the
tokens, so theme changes apply live with no runtime style generation.
*/

export const ANTD6_PREFIX = 'ant6';

// Defaults mirroring web/style-definitions/tokens (default theme). The JSON
// is shared with build-scripts/extract-antd6-styles.js, which bakes these
// values into the pre-extracted stylesheet as the pre-hydration fallback.
const defaultTheme = antd6DefaultTheme as ThemeConfig;

// Map the admin-customizable Owncast appearance variables (documented at
// https://owncast.online and set as CSS variables by Theme.tsx) onto v6
// design tokens. Only variables relevant to migrated components are mapped;
// extend this as more components move to v6.
type TokenOverrides = NonNullable<ThemeConfig['token']>;
type ComponentOverrides = NonNullable<ThemeConfig['components']>;
type ModalOverrides = NonNullable<ComponentOverrides['Modal']>;
type TabsOverrides = NonNullable<ComponentOverrides['Tabs']>;

function buildTheme(appearanceVariables: Record<string, string>): ThemeConfig {
  const vars = appearanceVariables || {};
  const token: TokenOverrides = { ...defaultTheme.token };
  const modal: ModalOverrides = { ...defaultTheme.components?.Modal };
  const tabs: TabsOverrides = { ...defaultTheme.components?.Tabs };

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

  return {
    token,
    components: { Modal: modal, Tabs: tabs },
    // Component style RULES come from the pre-extracted styles/antd6.css
    // (see build-scripts/extract-antd6-styles.js); the runtime only injects
    // the CSS-variable VALUES derived from the tokens above, which is what
    // keeps dynamic theming working with no runtime rule generation and no
    // FOUC on statically exported pages.
    zeroRuntime: true,
    // Must match the extraction script: without the hash class, selectors in
    // the static CSS line up with the classes components render.
    hashed: false,
  };
}

export type Antd6ProviderProps = {
  children: ReactNode;
};

export const Antd6Provider: FC<Antd6ProviderProps> = ({ children }) => {
  const clientConfig = useRecoilValue<ClientConfig>(clientConfigStateAtom);
  const { appearanceVariables } = clientConfig;

  // Stringify so the memo reacts to content changes, not object identity.
  const appearanceKey = JSON.stringify(appearanceVariables);
  const theme = useMemo(() => buildTheme(appearanceVariables), [appearanceKey]);

  return (
    <ConfigProvider prefixCls={ANTD6_PREFIX} iconPrefixCls={`${ANTD6_PREFIX}icon`} theme={theme}>
      {children}
    </ConfigProvider>
  );
};

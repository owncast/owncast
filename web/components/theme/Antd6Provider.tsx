import { FC, ReactNode, useMemo } from 'react';
import { ConfigProvider } from 'antd6';
import type { ThemeConfig } from 'antd6';
import { useRecoilValue } from 'recoil';
import { ClientConfig } from '../../interfaces/client-config.model';
import { clientConfigStateAtom } from '../stores/ClientConfigStore';

/*
Theme bridge for Ant Design v6 components during the incremental v4 -> v6
migration.

antd v6 is installed under the npm alias "antd6" so it can coexist with v4.
Components migrated to v6 must be wrapped in this provider. It does two
things:

1. Scopes v6 under the "ant6" class prefix so the global v4 Less styles
	 (.ant-*) cannot bleed into v6 components, and vice versa. When the
	 migration is complete and v4 is removed, drop the prefix.

2. Maps the Owncast theme onto v6 design tokens. Static defaults below
	 mirror web/style-definitions default theme values; runtime admin
	 customizations arrive through clientConfig.appearanceVariables (the same
	 values Theme.tsx sets as CSS variables) and override the defaults.

v6 runs in CSS-variables mode by default, so token changes apply live
without style recomputation.
*/

export const ANTD6_PREFIX = 'ant6';

// Defaults mirroring web/style-definitions/tokens (default theme). If the
// style-dictionary defaults change, these need to change with them.
const defaultTheme: ThemeConfig = {
  token: {
    colorPrimary: '#6544e9',
    colorLink: '#6544e9',
    colorLinkHover: '#7a5cf3',
    colorLinkActive: '#7a5cf3',
    borderRadius: 9,
    colorBgContainer: '#ffffff',
    colorBgLayout: '#e2e8f0',
    colorBgElevated: '#ffffff',
    colorText: '#12161d',
    colorTextSecondary: '#5d5f72',
    colorError: '#ff4b39',
    colorWarning: '#ffc655',
    fontFamily:
      "Inter, system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, 'Noto Sans', sans-serif",
  },
  components: {
    Modal: {
      headerBg: '#2d3748',
      titleColor: '#e2e8f0',
      contentBg: '#e2e8f0',
      // v4-style modal chrome: distinct header bar with its own background.
      wireframe: true,
    },
  },
};

// Map the admin-customizable Owncast appearance variables (documented at
// https://owncast.online and set as CSS variables by Theme.tsx) onto v6
// design tokens. Only variables relevant to migrated components are mapped;
// extend this as more components move to v6.
type TokenOverrides = NonNullable<ThemeConfig['token']>;
type ModalOverrides = NonNullable<NonNullable<ThemeConfig['components']>['Modal']>;

function buildTheme(appearanceVariables: Record<string, string>): ThemeConfig {
  const vars = appearanceVariables || {};
  const token: TokenOverrides = { ...defaultTheme.token };
  const modal: ModalOverrides = { ...defaultTheme.components?.Modal };

  if (vars['theme-color-action']) {
    token.colorPrimary = vars['theme-color-action'];
    token.colorLink = vars['theme-color-action'];
  }
  if (vars['theme-color-action-hover']) {
    token.colorLinkHover = vars['theme-color-action-hover'];
    token.colorLinkActive = vars['theme-color-action-hover'];
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

  return { token, components: { Modal: modal } };
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

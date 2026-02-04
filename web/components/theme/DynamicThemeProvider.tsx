import { FC, ReactNode, useMemo } from 'react';
import { ConfigProvider } from 'antd';
import { useRecoilValue } from 'recoil';
import { ClientConfig } from '../../interfaces/client-config.model';
import { clientConfigStateAtom } from '../stores/ClientConfigStore';
import { antdTheme } from './antdThemeConfig';

interface DynamicThemeProviderProps {
  children: ReactNode;
}

/**
 * Dynamic theme provider for the main (viewer-facing) pages.
 * Reads appearance variables from the client config and applies them
 * to Ant Design components via ConfigProvider.
 *
 * This is the proper Ant Design v5 way to handle dynamic theming,
 * rather than using CSS !important overrides.
 */
export const DynamicThemeProvider: FC<DynamicThemeProviderProps> = ({ children }) => {
  const clientConfig = useRecoilValue<ClientConfig>(clientConfigStateAtom);
  const { appearanceVariables } = clientConfig;

  // Stringify to ensure React detects object content changes
  const appearanceKey = JSON.stringify(appearanceVariables);

  const dynamicTheme = useMemo(() => {
    // Start with the base theme
    const themeConfig = { ...antdTheme };

    // If no custom appearance variables, use the default theme
    if (!appearanceVariables || Object.keys(appearanceVariables).length === 0) {
      return themeConfig;
    }

    // Map Owncast appearance variables to Ant Design tokens
    const tokens: Record<string, string | number> = {};
    const componentTokens: Record<string, Record<string, unknown>> = {};

    // Primary/Action colors
    if (appearanceVariables['theme-color-action']) {
      tokens.colorPrimary = appearanceVariables['theme-color-action'];
      tokens.colorLink = appearanceVariables['theme-color-action'];
    }
    if (appearanceVariables['theme-color-action-hover']) {
      tokens.colorLinkHover = appearanceVariables['theme-color-action-hover'];
      tokens.colorLinkActive = appearanceVariables['theme-color-action-hover'];
    }

    // Background colors
    if (appearanceVariables['theme-color-background-main']) {
      tokens.colorBgLayout = appearanceVariables['theme-color-background-main'];
    }
    if (appearanceVariables['theme-color-palette-4']) {
      tokens.colorBgContainer = appearanceVariables['theme-color-palette-4'];
      tokens.colorBgElevated = appearanceVariables['theme-color-palette-4'];
    }

    // Text colors
    if (appearanceVariables['theme-color-components-text-on-light']) {
      tokens.colorText = appearanceVariables['theme-color-components-text-on-light'];
      tokens.colorTextSecondary = appearanceVariables['theme-color-components-text-on-light'];
      tokens.colorTextHeading = appearanceVariables['theme-color-components-text-on-light'];
    }

    // Border radius
    if (appearanceVariables['theme-rounded-corners']) {
      const radius = parseInt(appearanceVariables['theme-rounded-corners'], 10);
      if (!Number.isNaN(radius)) {
        tokens.borderRadius = radius;
      }
    }

    // Button component tokens
    if (
      appearanceVariables['theme-color-components-primary-button-border'] ||
      appearanceVariables['theme-color-components-primary-button-text']
    ) {
      componentTokens.Button = {
        ...antdTheme.components?.Button,
      };
      if (appearanceVariables['theme-color-components-primary-button-border']) {
        componentTokens.Button.colorPrimaryBorder =
          appearanceVariables['theme-color-components-primary-button-border'];
      }
      if (appearanceVariables['theme-color-components-primary-button-text']) {
        componentTokens.Button.primaryColor =
          appearanceVariables['theme-color-components-primary-button-text'];
      }
    }

    // Modal component tokens
    if (
      appearanceVariables['theme-color-components-modal-header-background'] ||
      appearanceVariables['theme-color-components-modal-content-background']
    ) {
      componentTokens.Modal = {
        ...antdTheme.components?.Modal,
      };
      if (appearanceVariables['theme-color-components-modal-header-background']) {
        componentTokens.Modal.headerBg =
          appearanceVariables['theme-color-components-modal-header-background'];
      }
      if (appearanceVariables['theme-color-components-modal-header-text']) {
        componentTokens.Modal.titleColor =
          appearanceVariables['theme-color-components-modal-header-text'];
      }
      if (appearanceVariables['theme-color-components-modal-content-background']) {
        componentTokens.Modal.contentBg =
          appearanceVariables['theme-color-components-modal-content-background'];
      }
    }

    // Input component tokens
    if (appearanceVariables['theme-color-components-form-field-background']) {
      componentTokens.Input = {
        ...antdTheme.components?.Input,
        colorBgContainer: appearanceVariables['theme-color-components-form-field-background'],
      };
    }

    // Tabs component tokens
    if (appearanceVariables['theme-color-action']) {
      componentTokens.Tabs = {
        ...antdTheme.components?.Tabs,
        itemActiveColor: appearanceVariables['theme-color-action'],
        itemSelectedColor: appearanceVariables['theme-color-action'],
        inkBarColor: appearanceVariables['theme-color-action'],
      };
      if (appearanceVariables['theme-color-action-hover']) {
        componentTokens.Tabs.itemHoverColor = appearanceVariables['theme-color-action-hover'];
      }
    }

    return {
      ...themeConfig,
      token: {
        ...themeConfig.token,
        ...tokens,
      },
      components: {
        ...themeConfig.components,
        ...componentTokens,
      },
    };
  }, [appearanceKey]);

  return <ConfigProvider theme={dynamicTheme}>{children}</ConfigProvider>;
};

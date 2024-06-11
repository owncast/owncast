import React, { FC, ReactNode, useEffect, useState } from 'react';
import { ConfigProvider } from 'antd';

export type AntdThemeProps = {
  children: ReactNode;
};

const readCSSVar = (varName: string) =>
  getComputedStyle(document.documentElement).getPropertyValue(varName);

export const AntdTheme: FC<AntdThemeProps> = ({ children }) => {
  const [customToken, setCustomToken] = useState({});

  const updateToken = () => {
    // for global properties see https://ant.design/docs/react/customize-theme#theme
    // for per-component see https://ant.design/docs/react/migrate-less-variables
    setCustomToken({
      colorLink: readCSSVar('--link-color'),
      colorLinkHover: readCSSVar('--link-hover-color'),
      Modal: {
        headerBg: readCSSVar('--modal-header-bg'),
        contentBg: readCSSVar('--modal-content-bg'),
      },
      Alert: {
        colorErrorBg: readCSSVar('--alert-error-bg-color'),
        colorErrorBorder: readCSSVar('--alert-error-border-color'),
      },
      colorBgElevated: readCSSVar('--popover-background'),
      Tag: {
        defaultBg: readCSSVar('--tag-default-color'),
      },
      borderRadius: parseInt(readCSSVar('--border-radius-base'), 10),
      // background-color-light needs mapping
      colorIcon: readCSSVar('--modal-close-color'),

      colorPrimary: readCSSVar('--primary-color'),
      colorPrimaryHover: readCSSVar('--primary-color-hover'),
      colorPrimaryActive: readCSSVar('--primary-color-active'),
      // primary-$n needs mapping
      colorBgBase: readCSSVar('--component-background'),
      // body-background needs mapping
    });
  };

  const updateTokenOnMutation = (_, observer: MutationObserver) => {
    updateToken();
    observer.disconnect();
  };

  // removes the initial cache style, see pages/_document.tsx
  const purgeInitCache = () => {
    const cache = document.getElementById('antd-init-cache');
    if (cache !== null) {
      cache.remove();
    }
  };

  const observeCustomThemeChanges = () => {
    const customStyles = document.getElementById('custom-color-styles');
    if (customStyles === null) {
      return;
    }
    console.log(customStyles.textContent);
    const mutObserver = new MutationObserver(updateTokenOnMutation);
    // check if color styles have been updated by ./Theme
    mutObserver.observe(customStyles, {
      characterData: true,
      subtree: true,
    });
  };

  useEffect(() => {
    updateToken();
    purgeInitCache();
    observeCustomThemeChanges();
  }, []);

  return (
    <ConfigProvider theme={{ cssVar: true, hashed: false, token: customToken }}>
      {children}
    </ConfigProvider>
  );
};

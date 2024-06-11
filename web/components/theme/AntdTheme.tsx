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
    setCustomToken({
      colorPrimary: readCSSVar('--primary-color'),
      linkColor: readCSSVar('--link-color'),
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

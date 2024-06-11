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
    });
    console.log(customToken);
  };

  // removes the initial cache, see pages/_document.tsx
  const purgeInitCache = () => {
    const cache = document.getElementById('antd-init-cache');
    if (cache !== null) {
      cache.remove();
    }
  };

  useEffect(() => {
    updateToken();
    purgeInitCache();
  }, []);

  return (
    <ConfigProvider theme={{ cssVar: true, hashed: false, token: customToken }}>
      {children}
    </ConfigProvider>
  );
};

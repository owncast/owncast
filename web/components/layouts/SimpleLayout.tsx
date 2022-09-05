import { AppProps } from 'next/app';
import { FC } from 'react';

export const SimpleLayout: FC<AppProps> = ({ Component, pageProps }) => (
  <div>
    <Component {...pageProps} />
  </div>
);

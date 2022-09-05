import { AppProps } from 'next/app';
import { FC } from 'react';

const SimpleLayout: FC<AppProps> = ({ Component, pageProps }) => (
  <div>
    <Component {...pageProps} />
  </div>
);

export default SimpleLayout;

import { AppProps } from 'next/app';

const SimpleLayout = ({ Component, pageProps }: AppProps) => (
  <div>
    <Component {...pageProps} />
  </div>
);

export default SimpleLayout;

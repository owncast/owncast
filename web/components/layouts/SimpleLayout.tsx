import { AppProps } from 'next/app';

function SimpleLayout({ Component, pageProps }: AppProps) {
  return (
    <div>
      <Component {...pageProps} />
    </div>
  );
}

export default SimpleLayout;

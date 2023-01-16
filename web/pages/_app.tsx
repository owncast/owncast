// All these imports are almost exclusively for the Admin.
// We should not be loading them for the main frontend UI.

// order matters!
import '../styles/variables.css';
import '../styles/global.less';
import '../styles/globals.scss';
import '../styles/ant-overrides.scss';

// TODO: Move this videojs sass to the player component.
import '../components/video/VideoJS/VideoJS.scss';

import { AppProps } from 'next/app';
import { ReactElement, ReactNode } from 'react';
import { NextPage } from 'next';
import { RecoilRoot } from 'recoil';

export type NextPageWithLayout<P = {}, IP = P> = NextPage<P, IP> & {
  getLayout?: (page: ReactElement) => ReactNode;
};

type AppPropsWithLayout = AppProps & {
  Component: NextPageWithLayout;
};

export default function App({ Component, pageProps }: AppPropsWithLayout) {
  const layout = Component.getLayout ?? (page => page);

  return layout(
    <RecoilRoot>
      <Component {...pageProps} />
    </RecoilRoot>,
  );
}

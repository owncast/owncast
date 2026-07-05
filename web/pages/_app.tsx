// All these imports are almost exclusively for the Admin.
// We should not be loading them for the main frontend UI.

// order matters!
import '../styles/variables.css';
import '../styles/global.less';
import '../styles/globals.scss';
import '../styles/ant-overrides.scss';
// Pre-extracted Ant Design v6 styles (see build-scripts/extract-antd6-styles.js).
import '../styles/antd6.css';

// TODO: Move this videojs sass to the player component.
import '../components/video/VideoJS/VideoJS.scss';

import { AppProps } from 'next/app';
import { ReactElement, ReactNode, useEffect } from 'react';
import { NextPage } from 'next';
import { RecoilRoot } from 'recoil';
import { useRouter } from 'next/router';
import { useSelectedLanguage } from 'next-export-i18n';
import { loadViewerLocale, shouldCleanUrlAfterFlip } from '../utils/localeLoader';
import { Antd6Provider } from '../components/theme/Antd6Provider';

export type NextPageWithLayout<P = {}, IP = P> = NextPage<P, IP> & {
  getLayout?: (page: ReactElement) => ReactNode;
};

type AppPropsWithLayout = AppProps & {
  Component: NextPageWithLayout;
};

// Removes the ?lang= parameter the locale loader adds to trigger a language
// switch, once that switch has verifiably landed (this component's own
// selected language changed to it). Copied URLs stay clean instead of
// pinning the viewer's language onto whoever opens them. Explicit viewer-set
// ?lang= values are never registered for cleanup, so they stay put.
function LocaleSync() {
  const router = useRouter();
  const { lang } = useSelectedLanguage();
  useEffect(() => {
    if (shouldCleanUrlAfterFlip(lang)) {
      const rest = { ...router.query };
      delete rest.lang;
      router.replace({ query: rest }, undefined, { shallow: true, scroll: false });
    }
  }, [lang]);
  return null;
}

export default function App({ Component, pageProps }: AppPropsWithLayout) {
  const layout = Component.getLayout ?? (page => page);

  // Fetch the viewer's locale catalog (one small chunk) once the router can
  // report ?lang=. Everything renders in English until it lands.
  const router = useRouter();
  useEffect(() => {
    if (router.isReady) {
      loadViewerLocale(router);
    }
  }, [router.isReady]);

  return (
    <>
      <LocaleSync />
      {/* RecoilRoot and the Ant Design v6 theme bridge wrap the page
          layout too, not just the page component: layouts attached via
          getLayout (e.g. the admin's AdminLayout/MainLayout) also render
          antd6 components and must be inside the provider, or they fall
          back to the unscoped default prefix and collide with v4's global
          styles. */}
      <RecoilRoot>
        <Antd6Provider>{layout(<Component {...pageProps} />)}</Antd6Provider>
      </RecoilRoot>
    </>
  );
}

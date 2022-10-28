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
import { Router, useRouter } from 'next/router';

import { RecoilRoot } from 'recoil';
import { useEffect } from 'react';
import { AdminLayout } from '../components/layouts/AdminLayout';
import { SimpleLayout } from '../components/layouts/SimpleLayout';

const App = ({ Component, pageProps }: AppProps) => {
  useEffect(() => {
    if ('serviceWorker' in navigator) {
      window.addEventListener('load', () => {
        navigator.serviceWorker.register('/serviceWorker.js').then(
          registration => {
            console.debug(
              'Service Worker registration successful with scope: ',
              registration.scope,
            );
          },
          err => {
            console.error('Service Worker registration failed: ', err);
          },
        );
      });
    }
  }, []);

  const router = useRouter() as Router;
  if (router.pathname.startsWith('/admin')) {
    return <AdminLayout pageProps={pageProps} Component={Component} router={router} />;
  }
  return (
    <RecoilRoot>
      <SimpleLayout pageProps={pageProps} Component={Component} router={router} />
    </RecoilRoot>
  );
};

export default App;

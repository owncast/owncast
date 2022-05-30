// All these imports are almost exclusively for the Admin.
// We should not be loading them for the main frontend UI.

// order matters!
import '../styles/variables.css';
import '../styles/global.less';
import '../styles/globals.scss';
// import '../styles/ant-overrides.scss';
// import '../styles/markdown-editor.scss';

import '../styles/main-layout.scss';

import '../styles/form-textfields.scss';
import '../styles/form-misc-elements.scss';
import '../styles/config-socialhandles.scss';
import '../styles/config-storage.scss';
import '../styles/config-edit-string-tags.scss';
import '../styles/config-video-variants.scss';
import '../styles/config-public-details.scss';

import '../styles/home.scss';
import '../styles/chat.scss';
import '../styles/pages.scss';
import '../styles/offline-notice.scss';

import '../components/video/player.scss';

import { AppProps } from 'next/app';
import { Router, useRouter } from 'next/router';

import { RecoilRoot } from 'recoil';
import { useEffect } from 'react';
import AdminLayout from '../components/layouts/admin-layout';
import SimpleLayout from '../components/layouts/SimpleLayout';

function App({ Component, pageProps }: AppProps) {
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
}

export default App;

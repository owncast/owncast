/* eslint-disable react/no-invalid-html-attribute */
/* eslint-disable react/no-danger */
/* eslint-disable react/no-unescaped-entities */
import { useRecoilValue } from 'recoil';
import Head from 'next/head';
import { FC, useEffect, useRef } from 'react';
import { Layout } from 'antd';
import dynamic from 'next/dynamic';
import Script from 'next/script';
import { ErrorBoundary } from 'react-error-boundary';
import {
  ClientConfigStore,
  isChatAvailableSelector,
  clientConfigStateAtom,
  fatalErrorStateAtom,
  appStateAtom,
  serverStatusState,
} from '../../stores/ClientConfigStore';
import { Content } from '../../ui/Content/Content';
import { Header } from '../../ui/Header/Header';
import { ClientConfig } from '../../../interfaces/client-config.model';
import { DisplayableError } from '../../../types/displayable-error';
import setupNoLinkReferrer from '../../../utils/no-link-referrer';
import { TitleNotifier } from '../../TitleNotifier/TitleNotifier';
import { ServerRenderedHydration } from '../../ServerRendered/ServerRenderedHydration';
import { Theme } from '../../theme/Theme';
import styles from './Main.module.scss';
import { PushNotificationServiceWorker } from '../../workers/PushNotificationServiceWorker/PushNotificationServiceWorker';
import { AppStateOptions } from '../../stores/application-state';
import { Noscript } from '../../ui/Noscript/Noscript';
import { ServerStatus } from '../../../interfaces/server-status.model';

// Lazy loaded components

const FatalErrorStateModal = dynamic(
  () =>
    import('../../modals/FatalErrorStateModal/FatalErrorStateModal').then(
      mod => mod.FatalErrorStateModal,
    ),
  {
    ssr: false,
  },
);

export const Main: FC = () => {
  const clientConfig = useRecoilValue<ClientConfig>(clientConfigStateAtom);
  const clientStatus = useRecoilValue<ServerStatus>(serverStatusState);
  const { name } = clientConfig;
  const isChatAvailable = useRecoilValue<boolean>(isChatAvailableSelector);
  const fatalError = useRecoilValue<DisplayableError>(fatalErrorStateAtom);
  const appState = useRecoilValue<AppStateOptions>(appStateAtom);
  const layoutRef = useRef<HTMLDivElement>(null);
  const { chatDisabled } = clientConfig;
  const { videoAvailable } = appState;
  const { online, streamTitle } = clientStatus;

  useEffect(() => {
    setupNoLinkReferrer(layoutRef.current);
  }, []);

  const isProduction = process.env.NODE_ENV === 'production';
  const headerText = online ? streamTitle || name : name;

  return (
    <>
      <Head>
        {isProduction && <ServerRenderedHydration />}

        <link rel="icon" href="/favicon.ico" />
        <link rel="manifest" href="/manifest.json" />
        <link rel="authorization_endpoint" href="/api/auth/provider/indieauth" />
        <meta
          name="viewport"
          content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=no"
        />
        <meta name="apple-mobile-web-app-status-bar-style" content="black-translucent" />

        <base target="_blank" />
      </Head>

      {isProduction ? (
        <Head>{name ? <title>{name}</title> : <title>{'{{.Name}}'}</title>}</Head>
      ) : (
        <Head>
          <title>{name}</title>
        </Head>
      )}
      <ErrorBoundary
        // eslint-disable-next-line react/no-unstable-nested-components
        fallbackRender={({ error }) => (
          <FatalErrorStateModal
            title="Error"
            message={`There was an unexpected error. Please refresh the page to retry. If this error continues please file a bug with the Owncast project: ${error}`}
          />
        )}
      >
        <ClientConfigStore />
      </ErrorBoundary>
      <PushNotificationServiceWorker />
      <TitleNotifier name={name} />
      <Theme />
      <Script strategy="afterInteractive" src="/customjavascript" />
      <Layout ref={layoutRef} className={styles.layout}>
        <Header
          name={headerText}
          chatAvailable={isChatAvailable}
          chatDisabled={chatDisabled}
          online={videoAvailable}
        />
        <Content />
        {fatalError && (
          <FatalErrorStateModal title={fatalError.title} message={fatalError.message} />
        )}
      </Layout>
      <Noscript />
    </>
  );
};

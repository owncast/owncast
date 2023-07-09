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

        <link rel="apple-touch-icon" sizes="57x57" href="/img/favicon/apple-icon-57x57.png" />
        <link rel="apple-touch-icon" sizes="60x60" href="/img/favicon/apple-icon-60x60.png" />
        <link rel="apple-touch-icon" sizes="72x72" href="/img/favicon/apple-icon-72x72.png" />
        <link rel="apple-touch-icon" sizes="76x76" href="/img/favicon/apple-icon-76x76.png" />
        <link rel="apple-touch-icon" sizes="114x114" href="/img/favicon/apple-icon-114x114.png" />
        <link rel="apple-touch-icon" sizes="120x120" href="/img/favicon/apple-icon-120x120.png" />
        <link rel="apple-touch-icon" sizes="144x144" href="/img/favicon/apple-icon-144x144.png" />
        <link rel="apple-touch-icon" sizes="152x152" href="/img/favicon/apple-icon-152x152.png" />
        <link rel="apple-touch-icon" sizes="180x180" href="/img/favicon/apple-icon-180x180.png" />
        <link
          rel="icon"
          type="image/png"
          sizes="192x192"
          href="/img/favicon/android-icon-192x192.png"
        />
        <link rel="icon" type="image/png" sizes="32x32" href="/img/favicon/favicon-32x32.png" />
        <link rel="icon" type="image/png" sizes="96x96" href="/img/favicon/favicon-96x96.png" />
        <link rel="icon" type="image/png" sizes="16x16" href="/img/favicon/favicon-16x16.png" />
        <link rel="manifest" href="/manifest.json" />
        <link rel="authorization_endpoint" href="/api/auth/provider/indieauth" />
        <meta name="msapplication-TileColor" content="#ffffff" />
        <meta name="msapplication-TileImage" content="/img/favicon/ms-icon-144x144.png" />
        <meta
          name="viewport"
          content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=no"
        />
        <meta name="apple-mobile-web-app-status-bar-style" content="black-translucent" />

        <base target="_blank" />
      </Head>

      {isProduction ? (
        <Head>
          {name ? <title>{name}</title> : <title>{'{{.Name}}'}</title>}
          <meta name="description" content="{{.Summary}}" />

          <meta property="og:title" content="{{.Name}}" />
          <meta property="og:site_name" content="{{.Name}}" />
          <meta property="og:url" content="{{.RequestedURL}}" />
          <meta property="og:description" content="{{.Summary}}" />
          <meta property="og:type" content="video.other" />
          <meta property="video:tag" content="{{.TagsString}}" />

          <meta property="og:image" content="{{.RequestedURL}}{{.Thumbnail}}" />
          <meta property="og:image:url" content="{{.RequestedURL}}{{.Thumbnail}}" />
          <meta property="og:image:alt" content="{{.RequestedURL}}{{.Image}}" />

          <meta property="og:video" content="{{.RequestedURL}}{{.EmbedVideo}}" />
          <meta property="og:video:secure_url" content="{{.RequestedURL}}{{.EmbedVideo}}" />
          <meta property="og:video:height" content="315" />
          <meta property="og:video:width" content="560" />
          <meta property="og:video:type" content="text/html" />
          <meta property="og:video:actor" content="{{.Name}}" />

          <meta property="twitter:title" content="{{.Name}}" />
          <meta property="twitter:url" content="{{.RequestedURL}}" />
          <meta property="twitter:description" content="{{.Summary}}" />
          <meta property="twitter:image" content="{{.Image}}" />
          <meta property="twitter:card" content="player" />
          <meta property="twitter:player" content="{{.RequestedURL}}{{.EmbedVideo}}" />
          <meta property="twitter:player:width" content="560" />
          <meta property="twitter:player:height" content="315" />
        </Head>
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

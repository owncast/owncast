/* eslint-disable react/no-danger */
/* eslint-disable react/no-unescaped-entities */
import { Layout } from 'antd';
import { useRecoilValue } from 'recoil';
import Head from 'next/head';
import { FC, useEffect, useRef } from 'react';
import {
  ClientConfigStore,
  isChatAvailableSelector,
  clientConfigStateAtom,
  fatalErrorStateAtom,
} from '../stores/ClientConfigStore';
import { Content } from '../ui/Content/Content';
import { Header } from '../ui/Header/Header';
import { ClientConfig } from '../../interfaces/client-config.model';
import { DisplayableError } from '../../types/displayable-error';
import { FatalErrorStateModal } from '../modals/FatalErrorStateModal/FatalErrorStateModal';
import setupNoLinkReferrer from '../../utils/no-link-referrer';
import { ServerRenderedMetadata } from '../ServerRendered/ServerRenderedMetadata';
import { ServerRenderedHydration } from '../ServerRendered/ServerRenderedHydration';

export const Main: FC = () => {
  const clientConfig = useRecoilValue<ClientConfig>(clientConfigStateAtom);
  const { name, title, customStyles } = clientConfig;
  const isChatAvailable = useRecoilValue<boolean>(isChatAvailableSelector);
  const fatalError = useRecoilValue<DisplayableError>(fatalErrorStateAtom);

  const layoutRef = useRef<HTMLDivElement>(null);
  const { chatDisabled } = clientConfig;

  useEffect(() => {
    setupNoLinkReferrer(layoutRef.current);
  }, []);

  const isProduction = process.env.NODE_ENV === 'production';

  const hydrationScript = isProduction
    ? `
	window.statusHydration = {{.StatusJSON}};
	window.configHydration = {{.ServerConfigJSON}};
	`
    : null;

  return (
    <>
      <Head>
        {isProduction && <ServerRenderedMetadata />}

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

        <link href="/api/auth/provider/indieauth" />

        <meta name="msapplication-TileColor" content="#ffffff" />
        <meta name="msapplication-TileImage" content="/img/favicon/ms-icon-144x144.png" />
        <meta name="theme-color" content="#ffffff" />

        <style>{customStyles}</style>
        {isProduction && <ServerRenderedHydration hydrationScript={hydrationScript} />}
      </Head>

      <ClientConfigStore />
      <Layout ref={layoutRef} style={{ minHeight: '100vh' }}>
        <Header name={title || name} chatAvailable={isChatAvailable} chatDisabled={chatDisabled} />
        <Content />
        {fatalError && (
          <FatalErrorStateModal title={fatalError.title} message={fatalError.message} />
        )}
      </Layout>
    </>
  );
};

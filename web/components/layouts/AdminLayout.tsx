import { AppProps } from 'next/app';
import { FC } from 'react';
import ServerStatusProvider from '../../utils/server-status-context';
import AlertMessageProvider from '../../utils/alert-message-context';
import MainLayout from '../MainLayout';

const AdminLayout: FC<AppProps> = ({ Component, pageProps }) => (
  <ServerStatusProvider>
    <AlertMessageProvider>
      <MainLayout>
        <Component {...pageProps} />
      </MainLayout>
    </AlertMessageProvider>
  </ServerStatusProvider>
);

export default AdminLayout;

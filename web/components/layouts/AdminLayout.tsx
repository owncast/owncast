import { AppProps } from 'next/app';
import ServerStatusProvider from '../../utils/server-status-context';
import AlertMessageProvider from '../../utils/alert-message-context';
import MainLayout from '../MainLayout';

const AdminLayout = ({ Component, pageProps }: AppProps) => (
  <ServerStatusProvider>
    <AlertMessageProvider>
      <MainLayout>
        <Component {...pageProps} />
      </MainLayout>
    </AlertMessageProvider>
  </ServerStatusProvider>
);

export default AdminLayout;

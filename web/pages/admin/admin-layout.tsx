import { AppProps } from 'next/app';
import ServerStatusProvider from '../../utils/server-status-context';
import AlertMessageProvider from '../../utils/alert-message-context';
import MainLayout from '../../components/main-layout';

function AdminLayout({ Component, pageProps }: AppProps) {
  return (
    <ServerStatusProvider>
      <AlertMessageProvider>
        <MainLayout>
          <Component {...pageProps} />
        </MainLayout>
      </AlertMessageProvider>
    </ServerStatusProvider>
  );
}

export default AdminLayout;

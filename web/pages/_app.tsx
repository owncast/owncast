// order matters!
import 'antd/dist/antd.css';
import '../styles/colors.scss';
import '../styles/globals.scss';
import '../styles/ant-overrides.scss';
import '../styles/markdown-editor.scss';

import '../styles/main-layout.scss';

import '../styles/form-textfields.scss';
import '../styles/form-toggleswitch.scss';
import '../styles/form-misc-elements.scss';
import '../styles/config-socialhandles.scss';
import '../styles/config-storage.scss';
import '../styles/config-tags.scss';
import '../styles/config-video-variants.scss';

import '../styles/home.scss';
import '../styles/chat.scss';
import '../styles/config.scss';

import { AppProps } from 'next/app';
import ServerStatusProvider from '../utils/server-status-context';
import AlertMessageProvider from '../utils/alert-message-context';

import MainLayout from '../components/main-layout';

function App({ Component, pageProps }: AppProps) {
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

export default App;

import 'antd/dist/antd.css';
import '../styles/colors.scss';
import '../styles/globals.scss';

import '../styles/home.scss';
import '../styles/chat.scss';
import '../styles/config.scss';

import { AppProps } from 'next/app';
import ServerStatusProvider from '../utils/server-status-context';
import MainLayout from './components/main-layout';


function App({ Component, pageProps }: AppProps) {
  return (
    <ServerStatusProvider>
      <MainLayout>
        <Component {...pageProps} />
      </MainLayout>
    </ServerStatusProvider>
    
  )
}

export default App;
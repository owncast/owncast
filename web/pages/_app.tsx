import 'antd/dist/antd.compact.css';
import "../styles/globals.scss";

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
import 'antd/dist/antd.dark.css';
import 'antd/dist/antd.compact.css';
import "../styles/globals.scss";

import { AppProps } from 'next/app';
import BroadcastStatusProvider from './utils/broadcast-status-context';
import MainLayout from './components/main-layout';


function App({ Component, pageProps }: AppProps) {
  return (
    <BroadcastStatusProvider>
      <MainLayout>
        <Component {...pageProps} />
      </MainLayout>
    </BroadcastStatusProvider>
    
  )
}

export default App;
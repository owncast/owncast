import 'antd/dist/antd.compact.css';
import '../styles/colors.scss';
import '../styles/globals.scss';

// GW: I can't override ant design styles through components using NextJS's built-in CSS modules. So I'll just import styles here for now and figure out enabling SASS modules later.
import '../styles/home.scss';

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
// import 'antd/dist/antd.css';
// import '../styles/globals.scss'

import 'antd/dist/antd.dark.css';
import 'antd/dist/antd.compact.css';

import "../styles/globals.scss";

import { AppProps } from 'next/app'

function App({ Component, pageProps }: AppProps) {
  return <Component {...pageProps} />
}

export default App
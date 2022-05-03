import { Layout } from 'antd';
import { ServerStatusStore } from '../stores/ServerStatusStore';
import { ClientConfigStore } from '../stores/ClientConfigStore';
import { Content, Header } from '../ui';

function Main() {
  return (
    <>
      <ServerStatusStore />
      <ClientConfigStore />
      <Layout>
        <Header />
        <Content />
      </Layout>
    </>
  );
}

export default Main;

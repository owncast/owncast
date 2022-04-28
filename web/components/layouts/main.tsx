import { Layout } from 'antd';
import { ServerStatusStore } from '../stores/ServerStatusStore';
import { ClientConfigStore } from '../stores/ClientConfigStore';
import { Content, Footer, Header, Sidebar } from '../ui';

function Main() {
  return (
    <>
      <ServerStatusStore />
      <ClientConfigStore />
      <Layout>
        <Sidebar />
        <Header />
        <Layout className="site-layout" style={{ marginRight: 200 }}>
          <Content />
          <Footer />
        </Layout>
      </Layout>
    </>
  );
}

export default Main;

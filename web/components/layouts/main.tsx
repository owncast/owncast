import { useRecoilValue } from 'recoil';
import { Layout, Row, Col } from 'antd';
import { useState } from 'react';
import { ServerStatus } from '../../models/ServerStatus';
import { ServerStatusStore, serverStatusState } from '../stores/ServerStatusStore';
import { ClientConfigStore, clientConfigState } from '../stores/ClientConfigStore';
import { ClientConfig } from '../../models/ClientConfig';

const { Header, Content, Footer, Sider } = Layout;

function Main() {
  const serverStatus = useRecoilValue<ServerStatus>(serverStatusState);
  const clientConfig = useRecoilValue<ClientConfig>(clientConfigState);

  const { name, version, extraPageContent } = clientConfig;
  const [chatCollapsed, setChatCollapsed] = useState(false);

  const toggleChatCollapsed = () => {
    setChatCollapsed(!chatCollapsed);
  };

  return (
    <>
      <ServerStatusStore />
      <ClientConfigStore />

      <Layout>
        <Sider
          collapsed={chatCollapsed}
          width={300}
          style={{
            position: 'fixed',
            right: 0,
            top: 0,
            bottom: 0,
          }}
        />
        <Layout className="site-layout" style={{ marginRight: 200 }}>
          <Header
            className="site-layout-background"
            style={{ position: 'fixed', zIndex: 1, width: '100%' }}
          >
            {name}
            <button onClick={toggleChatCollapsed}>Toggle Chat</button>
          </Header>
          <Content style={{ margin: '80px 16px 0', overflow: 'initial' }}>
            <div>
              <Row>
                <Col span={24}>Video player goes here</Col>
              </Row>
              <Row>
                <Col span={24}>
                  <Content dangerouslySetInnerHTML={{ __html: extraPageContent }} />
                </Col>
              </Row>
            </div>
          </Content>
          <Footer style={{ textAlign: 'center' }}>Footer: Owncast {version}</Footer>
        </Layout>
      </Layout>
    </>
  );
}

export default Main;

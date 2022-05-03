import { useRecoilValue } from 'recoil';
import { Layout, Row, Col, Tabs } from 'antd';
import { clientConfigStateAtom } from '../../stores/ClientConfigStore';
import { ClientConfig } from '../../../interfaces/client-config.model';
import CustomPageContent from '../../CustomPageContent';
import OwncastPlayer from '../../video/OwncastPlayer';
import FollowerCollection from '../../FollowersCollection';

const { TabPane } = Tabs;

const { Content } = Layout;

export default function FooterComponent() {
  const clientConfig = useRecoilValue<ClientConfig>(clientConfigStateAtom);
  const { extraPageContent } = clientConfig;

  return (
    <Content style={{ margin: '80px 16px 0', overflow: 'initial' }}>
      <div>
        <Row>
          <Col span={24}>
            <OwncastPlayer source="https://watch.owncast.online" />
          </Col>
        </Row>
        <Row>
          <Col span={24}>
            <Tabs defaultActiveKey="1" type="card">
              <TabPane tab="About" key="1">
                <CustomPageContent content={extraPageContent} />
              </TabPane>
              <TabPane tab="Followers" key="2">
                <FollowerCollection />
              </TabPane>
            </Tabs>
          </Col>
        </Row>
      </div>
    </Content>
  );
}

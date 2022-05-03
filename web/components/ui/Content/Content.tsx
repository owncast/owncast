import { useRecoilValue } from 'recoil';
import { Layout, Row, Col, Tabs } from 'antd';
import { chatVisibilityAtom, clientConfigStateAtom } from '../../stores/ClientConfigStore';
import { ClientConfig } from '../../../interfaces/client-config.model';
import CustomPageContent from '../../CustomPageContent';
import OwncastPlayer from '../../video/OwncastPlayer';
import FollowerCollection from '../../FollowersCollection';
import s from './Content.module.scss';
import Sidebar from '../Sidebar';
import { ChatVisibilityState } from '../../../interfaces/application-state';
import Footer from '../Footer';
import Grid from 'antd/lib/card/Grid';

const { TabPane } = Tabs;

const { Content } = Layout;

export default function FooterComponent() {
  const clientConfig = useRecoilValue<ClientConfig>(clientConfigStateAtom);
  const chatOpen = useRecoilValue<ChatVisibilityState>(chatVisibilityAtom);
  const { extraPageContent } = clientConfig;

  return (
    <Content className={`${s.root}`}>
      <Col className={`${s.leftCol}`}>
        <OwncastPlayer source="https://watch.owncast.online" />
        <div className={`${s.lowerRow}`}>
          <Tabs defaultActiveKey="1" type="card">
            <TabPane tab="About" key="1">
              <CustomPageContent content={extraPageContent} />
            </TabPane>
            <TabPane tab="Followers" key="2">
              <FollowerCollection />
            </TabPane>
          </Tabs>
          <Footer />
        </div>
      </Col>
      {chatOpen && (
        <Col>
          <Sidebar />
        </Col>
      )}
    </Content>
  );
}

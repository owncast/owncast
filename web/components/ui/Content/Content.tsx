import { useRecoilValue } from 'recoil';
import { Layout, Tabs } from 'antd';
import { chatVisibilityAtom, clientConfigStateAtom } from '../../stores/ClientConfigStore';
import { ClientConfig } from '../../../interfaces/client-config.model';
import CustomPageContent from '../../CustomPageContent';
import OwncastPlayer from '../../video/OwncastPlayer';
import FollowerCollection from '../../FollowersCollection';
import s from './Content.module.scss';
import Sidebar from '../Sidebar';
import Footer from '../Footer';
import ChatContainer from '../../chat/ChatContainer';
import { ChatMessage } from '../../../interfaces/chat-message.model';
import { chatMessagesAtom, chatStateAtom } from '../../stores/ClientConfigStore';
import { ChatState, ChatVisibilityState } from '../../../interfaces/application-state';
import ChatTextField from '../../chat/ChatTextField/ChatTextField';

const { TabPane } = Tabs;

const { Content } = Layout;

export default function FooterComponent() {
  const clientConfig = useRecoilValue<ClientConfig>(clientConfigStateAtom);
  const chatOpen = useRecoilValue<ChatVisibilityState>(chatVisibilityAtom);
  const messages = useRecoilValue<ChatMessage[]>(chatMessagesAtom);
  const chatState = useRecoilValue<ChatState>(chatStateAtom);

  const { extraPageContent } = clientConfig;

  return (
    <Content className={`${s.root}`} data-columns={chatOpen ? 2 : 1}>
      <div className={`${s.leftCol}`}>
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
          {chatOpen && (
            <div className={`${s.mobileChat}`}>
              <ChatContainer messages={messages} state={chatState} />
              <ChatTextField />
            </div>
          )}
          <Footer />
        </div>
      </div>
      {chatOpen && <Sidebar />}
    </Content>
  );
}

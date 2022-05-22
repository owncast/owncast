import { useRecoilValue } from 'recoil';
import { Layout, Button, Tabs } from 'antd';
import { NotificationFilled, HeartFilled } from '@ant-design/icons';
import {
  chatVisibilityAtom,
  clientConfigStateAtom,
  chatMessagesAtom,
  chatStateAtom,
  serverStatusState,
} from '../../stores/ClientConfigStore';
import { ClientConfig } from '../../../interfaces/client-config.model';
import CustomPageContent from '../../CustomPageContent';
import OwncastPlayer from '../../video/OwncastPlayer';
import FollowerCollection from '../../FollowersCollection';
import s from './Content.module.scss';
import Sidebar from '../Sidebar';
import Footer from '../Footer';
import ChatContainer from '../../chat/ChatContainer';
import { ChatMessage } from '../../../interfaces/chat-message.model';
import { ChatState, ChatVisibilityState } from '../../../interfaces/application-state';
import ChatTextField from '../../chat/ChatTextField/ChatTextField';
import ActionButtonRow from '../../action-buttons/ActionButtonRow';
import ActionButton from '../../action-buttons/ActionButton';
import Statusbar from '../Statusbar/Statusbar';
import { ServerStatus } from '../../../interfaces/server-status.model';
import { Follower } from '../../../interfaces/follower';
import SocialLinks from '../SocialLinks/SocialLinks';
import NotifyReminderPopup from '../NotifyReminderPopup/NotifyReminderPopup';
import ServerLogo from '../Logo/Logo';
import CategoryIcon from '../CategoryIcon/CategoryIcon';

const { TabPane } = Tabs;
const { Content } = Layout;

export default function ContentComponent() {
  const status = useRecoilValue<ServerStatus>(serverStatusState);
  const clientConfig = useRecoilValue<ClientConfig>(clientConfigStateAtom);
  const chatVisibility = useRecoilValue<ChatVisibilityState>(chatVisibilityAtom);
  const messages = useRecoilValue<ChatMessage[]>(chatMessagesAtom);
  const chatState = useRecoilValue<ChatState>(chatStateAtom);

  const { extraPageContent, version, socialHandles, name, title, tags } = clientConfig;
  const { online, viewerCount, lastConnectTime, lastDisconnectTime } = status;

  const followers: Follower[] = [];

  const total = 0;

  const chatVisible =
    chatState === ChatState.Available && chatVisibility === ChatVisibilityState.Visible;

  // This is example content. It should be removed.
  const externalActions = [
    {
      url: 'https://owncast.online/docs',
      title: 'Example button',
      description: 'Example button description',
      icon: 'https://owncast.online/images/logo.svg',
      color: '#5232c8',
      openExternally: false,
    },
  ];

  const externalActionButtons = externalActions.map(action => (
    <ActionButton key={action.url} action={action} />
  ));

  return (
    <Content className={`${s.root}`}>
      <div className={`${s.leftCol}`}>
        <OwncastPlayer source="/hls/stream.m3u8" online={online} />
        <Statusbar
          online={online}
          lastConnectTime={lastConnectTime}
          lastDisconnectTime={lastDisconnectTime}
          viewerCount={viewerCount}
        />
        <div className={s.buttonsLogoTitleSection}>
          <ActionButtonRow>
            {externalActionButtons}
            <Button icon={<HeartFilled />}>Follow</Button>
            <NotifyReminderPopup
              visible
              notificationClicked={() => {}}
              notificationClosed={() => {}}
            >
              <Button icon={<NotificationFilled />}>Notify</Button>
            </NotifyReminderPopup>
          </ActionButtonRow>

          <div className={`${s.lowerRow}`}>
            <div className={s.logoTitleSection}>
              <ServerLogo src="/logo" />
              <div className={s.titleSection}>
                <div className={s.title}>{name}</div>
                <div className={s.subtitle}>
                  {title}
                  <CategoryIcon tags={tags} />
                </div>
                <div>{tags.length > 0 && tags.map(tag => <span key={tag}>#{tag}&nbsp;</span>)}</div>
              </div>
            </div>
          </div>

          <Tabs defaultActiveKey="1">
            <TabPane tab="About" key="1" className={`${s.pageContentSection}`}>
              <SocialLinks links={socialHandles} />
              <CustomPageContent content={extraPageContent} />
            </TabPane>
            <TabPane tab="Followers" key="2" className={`${s.pageContentSection}`}>
              <FollowerCollection total={total} followers={followers} />
            </TabPane>
          </Tabs>
          {chatVisibility && (
            <div className={`${s.mobileChat}`}>
              <ChatContainer messages={messages} state={chatState} />
              <ChatTextField />
            </div>
          )}
          <Footer version={version} />
        </div>
      </div>
      {chatVisible && <Sidebar />}
    </Content>
  );
}

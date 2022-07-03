/* eslint-disable react/no-danger */
import { useRecoilValue, useSetRecoilState } from 'recoil';
import { Layout, Tabs, Spin } from 'antd';
import { useEffect, useState } from 'react';
import cn from 'classnames';
import { LOCAL_STORAGE_KEYS, getLocalStorage, setLocalStorage } from '../../../utils/localStorage';

import {
  clientConfigStateAtom,
  chatMessagesAtom,
  chatDisplayNameAtom,
  chatUserIdAtom,
  isChatVisibleSelector,
  serverStatusState,
  appStateAtom,
  isOnlineSelector,
  isMobileAtom,
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
import OfflineBanner from '../OfflineBanner/OfflineBanner';
import { AppStateOptions } from '../../stores/application-state';
import FollowButton from '../../action-buttons/FollowButton';
import NotifyButton from '../../action-buttons/NotifyButton';
import Modal from '../Modal/Modal';
import BrowserNotifyModal from '../../modals/BrowserNotify/BrowserNotifyModal';

const { TabPane } = Tabs;
const { Content } = Layout;

export default function ContentComponent() {
  const appState = useRecoilValue<AppStateOptions>(appStateAtom);
  const status = useRecoilValue<ServerStatus>(serverStatusState);
  const clientConfig = useRecoilValue<ClientConfig>(clientConfigStateAtom);
  const isChatVisible = useRecoilValue<boolean>(isChatVisibleSelector);
  const setIsMobile = useSetRecoilState<boolean>(isMobileAtom);
  const isMobile = useRecoilValue(isMobileAtom);
  const messages = useRecoilValue<ChatMessage[]>(chatMessagesAtom);
  const online = useRecoilValue<boolean>(isOnlineSelector);
  const chatDisplayName = useRecoilValue<string>(chatDisplayNameAtom);
  const chatUserId = useRecoilValue<string>(chatUserIdAtom);

  const { extraPageContent, version, socialHandles, name, title, tags, summary } = clientConfig;
  const { viewerCount, lastConnectTime, lastDisconnectTime } = status;
  const [showNotifyReminder, setShowNotifyReminder] = useState(false);
  const [showNotifyPopup, setShowNotifyPopup] = useState(false);

  const followers: Follower[] = [];

  const total = 0;

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

  const incrementVisitCounter = () => {
    let visits = parseInt(getLocalStorage(LOCAL_STORAGE_KEYS.userVisitCount), 10);
    if (Number.isNaN(visits)) {
      visits = 0;
    }

    setLocalStorage(LOCAL_STORAGE_KEYS.userVisitCount, visits + 1);

    if (visits > 2 && !getLocalStorage(LOCAL_STORAGE_KEYS.hasDisplayedNotificationModal)) {
      setShowNotifyReminder(true);
    }
  };

  const disableNotifyReminderPopup = () => {
    setShowNotifyPopup(false);
    setShowNotifyReminder(false);
    setLocalStorage(LOCAL_STORAGE_KEYS.hasDisplayedNotificationModal, true);
  };

  const checkIfMobile = () => {
    const w = window.innerWidth;
    if (isMobile === undefined) {
      if (w <= 768) setIsMobile(true);
      else setIsMobile(false);
    }
    if (!isMobile && w <= 768) setIsMobile(true);
    if (isMobile && w > 768) setIsMobile(false);
  };

  useEffect(() => {
    incrementVisitCounter();
    checkIfMobile();
    window.addEventListener('resize', checkIfMobile);
  }, []);

  const rootClassName = cn(s.root, {
    [s.mobile]: isMobile,
  });

  return (
    <Content className={rootClassName}>
      <Spin className={s.loadingSpinner} size="large" spinning={appState.appLoading} />

      <div className={s.leftCol}>
        {online && <OwncastPlayer source="/hls/stream.m3u8" online={online} />}
        {!online && (
          <OfflineBanner
            name={name}
            text="Stream is offline text goes here. Will create a new form to set it in the Admin."
          />
        )}

        <Statusbar
          online={online}
          lastConnectTime={lastConnectTime}
          lastDisconnectTime={lastDisconnectTime}
          viewerCount={viewerCount}
        />
        <div className={s.buttonsLogoTitleSection}>
          <ActionButtonRow>
            {externalActionButtons}
            <FollowButton />
            <NotifyReminderPopup
              visible={showNotifyReminder}
              notificationClicked={() => setShowNotifyPopup(true)}
              notificationClosed={() => disableNotifyReminderPopup()}
            >
              <NotifyButton onClick={() => setShowNotifyPopup(true)} />
            </NotifyReminderPopup>
          </ActionButtonRow>

          <Modal
            title="Notify"
            visible={showNotifyPopup}
            afterClose={() => disableNotifyReminderPopup()}
            handleCancel={() => disableNotifyReminderPopup()}
          >
            <BrowserNotifyModal />
          </Modal>

          <div className={s.streamInfo}>
            <div className={s.logoTitleSection}>
              <ServerLogo src="/logo" />
              <div className={s.titleSection}>
                <div className={s.title}>{name}</div>
                <div className={s.subtitle}>
                  {title}
                  <CategoryIcon tags={tags} />
                </div>
                <div>{tags.length > 0 && tags.map(tag => <span key={tag}>#{tag}&nbsp;</span>)}</div>
                <SocialLinks links={socialHandles} />
              </div>
            </div>
          </div>

          <Tabs defaultActiveKey="1">
            <TabPane tab="About" key="1" className={s.pageContentSection}>
              <div dangerouslySetInnerHTML={{ __html: summary }} />
              <CustomPageContent content={extraPageContent} />
            </TabPane>
            <TabPane tab="Followers" key="2" className={s.pageContentSection}>
              <FollowerCollection total={total} followers={followers} />
            </TabPane>
          </Tabs>
          {isChatVisible && isMobile && (
            <div style={{ position: 'relative' }}>
              <div className={s.mobileChat}>
                <ChatContainer
                  messages={messages}
                  loading={appState.chatLoading}
                  usernameToHighlight={chatDisplayName}
                  chatUserId={chatUserId}
                  isModerator={false}
                  isMobile={isMobile}
                />
              </div>
              <ChatTextField />
            </div>
          )}
          <Footer version={version} />
        </div>
      </div>
      {isChatVisible && !isMobile && <Sidebar />}
    </Content>
  );
}

/* eslint-disable react/no-danger */
import { useRecoilState, useRecoilValue } from 'recoil';
import { Layout, Tabs, Spin } from 'antd';
import { useEffect, useState } from 'react';
import cn from 'classnames';
import { LOCAL_STORAGE_KEYS, getLocalStorage, setLocalStorage } from '../../../utils/localStorage';

import {
  clientConfigStateAtom,
  // chatMessagesAtom,
  // chatDisplayNameAtom,
  // chatUserIdAtom,
  isChatVisibleSelector,
  appStateAtom,
  isOnlineSelector,
  isMobileAtom,
  serverStatusState,
} from '../../stores/ClientConfigStore';
import { ClientConfig } from '../../../interfaces/client-config.model';
import CustomPageContent from '../CustomPageContent/CustomPageContent';
import OwncastPlayer from '../../video/OwncastPlayer';
import FollowerCollection from '../Followers/FollowersCollection';
import s from './Content.module.scss';
import Sidebar from '../Sidebar';
import Footer from '../Footer';
// import ChatContainer from '../../chat/ChatContainer';
// import { ChatMessage } from '../../../interfaces/chat-message.model';
// import ChatTextField from '../../chat/ChatTextField/ChatTextField';
import ActionButtonRow from '../../action-buttons/ActionButtonRow';
import ActionButton from '../../action-buttons/ActionButton';
import NotifyReminderPopup from '../NotifyReminderPopup/NotifyReminderPopup';
import OfflineBanner from '../OfflineBanner/OfflineBanner';
import { AppStateOptions } from '../../stores/application-state';
import FollowButton from '../../action-buttons/FollowButton';
import NotifyButton from '../../action-buttons/NotifyButton';
import Modal from '../Modal/Modal';
import BrowserNotifyModal from '../../modals/BrowserNotify/BrowserNotifyModal';
import ContentHeader from '../../common/ContentHeader';
import { ServerStatus } from '../../../interfaces/server-status.model';
import { StatusBar } from '..';

const { TabPane } = Tabs;
const { Content } = Layout;

export default function ContentComponent() {
  const appState = useRecoilValue<AppStateOptions>(appStateAtom);
  const clientConfig = useRecoilValue<ClientConfig>(clientConfigStateAtom);
  const isChatVisible = useRecoilValue<boolean>(isChatVisibleSelector);
  const [isMobile, setIsMobile] = useRecoilState<boolean | undefined>(isMobileAtom);
  // const messages = useRecoilValue<ChatMessage[]>(chatMessagesAtom);
  const online = useRecoilValue<boolean>(isOnlineSelector);
  // const chatDisplayName = useRecoilValue<string>(chatDisplayNameAtom);
  // const chatUserId = useRecoilValue<string>(chatUserIdAtom);
  const { viewerCount, lastConnectTime, lastDisconnectTime, streamTitle } =
    useRecoilValue<ServerStatus>(serverStatusState);
  const {
    extraPageContent,
    version,
    name,
    summary,
    socialHandles,
    tags,
    externalActions,
    offlineMessage,
  } = clientConfig;
  const [showNotifyReminder, setShowNotifyReminder] = useState(false);
  const [showNotifyPopup, setShowNotifyPopup] = useState(false);

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
    <div>
      <Content className={rootClassName}>
        <div className={s.leftContent}>
          <Spin className={s.loadingSpinner} size="large" spinning={appState.appLoading} />

          <div className={s.topSection}>
            {online && <OwncastPlayer source="/hls/stream.m3u8" online={online} />}
            {!online && (
              <OfflineBanner
                name={name}
                text={
                  offlineMessage || 'Please follow and ask to get notified when the stream is live.'
                }
              />
            )}
            <StatusBar
              online={online}
              lastConnectTime={lastConnectTime}
              lastDisconnectTime={lastDisconnectTime}
              viewerCount={viewerCount}
            />
          </div>
          <div className={s.midSection}>
            <div className={s.buttonsLogoTitleSection}>
              <ActionButtonRow>
                {externalActionButtons}
                <FollowButton size="small" />
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
            </div>
          </div>
          <div className={s.lowerHalf}>
            <ContentHeader
              name={name}
              title={streamTitle}
              summary={summary}
              tags={tags}
              links={socialHandles}
              logo="/logo"
            />
          </div>

          <div className={s.lowerSection}>
            <Tabs defaultActiveKey="0" style={{ height: '100%' }}>
              {isChatVisible && isMobile && (
                <TabPane tab="Chat" key="0" style={{ height: '100%' }}>
                  {/* <div style={{ position: 'relative', height: '100%' }}>
                  <div className={s.mobileChat}>
                    <ChatContainer
                      messages={messages}
                      // loading={appState.chatLoading}
                      usernameToHighlight={chatDisplayName}
                      chatUserId={chatUserId}
                      isModerator={false}
                    />
                  </div>
                </div> */}
                </TabPane>
              )}
              <TabPane tab="About" key="2" className={s.pageContentSection}>
                <CustomPageContent content={extraPageContent} />
              </TabPane>
              <TabPane tab="Followers" key="3" className={s.pageContentSection}>
                <FollowerCollection />
              </TabPane>
            </Tabs>
          </div>
        </div>
        {isChatVisible && !isMobile && <Sidebar />}
      </Content>
      {!isMobile && <Footer version={version} />}
    </div>
  );
}

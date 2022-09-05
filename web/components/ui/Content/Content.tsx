/* eslint-disable react/no-danger */
import { useRecoilState, useRecoilValue } from 'recoil';
import { Layout, Tabs, Spin } from 'antd';
import { FC, useEffect, useState } from 'react';
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
import { CustomPageContent } from '../CustomPageContent/CustomPageContent';
import { OwncastPlayer } from '../../video/OwncastPlayer/OwncastPlayer';
import { FollowerCollection } from '../followers/FollowerCollection/FollowerCollection';
import styles from './Content.module.scss';
import { Sidebar } from '../Sidebar/Sidebar';
import { Footer } from '../Footer/Footer';
// import ChatContainer from '../../chat/ChatContainer';
// import { ChatMessage } from '../../../interfaces/chat-message.model';
// import ChatTextField from '../../chat/ChatTextField/ChatTextField';
import { ActionButtonRow } from '../../action-buttons/ActionButtonRow/ActionButtonRow';
import { ActionButton } from '../../action-buttons/ActionButton/ActionButton';
import { NotifyReminderPopup } from '../NotifyReminderPopup/NotifyReminderPopup';
import { OfflineBanner } from '../OfflineBanner/OfflineBanner';
import { AppStateOptions } from '../../stores/application-state';
import { FollowButton } from '../../action-buttons/FollowButton';
import { NotifyButton } from '../../action-buttons/NotifyButton';
import { Modal } from '../Modal/Modal';
import { BrowserNotifyModal } from '../../modals/BrowserNotifyModal/BrowserNotifyModal';
import { ContentHeader } from '../../common/ContentHeader/ContentHeader';
import { ServerStatus } from '../../../interfaces/server-status.model';
import { Statusbar } from '../Statusbar/Statusbar';

const { TabPane } = Tabs;
const { Content: AntContent } = Layout;

export const Content: FC = () => {
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

  const rootClassName = cn(styles.root, {
    [styles.mobile]: isMobile,
  });

  return (
    <div>
      <AntContent className={rootClassName}>
        <div className={styles.leftContent}>
          <Spin className={styles.loadingSpinner} size="large" spinning={appState.appLoading} />

          <div className={styles.topSection}>
            {online && <OwncastPlayer source="/hls/stream.m3u8" online={online} />}
            {!online && (
              <OfflineBanner
                name={name}
                text={
                  offlineMessage || 'Please follow and ask to get notified when the stream is live.'
                }
              />
            )}
            <Statusbar
              online={online}
              lastConnectTime={lastConnectTime}
              lastDisconnectTime={lastDisconnectTime}
              viewerCount={viewerCount}
            />
          </div>
          <div className={styles.midSection}>
            <div className={styles.buttonsLogoTitleSection}>
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
          <div className={styles.lowerHalf}>
            <ContentHeader
              name={name}
              title={streamTitle}
              summary={summary}
              tags={tags}
              links={socialHandles}
              logo="/logo"
            />
          </div>

          <div className={styles.lowerSection}>
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
              <TabPane tab="About" key="2" className={styles.pageContentSection}>
                <CustomPageContent content={extraPageContent} />
              </TabPane>
              <TabPane tab="Followers" key="3" className={styles.pageContentSection}>
                <FollowerCollection />
              </TabPane>
            </Tabs>
          </div>
        </div>
        {isChatVisible && !isMobile && <Sidebar />}
      </AntContent>
      {!isMobile && <Footer version={version} />}
    </div>
  );
};
export default Content;

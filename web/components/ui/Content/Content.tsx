import { useRecoilState, useRecoilValue } from 'recoil';
import { Layout, Tabs, Spin } from 'antd';
import { FC, useEffect, useState } from 'react';
import classNames from 'classnames';
import { LOCAL_STORAGE_KEYS, getLocalStorage, setLocalStorage } from '../../../utils/localStorage';

import {
  clientConfigStateAtom,
  chatMessagesAtom,
  chatDisplayNameAtom,
  chatUserIdAtom,
  isChatAvailableSelector,
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

import { ActionButtonRow } from '../../action-buttons/ActionButtonRow/ActionButtonRow';
import { ActionButton } from '../../action-buttons/ActionButton/ActionButton';
import { NotifyReminderPopup } from '../NotifyReminderPopup/NotifyReminderPopup';
import { OfflineBanner } from '../OfflineBanner/OfflineBanner';
import { AppStateOptions } from '../../stores/application-state';
import { FollowButton } from '../../action-buttons/FollowButton';
import { NotifyButton } from '../../action-buttons/NotifyButton';
import { Modal } from '../../atomic/molecules/Modal/Modal';
import { BrowserNotifyModal } from '../../modals/BrowserNotifyModal/BrowserNotifyModal';
import { ContentHeader } from '../../common/ContentHeader/ContentHeader';
import { ServerStatus } from '../../../interfaces/server-status.model';
import { Statusbar } from '../Statusbar/Statusbar';
import { ChatContainer } from '../../chat/ChatContainer/ChatContainer';
import { ChatMessage } from '../../../interfaces/chat-message.model';

const { TabPane } = Tabs;
const { Content: AntContent } = Layout;

const DesktopContent = ({ name, streamTitle, summary, tags, socialHandles, extraPageContent }) => (
  <>
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
      <Tabs defaultActiveKey="0">
        <TabPane tab="About" key="2">
          <CustomPageContent content={extraPageContent} />
        </TabPane>
        <TabPane tab="Followers" key="3">
          <FollowerCollection />
        </TabPane>
      </Tabs>
    </div>
  </>
);

const MobileContent = ({
  name,
  streamTitle,
  summary,
  tags,
  socialHandles,
  extraPageContent,
  messages,
  chatDisplayName,
  chatUserId,
  showChat,
}) => (
  <div className={classNames(styles.lowerSectionMobile)}>
    <Tabs defaultActiveKey="0">
      {showChat && (
        <TabPane tab="Chat" key="1">
          <ChatContainer
            messages={messages}
            usernameToHighlight={chatDisplayName}
            chatUserId={chatUserId}
            isModerator={false}
            height="40vh"
          />
        </TabPane>
      )}
      <TabPane tab="About" key="2">
        <ContentHeader
          name={name}
          title={streamTitle}
          summary={summary}
          tags={tags}
          links={socialHandles}
          logo="/logo"
        />
        <CustomPageContent content={extraPageContent} />
      </TabPane>
      <TabPane tab="Followers" key="3">
        <FollowerCollection />
      </TabPane>
    </Tabs>
  </div>
);

export const Content: FC = () => {
  const appState = useRecoilValue<AppStateOptions>(appStateAtom);
  const clientConfig = useRecoilValue<ClientConfig>(clientConfigStateAtom);
  const isChatVisible = useRecoilValue<boolean>(isChatVisibleSelector);
  const isChatAvailable = useRecoilValue<boolean>(isChatAvailableSelector);

  const [isMobile, setIsMobile] = useRecoilState<boolean | undefined>(isMobileAtom);
  const messages = useRecoilValue<ChatMessage[]>(chatMessagesAtom);
  const online = useRecoilValue<boolean>(isOnlineSelector);
  const chatDisplayName = useRecoilValue<string>(chatDisplayNameAtom);
  const chatUserId = useRecoilValue<string>(chatUserIdAtom);
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
    chatDisabled,
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

  let offlineMessageBody =
    !appState.appLoading && 'Please follow and ask to get notified when the stream is live.';
  if (offlineMessage && !appState.appLoading) {
    offlineMessageBody = offlineMessage;
  }

  const offlineTitle = !appState.appLoading && `${name} is currently offline`;
  const showChat = !chatDisabled && isChatAvailable && isChatVisible;

  return (
    <div>
      <Spin className={styles.loadingSpinner} size="large" spinning={appState.appLoading}>
        <AntContent className={styles.root}>
          <div className={styles.leftContent}>
            <div className={styles.topSection}>
              {online && <OwncastPlayer source="/hls/stream.m3u8" online={online} />}
              {!online && !appState.appLoading && (
                <OfflineBanner title={offlineTitle} text={offlineMessageBody} />
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
            {isMobile && isChatVisible ? (
              <MobileContent
                name={name}
                streamTitle={streamTitle}
                summary={summary}
                tags={tags}
                socialHandles={socialHandles}
                extraPageContent={extraPageContent}
                messages={messages}
                chatDisplayName={chatDisplayName}
                chatUserId={chatUserId}
                showChat={showChat}
              />
            ) : (
              <DesktopContent
                name={name}
                streamTitle={streamTitle}
                summary={summary}
                tags={tags}
                socialHandles={socialHandles}
                extraPageContent={extraPageContent}
              />
            )}
          </div>
          {showChat && !isMobile && <Sidebar />}
        </AntContent>
        {(!isMobile || !showChat) && <Footer version={version} />}
      </Spin>
    </div>
  );
};
export default Content;

import { useRecoilState, useRecoilValue } from 'recoil';
import { Layout, Tabs, Spin } from 'antd';
import { FC, useEffect, useState } from 'react';
import classNames from 'classnames';
import dynamic from 'next/dynamic';
import { LOCAL_STORAGE_KEYS, getLocalStorage, setLocalStorage } from '~/utils/localStorage';

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
} from '~/components/stores/ClientConfigStore';
import { ClientConfig } from '~/interfaces/client-config.model';
import { CustomPageContent } from '~/components/ui/CustomPageContent/CustomPageContent';
import { OwncastPlayer } from '~/components/video/OwncastPlayer/OwncastPlayer';
import { Sidebar } from '~/components/ui/Sidebar/Sidebar';
import { Footer } from '~/components/ui/Footer/Footer';
import { ActionButtonRow } from '~/components/action-buttons/ActionButtonRow/ActionButtonRow';
import { ActionButton } from '~/components/action-buttons/ActionButton/ActionButton';
import { OfflineBanner } from '~/components/ui/OfflineBanner/OfflineBanner';
import { AppStateOptions } from '~/components/stores/application-state';
import { FollowButton } from '~/components/action-buttons/FollowButton';
import { NotifyButton } from '~/components/action-buttons/NotifyButton';
import { ContentHeader } from '~/components/common/ContentHeader/ContentHeader';
import { ServerStatus } from '~/interfaces/server-status.model';
import { Statusbar } from '~/components/ui/Statusbar/Statusbar';
import { ChatMessage } from '~/interfaces/chat-message.model';
import styles from './Content.module.scss';

const { TabPane } = Tabs;
const { Content: AntContent } = Layout;

// Lazy loaded components

const Modal = dynamic(() => import('../Modal/Modal').then(mod => mod.Modal));

const BrowserNotifyModal = dynamic(() =>
  import('~/components/modals/BrowserNotifyModal/BrowserNotifyModal').then(
    mod => mod.BrowserNotifyModal,
  ),
);

const NotifyReminderPopup = dynamic(() =>
  import('~/components/ui/NotifyReminderPopup/NotifyReminderPopup').then(
    mod => mod.NotifyReminderPopup,
  ),
);

const FollowerCollection = dynamic(() =>
  import('~/components/ui/followers/FollowerCollection/FollowerCollection').then(
    mod => mod.FollowerCollection,
  ),
);

// We only need to load the chat container here if we're in mobile or narrow
// windows, so lazy loading it makes sense.
const ChatContainer = dynamic(() =>
  import('~/components/chat/ChatContainer/ChatContainer').then(mod => mod.ChatContainer),
);

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

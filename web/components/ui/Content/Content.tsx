import { useRecoilState, useRecoilValue } from 'recoil';
import { Layout, Tabs, Spin } from 'antd';
import { FC, MutableRefObject, useEffect, useRef, useState } from 'react';
import cn from 'classnames';
import dynamic from 'next/dynamic';
import { LOCAL_STORAGE_KEYS, getLocalStorage, setLocalStorage } from '../../../utils/localStorage';

import {
  clientConfigStateAtom,
  chatMessagesAtom,
  currentUserAtom,
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
import styles from './Content.module.scss';
import { Sidebar } from '../Sidebar/Sidebar';
import { Footer } from '../Footer/Footer';

import { ActionButtonRow } from '../../action-buttons/ActionButtonRow/ActionButtonRow';
import { ActionButton } from '../../action-buttons/ActionButton/ActionButton';
import { OfflineBanner } from '../OfflineBanner/OfflineBanner';
import { AppStateOptions } from '../../stores/application-state';
import { FollowButton } from '../../action-buttons/FollowButton';
import { NotifyButton } from '../../action-buttons/NotifyButton';
import { ContentHeader } from '../../common/ContentHeader/ContentHeader';
import { ServerStatus } from '../../../interfaces/server-status.model';
import { Statusbar } from '../Statusbar/Statusbar';
import { ChatMessage } from '../../../interfaces/chat-message.model';
import { FollowerCollection } from '../followers/FollowerCollection/FollowerCollection';
import { ExternalAction } from '../../../interfaces/external-action';
import { Modal } from '../Modal/Modal';

const { Content: AntContent } = Layout;

const BrowserNotifyModal = dynamic(() =>
  import('../../modals/BrowserNotifyModal/BrowserNotifyModal').then(mod => mod.BrowserNotifyModal),
);

const NotifyReminderPopup = dynamic(() =>
  import('../NotifyReminderPopup/NotifyReminderPopup').then(mod => mod.NotifyReminderPopup),
);

// We only need to load the chat container here if we're in mobile or narrow
// windows, so lazy loading it makes sense.
const ChatContainer = dynamic(() =>
  import('../../chat/ChatContainer/ChatContainer').then(mod => mod.ChatContainer),
);

const DesktopContent = ({ name, streamTitle, summary, tags, socialHandles, extraPageContent }) => {
  const aboutTabContent = <CustomPageContent content={extraPageContent} />;
  const followersTabContent = <FollowerCollection name={name} />;

  const items = [
    { label: 'About', key: '2', children: aboutTabContent },
    { label: 'Followers', key: '3', children: followersTabContent },
  ];

  return (
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
        <Tabs defaultActiveKey="0" items={items} />
      </div>
    </>
  );
};

function useHeight(ref: MutableRefObject<HTMLDivElement>) {
  const [contentH, setContentH] = useState(0);
  const handleResize = () => {
    if (!ref.current) return;
    const fromTop = ref.current.getBoundingClientRect().top;
    const { innerHeight } = window;
    setContentH(innerHeight - fromTop);
  };

  useEffect(() => {
    handleResize();
    window.addEventListener('resize', handleResize);

    return () => {
      window.removeEventListener('resize', handleResize);
    };
  }, []);

  return contentH;
}

const MobileContent = ({
  name,
  streamTitle,
  summary,
  tags,
  socialHandles,
  extraPageContent,
  messages,
  currentUser,
  showChat,
}) => {
  if (!currentUser) {
    return null;
  }

  const mobileContentRef = useRef<HTMLDivElement>();

  const { id, displayName } = currentUser;

  const chatContent = showChat && (
    <ChatContainer
      messages={messages}
      usernameToHighlight={displayName}
      chatUserId={id}
      isModerator={false}
    />
  );

  const aboutTabContent = (
    <>
      <ContentHeader
        name={name}
        title={streamTitle}
        summary={summary}
        tags={tags}
        links={socialHandles}
        logo="/logo"
      />
      <CustomPageContent content={extraPageContent} />
    </>
  );
  const followersTabContent = <FollowerCollection name={name} />;

  const items = [
    { label: 'Chat', key: '1', children: chatContent },
    { label: 'About', key: '2', children: aboutTabContent },
    { label: 'Followers', key: '3', children: followersTabContent },
  ];

  const height = `${useHeight(mobileContentRef)}px`;

  return (
    <div className={cn(styles.lowerSectionMobile)} ref={mobileContentRef} style={{ height }}>
      <Tabs defaultActiveKey="0" items={items} />
    </div>
  );
};

const ExternalModal = ({ externalActionToDisplay, setExternalActionToDisplay }) => {
  const { title, description, url } = externalActionToDisplay;
  return (
    <Modal
      title={description || title}
      url={url}
      open={!!externalActionToDisplay}
      height="80vh"
      handleCancel={() => setExternalActionToDisplay(null)}
    />
  );
};

export const Content: FC = () => {
  const appState = useRecoilValue<AppStateOptions>(appStateAtom);
  const clientConfig = useRecoilValue<ClientConfig>(clientConfigStateAtom);
  const isChatVisible = useRecoilValue<boolean>(isChatVisibleSelector);
  const isChatAvailable = useRecoilValue<boolean>(isChatAvailableSelector);
  const currentUser = useRecoilValue(currentUserAtom);

  const [isMobile, setIsMobile] = useRecoilState<boolean | undefined>(isMobileAtom);
  const messages = useRecoilValue<ChatMessage[]>(chatMessagesAtom);
  const online = useRecoilValue<boolean>(isOnlineSelector);

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
    federation,
    notifications,
  } = clientConfig;
  const [showNotifyReminder, setShowNotifyReminder] = useState(false);
  const [showNotifyPopup, setShowNotifyPopup] = useState(false);
  const { account: fediverseAccount } = federation;
  const { browser: browserNotifications } = notifications;
  const { enabled: browserNotificationsEnabled } = browserNotifications;
  const [externalActionToDisplay, setExternalActionToDisplay] = useState<ExternalAction>(null);

  const externalActionSelected = (action: ExternalAction) => {
    const { openExternally, url } = action;
    if (openExternally) {
      window.open(url, '_blank');
    } else {
      setExternalActionToDisplay(action);
    }
  };

  const externalActionButtons = externalActions.map(action => (
    <ActionButton
      key={action.url}
      action={action}
      externalActionSelected={externalActionSelected}
    />
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
    return () => {
      window.removeEventListener('resize', checkIfMobile);
    };
  }, []);

  const showChat = !chatDisabled && isChatAvailable && isChatVisible;

  return (
    <>
      <div className={styles.main}>
        <Spin wrapperClassName={styles.loadingSpinner} size="large" spinning={appState.appLoading}>
          <AntContent className={styles.root}>
            <div className={styles.mainSection}>
              <div className={styles.topSection}>
                {online && <OwncastPlayer source="/hls/stream.m3u8" online={online} />}
                {!online && !appState.appLoading && (
                  <OfflineBanner
                    streamName={name}
                    customText={offlineMessage}
                    notificationsEnabled={browserNotificationsEnabled}
                    fediverseAccount={fediverseAccount}
                    lastLive={lastDisconnectTime}
                    onNotifyClick={() => setShowNotifyPopup(true)}
                  />
                )}
                {online && (
                  <Statusbar
                    online={online}
                    lastConnectTime={lastConnectTime}
                    lastDisconnectTime={lastDisconnectTime}
                    viewerCount={viewerCount}
                  />
                )}
              </div>
              <div className={styles.midSection}>
                <div className={styles.buttonsLogoTitleSection}>
                  <ActionButtonRow>
                    {externalActionButtons}
                    <FollowButton size="small" />
                    <NotifyReminderPopup
                      open={showNotifyReminder}
                      notificationClicked={() => setShowNotifyPopup(true)}
                      notificationClosed={() => disableNotifyReminderPopup()}
                    >
                      <NotifyButton onClick={() => setShowNotifyPopup(true)} />
                    </NotifyReminderPopup>
                  </ActionButtonRow>

                  <Modal
                    title="Notify"
                    open={showNotifyPopup}
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
                  currentUser={currentUser}
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
        </Spin>
        {!isMobile && <Footer version={version} />}
      </div>

      {externalActionToDisplay && (
        <ExternalModal
          externalActionToDisplay={externalActionToDisplay}
          setExternalActionToDisplay={setExternalActionToDisplay}
        />
      )}
    </>
  );
};
export default Content;

import { useRecoilState, useRecoilValue } from 'recoil';
import { Skeleton } from 'antd';
import { FC, useEffect, useState } from 'react';
import dynamic from 'next/dynamic';
import classnames from 'classnames';
import { LOCAL_STORAGE_KEYS, getLocalStorage, setLocalStorage } from '../../../utils/localStorage';
import isPushNotificationSupported from '../../../utils/browserPushNotifications';

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
import { ExternalAction } from '../../../interfaces/external-action';
import { Modal } from '../Modal/Modal';
import { ActionButtonMenu } from '../../action-buttons/ActionButtonMenu/ActionButtonMenu';

// Lazy loaded components

const FollowerCollection = dynamic(
  () =>
    import('../followers/FollowerCollection/FollowerCollection').then(
      mod => mod.FollowerCollection,
    ),
  {
    ssr: false,
  },
);

const FollowModal = dynamic(
  () => import('../../modals/FollowModal/FollowModal').then(mod => mod.FollowModal),
  {
    ssr: false,
    loading: () => <Skeleton loading active paragraph={{ rows: 8 }} />,
  },
);

const BrowserNotifyModal = dynamic(
  () =>
    import('../../modals/BrowserNotifyModal/BrowserNotifyModal').then(
      mod => mod.BrowserNotifyModal,
    ),
  {
    ssr: false,
    loading: () => <Skeleton loading active paragraph={{ rows: 6 }} />,
  },
);

const NotifyReminderPopup = dynamic(
  () => import('../NotifyReminderPopup/NotifyReminderPopup').then(mod => mod.NotifyReminderPopup),
  {
    ssr: false,
    loading: () => <Skeleton loading active paragraph={{ rows: 8 }} />,
  },
);

const OwncastPlayer = dynamic(
  () => import('../../video/OwncastPlayer/OwncastPlayer').then(mod => mod.OwncastPlayer),
  {
    ssr: false,
    loading: () => <Skeleton loading active paragraph={{ rows: 12 }} />,
  },
);

const ChatContainer = dynamic(
  () => import('../../chat/ChatContainer/ChatContainer').then(mod => mod.ChatContainer),
  {
    ssr: false,
  },
);

const Tabs = dynamic(() => import('antd').then(mod => mod.Tabs), {
  ssr: false,
});

const DesktopContent = ({
  name,
  streamTitle,
  summary,
  tags,
  socialHandles,
  extraPageContent,
  setShowFollowModal,
  supportFediverseFeatures,
}) => {
  const aboutTabContent = <CustomPageContent content={extraPageContent} />;
  const followersTabContent = (
    <div>
      <FollowerCollection name={name} onFollowButtonClick={() => setShowFollowModal(true)} />
    </div>
  );

  const items = [{ label: 'About', key: '2', children: aboutTabContent }];
  if (supportFediverseFeatures) {
    items.push({ label: 'Followers', key: '3', children: followersTabContent });
  }

  return (
    <>
      <div className={styles.lowerHalf} id="skip-to-content">
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
        {items.length > 1 ? <Tabs defaultActiveKey="0" items={items} /> : aboutTabContent}
      </div>
    </>
  );
};

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
  actions,
  setExternalActionToDisplay,
  setShowNotifyPopup,
  setShowFollowModal,
  supportFediverseFeatures,
  supportsBrowserNotifications,
}) => {
  if (!currentUser) {
    return <Skeleton loading active paragraph={{ rows: 7 }} />;
  }
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
  const followersTabContent = (
    <FollowerCollection name={name} onFollowButtonClick={() => setShowFollowModal(true)} />
  );

  const items = [
    showChat && { label: 'Chat', key: '0', children: chatContent },
    { label: 'About', key: '2', children: aboutTabContent },
    { label: 'Followers', key: '3', children: followersTabContent },
  ];

  const replacementTabBar = (props, DefaultTabBar) => (
    <div className={styles.replacementBar}>
      <DefaultTabBar {...props} className={styles.defaultTabBar} />
      <ActionButtonMenu
        className={styles.actionButtonMenu}
        showFollowItem={supportFediverseFeatures}
        showNotifyItem={supportsBrowserNotifications}
        actions={actions}
        externalActionSelected={setExternalActionToDisplay}
        notifyItemSelected={() => setShowNotifyPopup(true)}
        followItemSelected={() => setShowFollowModal(true)}
      />
    </div>
  );

  return (
    <div className={styles.lowerSectionMobile}>
      <Tabs
        className={styles.tabs}
        defaultActiveKey="0"
        items={items}
        renderTabBar={replacementTabBar}
      />
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
  const serverStatus = useRecoilValue<ServerStatus>(serverStatusState);
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
  const [showNotifyModal, setShowNotifyModal] = useState(false);
  const [showFollowModal, setShowFollowModal] = useState(false);
  const { account: fediverseAccount, enabled: fediverseEnabled } = federation;
  const { browser: browserNotifications } = notifications;
  const { enabled: browserNotificationsEnabled } = browserNotifications;
  const { online: isStreamLive } = serverStatus;
  const [externalActionToDisplay, setExternalActionToDisplay] = useState<ExternalAction>(null);

  const [supportsBrowserNotifications, setSupportsBrowserNotifications] = useState(false);
  const supportFediverseFeatures = fediverseEnabled;

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
    setShowNotifyModal(false);
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

  useEffect(() => {
    // isPushNotificationSupported relies on `navigator` so that needs to be
    // fired from this useEffect.
    setSupportsBrowserNotifications(isPushNotificationSupported() && browserNotificationsEnabled);
  }, [browserNotificationsEnabled]);

  const showChat = !chatDisabled && isChatAvailable && isChatVisible;

  return (
    <>
      <div className={styles.root}>
        <div className={classnames(styles.mainSection, { [styles.offline]: !online })}>
          {appState.appLoading ? (
            <Skeleton loading active paragraph={{ rows: 7 }} className={styles.topSectionElement} />
          ) : (
            <div className="skeleton-placeholder" />
          )}
          {online && (
            <OwncastPlayer
              source="/hls/stream.m3u8"
              online={online}
              title={streamTitle || name}
              className={styles.topSectionElement}
            />
          )}
          {!online && !appState.appLoading && (
            <div id="offline-message">
              <OfflineBanner
                showsHeader={false}
                streamName={name}
                customText={offlineMessage}
                notificationsEnabled={browserNotificationsEnabled}
                fediverseAccount={fediverseAccount}
                lastLive={lastDisconnectTime}
                onNotifyClick={() => setShowNotifyModal(true)}
                onFollowClick={() => setShowFollowModal(true)}
                className={styles.topSectionElement}
              />
            </div>
          )}
          {isStreamLive ? (
            <Statusbar
              online={online}
              lastConnectTime={lastConnectTime}
              lastDisconnectTime={lastDisconnectTime}
              viewerCount={viewerCount}
              className={classnames(styles.topSectionElement, styles.statusBar)}
            />
          ) : (
            <div className="statusbar-placeholder" />
          )}
          <div className={styles.midSection}>
            <div className={styles.buttonsLogoTitleSection}>
              {!isMobile && (
                <ActionButtonRow>
                  {externalActionButtons}
                  {supportFediverseFeatures && (
                    <FollowButton size="small" onClick={() => setShowFollowModal(true)} />
                  )}
                  {supportsBrowserNotifications && (
                    <NotifyReminderPopup
                      open={showNotifyReminder}
                      notificationClicked={() => setShowNotifyModal(true)}
                      notificationClosed={() => disableNotifyReminderPopup()}
                    >
                      <NotifyButton onClick={() => setShowNotifyModal(true)} />
                    </NotifyReminderPopup>
                  )}
                </ActionButtonRow>
              )}

              <Modal
                title="Browser Notifications"
                open={showNotifyModal}
                afterClose={() => disableNotifyReminderPopup()}
                handleCancel={() => disableNotifyReminderPopup()}
              >
                <BrowserNotifyModal />
              </Modal>
            </div>
          </div>
          {isMobile ? (
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
              actions={externalActions}
              setExternalActionToDisplay={externalActionSelected}
              setShowNotifyPopup={setShowNotifyModal}
              setShowFollowModal={setShowFollowModal}
              supportFediverseFeatures={supportFediverseFeatures}
              supportsBrowserNotifications={supportsBrowserNotifications}
            />
          ) : (
            <DesktopContent
              name={name}
              streamTitle={streamTitle}
              summary={summary}
              tags={tags}
              socialHandles={socialHandles}
              extraPageContent={extraPageContent}
              setShowFollowModal={setShowFollowModal}
              supportFediverseFeatures={supportFediverseFeatures}
            />
          )}
          {!isMobile && <Footer version={version} />}
        </div>
        {showChat && !isMobile && <Sidebar />}
      </div>
      {externalActionToDisplay && (
        <ExternalModal
          externalActionToDisplay={externalActionToDisplay}
          setExternalActionToDisplay={setExternalActionToDisplay}
        />
      )}
      <Modal
        title={`Follow ${name}`}
        open={showFollowModal}
        handleCancel={() => setShowFollowModal(false)}
        width="550px"
      >
        <FollowModal
          account={fediverseAccount}
          name={name}
          handleClose={() => setShowFollowModal(false)}
        />
      </Modal>
    </>
  );
};
export default Content;

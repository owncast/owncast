import { useRecoilState, useRecoilValue } from 'recoil';
import { Skeleton, Col, Row } from 'antd';
import { FC, useEffect, useState } from 'react';
import dynamic from 'next/dynamic';
import classnames from 'classnames';
import ActionButtons from './ActionButtons';
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

import styles from './Content.module.scss';
import desktopStyles from './DesktopContent.module.scss';
import { Sidebar } from '../Sidebar/Sidebar';
import { OfflineBanner } from '../OfflineBanner/OfflineBanner';
import { AppStateOptions } from '../../stores/application-state';
import { ServerStatus } from '../../../interfaces/server-status.model';
import { Statusbar } from '../Statusbar/Statusbar';
import { ChatMessage } from '../../../interfaces/chat-message.model';
import { ExternalAction } from '../../../interfaces/external-action';
import { Modal } from '../Modal/Modal';
import { DesktopContent } from './DesktopContent';
import { MobileContent } from './MobileContent';

// Lazy loaded components

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

const OwncastPlayer = dynamic(
  () => import('../../video/OwncastPlayer/OwncastPlayer').then(mod => mod.OwncastPlayer),
  {
    ssr: false,
    loading: () => <Skeleton loading active paragraph={{ rows: 12 }} />,
  },
);

const ExternalModal = ({ externalActionToDisplay, setExternalActionToDisplay }) => {
  const { title, description, url, html } = externalActionToDisplay;
  return (
    <Modal
      title={description || title}
      url={url}
      open={!!externalActionToDisplay}
      height="80vh"
      handleCancel={() => setExternalActionToDisplay(null)}
    >
      {html ? (
        <div
          // eslint-disable-next-line react/no-danger
          dangerouslySetInnerHTML={{ __html: html }}
          style={{
            height: '100%',
            width: '100%',
            overflow: 'auto',
          }}
        />
      ) : null}
    </Modal>
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
    // apply openExternally only if we don't have an HTML embed
    if (openExternally && url) {
      window.open(url, '_blank');
    } else {
      setExternalActionToDisplay(action);
    }
  };

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

  const showChat = online && !chatDisabled && isChatVisible;

  // accounts for sidebar width when online in desktop
  const dynamicPadding = online && !isMobile ? '320px' : '0px';

  return (
    <>
      <>
        {appState.appLoading && (
          <Skeleton loading active paragraph={{ rows: 7 }} className={styles.topSectionElement} />
        )}
        {showChat && !isMobile && <Sidebar />}
        <Row>
          <Col span={24} style={{ paddingRight: dynamicPadding }}>
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
                  notificationsEnabled={supportsBrowserNotifications}
                  fediverseAccount={fediverseAccount}
                  lastLive={lastDisconnectTime}
                  onNotifyClick={() => setShowNotifyModal(true)}
                  onFollowClick={() => setShowFollowModal(true)}
                  className={classnames([styles.topSectionElement, styles.offlineBanner])}
                />
              </div>
            )}
          </Col>
        </Row>
        <Row>
          <Col span={24} style={{ paddingRight: dynamicPadding }}>
            {isStreamLive && (
              <Statusbar
                online={online}
                lastConnectTime={lastConnectTime}
                lastDisconnectTime={lastDisconnectTime}
                viewerCount={viewerCount}
                className={classnames(styles.topSectionElement, styles.statusBar)}
              />
            )}
          </Col>
        </Row>
        <Row>
          <Col span={24} style={{ paddingRight: dynamicPadding }}>
            <ActionButtons
              isMobile={isMobile}
              supportFediverseFeatures={supportFediverseFeatures}
              supportsBrowserNotifications={supportsBrowserNotifications}
              showNotifyReminder={showNotifyReminder}
              setShowNotifyModal={setShowNotifyModal}
              disableNotifyReminderPopup={disableNotifyReminderPopup}
              externalActions={externalActions}
              setExternalActionToDisplay={setExternalActionToDisplay}
              setShowFollowModal={setShowFollowModal}
              externalActionSelected={externalActionSelected}
            />
          </Col>
        </Row>

        <Modal
          title="Browser Notifications"
          open={showNotifyModal}
          afterClose={() => disableNotifyReminderPopup()}
          handleCancel={() => disableNotifyReminderPopup()}
        >
          <BrowserNotifyModal />
        </Modal>
        <Row>
          {isMobile ? (
            <MobileContent
              name={name}
              summary={summary}
              tags={tags}
              socialHandles={socialHandles}
              extraPageContent={extraPageContent}
              messages={messages}
              currentUser={currentUser}
              showChat={showChat}
              setShowFollowModal={setShowFollowModal}
              supportFediverseFeatures={supportFediverseFeatures}
              chatEnabled={isChatAvailable}
            />
          ) : (
            <Col span={24} style={{ paddingRight: dynamicPadding }}>
              <div className={desktopStyles.bottomSectionContent}>
                <DesktopContent
                  name={name}
                  summary={summary}
                  tags={tags}
                  socialHandles={socialHandles}
                  extraPageContent={extraPageContent}
                  setShowFollowModal={setShowFollowModal}
                  supportFediverseFeatures={supportFediverseFeatures}
                />
              </div>
            </Col>
          )}
        </Row>
      </>
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

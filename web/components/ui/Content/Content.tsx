import { useAtom, useAtomValue } from 'jotai';
import { Skeleton, Row, Button, Spin } from 'antd';
import MessageFilled from '@ant-design/icons/MessageFilled';
import { FC, useCallback, useEffect, useState } from 'react';
import classnames from 'classnames';
import dynamic from 'next/dynamic';
import { useTranslation } from 'next-export-i18n';
import ActionButtons from './ActionButtons';
import { LOCAL_STORAGE_KEYS, getLocalStorage, setLocalStorage } from '../../../utils/localStorage';
import { canPushNotificationsBeSupported } from '../../../utils/browserPushNotifications';
import { getVideoPosition } from '../../../utils/video-latency';
import { Localization } from '../../../types/localization';

import {
  clientConfigStateAtom,
  currentUserAtom,
  ChatState,
  chatStateAtom,
  appStateAtom,
  isOnlineSelector,
  isMobileAtom,
  serverStatusState,
  isChatAvailableSelector,
  visibleChatMessagesSelector,
  chatAuthenticatedAtom,
  isClientConfigLoadedAtom,
  canShowClipButtonSelector,
  canCreateClipSelector,
} from '../../stores/ClientConfigStore';

import styles from './Content.module.scss';
import desktopStyles from './DesktopContent.module.scss';
import { OfflineBanner } from '../OfflineBanner/OfflineBanner';
import { Statusbar } from '../Statusbar/Statusbar';
import { ExternalAction } from '../../../interfaces/external-action';
import { Modal } from '../Modal/Modal';
import { DesktopContent } from './DesktopContent';
import { MobileContent } from './MobileContent';
import { Footer } from '../Footer/Footer';
import { useFederatedServers } from '../../../hooks/useFederatedServers';

// Lazy loaded components
const ChatContainer = dynamic(
  () => import('../../chat/ChatContainer/ChatContainer').then(mod => mod.ChatContainer),
  {
    ssr: false,
  },
);

// The follow modal renders its own antd v6 modal shell (the first surface
// migrated to v6), so it mounts nothing until opened; no loading skeleton.
const FollowModal = dynamic(
  () => import('../../modals/FollowModal/FollowModal').then(mod => mod.FollowModal),
  {
    ssr: false,
  },
);

// The notify modal renders its own antd v6 modal shell; it mounts nothing
// until opened, so no loading skeleton.
const BrowserNotifyModal = dynamic(
  () =>
    import('../../modals/BrowserNotifyModal/BrowserNotifyModal').then(
      mod => mod.BrowserNotifyModal,
    ),
  {
    ssr: false,
  },
);

const OwncastPlayer = dynamic(
  () => import('../../video/OwncastPlayer/OwncastPlayer').then(mod => mod.OwncastPlayer),
  {
    ssr: false,
    loading: () => <Skeleton loading active paragraph={{ rows: 12 }} />,
  },
);

const ChatModal = dynamic(
  () => import('../../modals/ChatModal/ChatModal').then(mod => mod.ChatModal),
  {
    ssr: false,
  },
);

const ClipModal = dynamic(
  () => import('../../modals/ClipModal/ClipModal').then(mod => mod.ClipModal),
  { ssr: false },
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
  const { t } = useTranslation();
  const appState = useAtomValue(appStateAtom);
  const clientConfig = useAtomValue(clientConfigStateAtom);
  const chatState = useAtomValue(chatStateAtom);
  const currentUser = useAtomValue(currentUserAtom);
  const serverStatus = useAtomValue(serverStatusState);
  const [isMobile, setIsMobile] = useAtom(isMobileAtom);
  const messages = useAtomValue(visibleChatMessagesSelector);
  const online = useAtomValue(isOnlineSelector);
  const configLoaded = useAtomValue(isClientConfigLoadedAtom);
  const isChatAvailable = useAtomValue(isChatAvailableSelector);
  const isUserAuthenticated = useAtomValue(chatAuthenticatedAtom);

  const { viewerCount, lastConnectTime, lastDisconnectTime, streamTitle, streamId } = serverStatus;
  const {
    extraPageContent,
    name,
    summary,
    socialHandles,
    tags,
    externalActions,
    offlineMessage,
    chatDisabled,
    chatRequireAuthentication,
    federation,
    notifications,
    pluginTabs,
    autoplay,
    clips,
  } = clientConfig;
  const { enabled: fediverseEnabled, account: fediverseAccount, hideFollowersTab } = federation;
  const [showNotifyReminder, setShowNotifyReminder] = useState(false);
  const [showNotifyModal, setShowNotifyModal] = useState(false);
  const [showFollowModal, setShowFollowModal] = useState(false);
  const [showClipModal, setShowClipModal] = useState(false);
  const [activeClip, setActiveClip] = useState<{
    streamId: string;
    startSeconds: number;
    remainingSeconds: number;
  } | null>(null);
  const [pendingClipWindow, setPendingClipWindow] = useState<{
    streamId: string;
    startSeconds: number;
    endSeconds: number;
  } | null>(null);
  const [completedClipId, setCompletedClipId] = useState('');
  const [clipError, setClipError] = useState('');
  const [clipSaving, setClipSaving] = useState(false);
  const { browser: browserNotifications } = notifications;
  const { enabled: browserNotificationsEnabled } = browserNotifications;
  const { online: isStreamLive } = serverStatus;
  const [externalActionToDisplay, setExternalActionToDisplay] = useState<ExternalAction>(null);
  const [currentBrowserWindowUrl, setCurrentBrowserWindowUrl] = useState('');

  const [supportsBrowserNotifications, setSupportsBrowserNotifications] = useState(false);
  const supportFediverseFeatures = fediverseEnabled;
  // The Followers tab can be hidden independently of the rest of the
  // social features: federation stays on (follow button, go-live posts,
  // engagement) but the public followers list is not shown.
  const showFollowersTab = fediverseEnabled && !hideFollowersTab;
  const { servers: federatedServers } = useFederatedServers();

  // Capture the explicit media-time window locally. Nothing is persisted until
  // the viewer confirms the title prompt after ending the capture.
  const endClip = useCallback((capture: { streamId: string; startSeconds: number }) => {
    const endSeconds = getVideoPosition()?.playheadSeconds;
    if (endSeconds === null || endSeconds === undefined || endSeconds <= capture.startSeconds) {
      return;
    }
    setActiveClip(null);
    setPendingClipWindow({
      streamId: capture.streamId,
      startSeconds: capture.startSeconds,
      endSeconds,
    });
    setClipError('');
    setShowClipModal(true);
  }, []);

  const saveClip = useCallback(
    async (title: string) => {
      if (!pendingClipWindow) {
        return;
      }
      setClipSaving(true);
      setClipError('');
      try {
        const response = await fetch('/api/clip', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            streamId: pendingClipWindow.streamId,
            clipTitle: title,
            relativeStartTimeSeconds: pendingClipWindow.startSeconds,
            relativeEndTimeSeconds: pendingClipWindow.endSeconds,
          }),
        });
        const result = await response.json();
        if (!response.ok || result.success === false) {
          throw new Error(result.message || 'Unable to create clip');
        }
        setPendingClipWindow(null);
        setCompletedClipId(result.id || '');
      } catch (error) {
        setClipError(error instanceof Error ? error.message : 'Unable to create clip');
      } finally {
        setClipSaving(false);
      }
    },
    [pendingClipWindow],
  );

  const startClip = useCallback(() => {
    const playheadSeconds = getVideoPosition()?.playheadSeconds;
    if (playheadSeconds === null || playheadSeconds === undefined || !streamId) {
      return;
    }
    setActiveClip({
      streamId,
      startSeconds: Math.max(playheadSeconds - 10, 0),
      remainingSeconds: Math.max(1, (clips?.maxDurationSeconds ?? 1) - 10),
    });
  }, [clips?.maxDurationSeconds, streamId]);

  useEffect(() => {
    if (!activeClip) {
      return undefined;
    }

    const { startSeconds } = activeClip;
    const timer = window.setInterval(() => {
      setActiveClip(current => {
        if (!current || current.startSeconds !== startSeconds) {
          return current;
        }
        if (current.remainingSeconds <= 1) {
          endClip(current);
          return null;
        }
        return { ...current, remainingSeconds: current.remainingSeconds - 1 };
      });
    }, 1000);

    return () => window.clearInterval(timer);
  }, [activeClip, endClip]);

  const discardClip = useCallback(() => {
    setActiveClip(null);
    setPendingClipWindow(null);
    setShowClipModal(false);
    setClipError('');
  }, []);

  const clipAction = activeClip ? () => endClip(activeClip) : startClip;
  const clipsEnabled = !!clips?.enabled;
  const canShowClipButton = useAtomValue(canShowClipButtonSelector);
  const canClipNow = useAtomValue(canCreateClipSelector);

  const [showChatModal, setShowChatModal] = useState(false);

  const externalActionSelected = (action: ExternalAction) => {
    const { openExternally, url } = action;

    if (url) {
      // Plugin-contributed actions can use root-relative URLs (e.g.
      // "/plugins/<name>/") that the host validates and rewrites. Pass
      // window.location.origin as the base so URL() accepts both
      // absolute external URLs and same-origin plugin paths.
      const updatedUrl = new URL(url, window.location.origin);
      updatedUrl.searchParams.append('instance', currentBrowserWindowUrl);

      if (currentUser) {
        const { id, displayName } = currentUser;

        // Append url, username and userId to params so the link knows where we
        // came from and who we are. Display names are not unique, so userId
        // gives external actions a stable identifier for the chat user.
        updatedUrl.searchParams.append('username', displayName);
        updatedUrl.searchParams.append('userId', id);
      }
      const fullUrl = updatedUrl.toString();
      // Overwrite URL with the updated one that includes the params.
      const updatedAction = {
        ...action,
        url: fullUrl,
      };

      // apply openExternally only if we don't have an HTML embed
      if (openExternally) {
        window.open(fullUrl, '_blank');
      } else {
        setExternalActionToDisplay(updatedAction);
      }
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
    setSupportsBrowserNotifications(
      canPushNotificationsBeSupported() && browserNotificationsEnabled,
    );
  }, [browserNotificationsEnabled]);

  useEffect(() => {
    setCurrentBrowserWindowUrl(window.location.href);
  }, []);

  const showChat = isChatAvailable && !chatDisabled && chatState === ChatState.VISIBLE;

  // Determine if chat input should be enabled based on authentication requirements.
  // Moderators bypass the authentication requirement.
  const chatInputEnabled = !!(
    isChatAvailable &&
    (!chatRequireAuthentication || isUserAuthenticated || currentUser?.isModerator)
  );
  const chatInputDisabledMessage = chatRequireAuthentication
    ? t(Localization.Frontend.Chat.authenticateToChat)
    : t(Localization.Frontend.chatDisabled);

  return (
    <div className={styles.main}>
      <div className={styles.mainColumn}>
        {appState.appLoading && (
          <div
            className={classnames([styles.topSectionElement, styles.centerSpinner])}
            style={{ height: '30vh' }}
          >
            <Spin delay={2} size="large" tip="One moment..." />
          </div>
        )}
        <Row>
          {online && configLoaded && (
            <OwncastPlayer
              source="/hls/stream.m3u8"
              online={online}
              title={streamTitle || name}
              autoplay={autoplay}
              className={styles.topSectionElement}
            />
          )}
          {!online && !appState.appLoading && (
            <div id="offline-message" style={{ width: '100%' }}>
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
        </Row>
        <Row>
          {isStreamLive && (
            <Statusbar
              online={online}
              lastConnectTime={lastConnectTime}
              lastDisconnectTime={lastDisconnectTime}
              viewerCount={viewerCount}
              className={classnames(styles.topSectionElement, styles.statusBar)}
            />
          )}
        </Row>
        <Row>
          <ActionButtons
            supportFediverseFeatures={supportFediverseFeatures}
            supportsBrowserNotifications={supportsBrowserNotifications}
            showNotifyReminder={showNotifyReminder}
            setShowNotifyModal={setShowNotifyModal}
            disableNotifyReminderPopup={disableNotifyReminderPopup}
            externalActions={externalActions || []}
            setShowFollowModal={setShowFollowModal}
            externalActionSelected={externalActionSelected}
            showClipButton={canShowClipButton}
            canClip={canClipNow}
            clipActive={!!activeClip}
            clipRemainingSeconds={activeClip?.remainingSeconds || 0}
            onClipAction={clipAction}
          />
        </Row>

        {showNotifyModal && (
          <BrowserNotifyModal
            open={showNotifyModal}
            handleClose={() => disableNotifyReminderPopup()}
          />
        )}
        <Row>
          {!name && <Skeleton active loading style={{ marginLeft: '10vw', marginRight: '10vw' }} />}
          {isMobile ? (
            <MobileContent
              name={name}
              summary={summary}
              tags={tags}
              socialHandles={socialHandles}
              extraPageContent={extraPageContent}
              pluginTabs={pluginTabs}
              setShowFollowModal={setShowFollowModal}
              showFollowersTab={showFollowersTab}
              showClipsTab={clipsEnabled}
              online={online}
              federatedServers={federatedServers}
            />
          ) : (
            <div className={desktopStyles.bottomSectionContent}>
              <DesktopContent
                name={name}
                summary={summary}
                tags={tags}
                socialHandles={socialHandles}
                extraPageContent={extraPageContent}
                pluginTabs={pluginTabs}
                setShowFollowModal={setShowFollowModal}
                showFollowersTab={showFollowersTab}
                showClipsTab={clipsEnabled}
                federatedServers={federatedServers}
              />
            </div>
          )}
        </Row>
        <div style={{ flex: '1 1' }} />
        <Footer />
      </div>
      {showChat && !isMobile && currentUser && (
        <ChatContainer
          messages={messages}
          usernameToHighlight={currentUser.displayName}
          chatUserId={currentUser.id}
          isModerator={currentUser.isModerator}
          chatAvailable={isChatAvailable}
          showInput={!!currentUser}
          inputEnabled={chatInputEnabled}
          inputDisabledPlaceholder={chatInputDisabledMessage}
          desktop
        />
      )}
      {externalActionToDisplay && (
        <ExternalModal
          externalActionToDisplay={externalActionToDisplay}
          setExternalActionToDisplay={setExternalActionToDisplay}
        />
      )}
      {showFollowModal && (
        <FollowModal
          open={showFollowModal}
          account={fediverseAccount}
          name={name}
          handleClose={() => setShowFollowModal(false)}
        />
      )}
      {(showClipModal || !!completedClipId) && (
        <ClipModal
          open={showClipModal || !!completedClipId}
          completedClipId={completedClipId}
          onSaveTitle={saveClip}
          onDiscard={discardClip}
          saving={clipSaving}
          error={clipError}
          handleClose={() => {
            setShowClipModal(false);
            setPendingClipWindow(null);
            setCompletedClipId('');
            setClipError('');
          }}
        />
      )}
      {isMobile && showChatModal && chatState === ChatState.VISIBLE && (
        <ChatModal
          messages={messages}
          currentUser={currentUser}
          handleClose={() => setShowChatModal(false)}
          inputEnabled={chatInputEnabled}
          inputDisabledPlaceholder={chatInputDisabledMessage}
        />
      )}
      {isMobile && isChatAvailable && !chatDisabled && (
        <Button
          id="mobile-chat-button"
          type="primary"
          onClick={() => setShowChatModal(true)}
          className={styles.floatingMobileChatModalButton}
        >
          Chat <MessageFilled />
        </Button>
      )}
    </div>
  );
};

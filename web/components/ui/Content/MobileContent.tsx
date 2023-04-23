import React, { ComponentType, FC } from 'react';
import dynamic from 'next/dynamic';
import { Skeleton, TabsProps } from 'antd';
import { ErrorBoundary } from 'react-error-boundary';
import { SocialLink } from '../../../interfaces/social-link.model';
import styles from './Content.module.scss';
import mobileStyles from './MobileContent.module.scss';
import { CustomPageContent } from '../CustomPageContent/CustomPageContent';
import { ContentHeader } from '../../common/ContentHeader/ContentHeader';
import { ChatMessage } from '../../../interfaces/chat-message.model';
import { CurrentUser } from '../../../interfaces/current-user';
import { ActionButtonMenu } from '../../action-buttons/ActionButtonMenu/ActionButtonMenu';
import { ExternalAction } from '../../../interfaces/external-action';
import { ComponentError } from '../ComponentError/ComponentError';

export type MobileContentProps = {
  name: string;
  summary: string;
  tags: string[];
  socialHandles: SocialLink[];
  extraPageContent: string;
  notifyItemSelected: () => void;
  followItemSelected: () => void;
  setExternalActionToDisplay: (action: ExternalAction) => void;
  setShowNotifyPopup: (show: boolean) => void;
  setShowFollowModal: (show: boolean) => void;
  supportFediverseFeatures: boolean;
  messages: ChatMessage[];
  currentUser: CurrentUser;
  showChat: boolean;
  chatEnabled: boolean;
  actions: ExternalAction[];
  externalActionSelected: (action: ExternalAction) => void;
  supportsBrowserNotifications: boolean;
};

// lazy loaded components

const Tabs: ComponentType<TabsProps> = dynamic(() => import('antd').then(mod => mod.Tabs), {
  ssr: false,
});

const FollowerCollection = dynamic(
  () =>
    import('../followers/FollowerCollection/FollowerCollection').then(
      mod => mod.FollowerCollection,
    ),
  {
    ssr: false,
  },
);

const ChatContainer = dynamic(
  () => import('../../chat/ChatContainer/ChatContainer').then(mod => mod.ChatContainer),
  {
    ssr: false,
  },
);

type ChatContentProps = {
  showChat: boolean;
  chatEnabled: boolean;
  messages: ChatMessage[];
  currentUser: CurrentUser;
};

const ComponentErrorFallback = ({ error, resetErrorBoundary }) => (
  <ComponentError
    message={error}
    componentName="MobileContent"
    retryFunction={resetErrorBoundary}
  />
);

const ChatContent: FC<ChatContentProps> = ({ showChat, chatEnabled, messages, currentUser }) => {
  const { id, displayName } = currentUser;

  return showChat && !!currentUser ? (
    <ChatContainer
      messages={messages}
      usernameToHighlight={displayName}
      chatUserId={id}
      isModerator={false}
      chatAvailable={chatEnabled}
    />
  ) : (
    <Skeleton loading active paragraph={{ rows: 7 }} />
  );
};

const ActionButton = ({
  supportFediverseFeatures,
  supportsBrowserNotifications,
  actions,
  setExternalActionToDisplay,
  setShowNotifyPopup,
  setShowFollowModal,
}) => (
  <ActionButtonMenu
    className={styles.actionButtonMenu}
    showFollowItem={supportFediverseFeatures}
    showNotifyItem={supportsBrowserNotifications}
    actions={actions}
    externalActionSelected={setExternalActionToDisplay}
    notifyItemSelected={() => setShowNotifyPopup(true)}
    followItemSelected={() => setShowFollowModal(true)}
  />
);

export const MobileContent: FC<MobileContentProps> = ({
  name,
  summary,
  tags,
  socialHandles,
  extraPageContent,
  messages,
  currentUser,
  showChat,
  chatEnabled,
  actions,
  setExternalActionToDisplay,
  setShowNotifyPopup,
  setShowFollowModal,
  supportFediverseFeatures,
  supportsBrowserNotifications,
}) => {
  const aboutTabContent = (
    <>
      <ContentHeader name={name} summary={summary} tags={tags} links={socialHandles} logo="/logo" />
      {!!extraPageContent && (
        <div className={styles.bottomPageContentContainer}>
          <CustomPageContent content={extraPageContent} />
        </div>
      )}
    </>
  );
  const followersTabContent = (
    <div className={styles.bottomPageContentContainer}>
      <FollowerCollection name={name} onFollowButtonClick={() => setShowFollowModal(true)} />
    </div>
  );

  const items = [];
  if (showChat && currentUser) {
    items.push({
      label: 'Chat',
      key: '0',
      children: (
        <ChatContent
          showChat={showChat}
          chatEnabled={chatEnabled}
          messages={messages}
          currentUser={currentUser}
        />
      ),
    });
  }
  items.push({ label: 'About', key: '2', children: aboutTabContent });
  if (supportFediverseFeatures) {
    items.push({ label: 'Followers', key: '3', children: followersTabContent });
  }

  const replacementTabBar = (props, DefaultTabBar) => (
    <div className={styles.replacementBar}>
      <DefaultTabBar {...props} className={styles.defaultTabBar} />
      <ActionButton
        supportFediverseFeatures={supportFediverseFeatures}
        supportsBrowserNotifications={supportsBrowserNotifications}
        actions={actions}
        setExternalActionToDisplay={setExternalActionToDisplay}
        setShowNotifyPopup={setShowNotifyPopup}
        setShowFollowModal={setShowFollowModal}
      />
    </div>
  );

  return (
    <ErrorBoundary
      // eslint-disable-next-line react/no-unstable-nested-components
      fallbackRender={({ error, resetErrorBoundary }) => (
        <ComponentErrorFallback error={error} resetErrorBoundary={resetErrorBoundary} />
      )}
    >
      <div className={styles.lowerSectionMobile}>
        {items.length > 1 ? (
          <Tabs
            className={styles.tabs}
            defaultActiveKey="0"
            items={items}
            renderTabBar={replacementTabBar}
          />
        ) : (
          <>
            <div className={mobileStyles.noTabsActionMenuButton}>
              <ActionButton
                supportFediverseFeatures={supportFediverseFeatures}
                supportsBrowserNotifications={supportsBrowserNotifications}
                actions={actions}
                setExternalActionToDisplay={setExternalActionToDisplay}
                setShowNotifyPopup={setShowNotifyPopup}
                setShowFollowModal={setShowFollowModal}
              />
            </div>
            <div className={mobileStyles.noTabsAboutContent}>{aboutTabContent}</div>
          </>
        )}
      </div>
    </ErrorBoundary>
  );
};

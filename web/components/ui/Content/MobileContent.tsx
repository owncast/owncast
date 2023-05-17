import React, { ComponentType, FC } from 'react';
import dynamic from 'next/dynamic';
import { Skeleton, TabsProps } from 'antd';
import { ErrorBoundary } from 'react-error-boundary';
import classNames from 'classnames';
import { SocialLink } from '../../../interfaces/social-link.model';
import styles from './Content.module.scss';
import mobileStyles from './MobileContent.module.scss';
import { CustomPageContent } from '../CustomPageContent/CustomPageContent';
import { ContentHeader } from '../../common/ContentHeader/ContentHeader';
import { ChatMessage } from '../../../interfaces/chat-message.model';
import { CurrentUser } from '../../../interfaces/current-user';
import { ComponentError } from '../ComponentError/ComponentError';
import Footer from '../Footer/Footer';

export type MobileContentProps = {
  name: string;
  summary: string;
  tags: string[];
  socialHandles: SocialLink[];
  extraPageContent: string;
  setShowFollowModal: (show: boolean) => void;
  supportFediverseFeatures: boolean;
  messages: ChatMessage[];
  currentUser: CurrentUser;
  showChat: boolean;
  chatEnabled: boolean;
  online: boolean;
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
  setShowFollowModal,
  supportFediverseFeatures,
  online,
}) => {
  const aboutTabContent = (
    <>
      <ContentHeader name={name} summary={summary} tags={tags} links={socialHandles} logo="/logo" />
      {!!extraPageContent && (
        <div className={styles.bottomPageContentContainer}>
          <CustomPageContent content={extraPageContent} />
        </div>
      )}
      {online && (
        <div className={mobileStyles.footerContainer}>
          <Footer />
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

  return (
    <ErrorBoundary
      // eslint-disable-next-line react/no-unstable-nested-components
      fallbackRender={({ error, resetErrorBoundary }) => (
        <ComponentErrorFallback error={error} resetErrorBoundary={resetErrorBoundary} />
      )}
    >
      {items.length > 1 ? (
        <div className={classNames([styles.lowerSectionMobileTabbed, online && styles.online])}>
          <Tabs defaultActiveKey="0" items={items} />
        </div>
      ) : (
        <div className={styles.lowerSectionMobileNoTabs}>{aboutTabContent}</div>
      )}
    </ErrorBoundary>
  );
};

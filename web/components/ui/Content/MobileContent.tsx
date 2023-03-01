import React, { ComponentType, FC } from 'react';
import dynamic from 'next/dynamic';
import { Skeleton, TabsProps } from 'antd';
import { SocialLink } from '../../../interfaces/social-link.model';
import styles from './Content.module.scss';
import { CustomPageContent } from '../CustomPageContent/CustomPageContent';
import { ContentHeader } from '../../common/ContentHeader/ContentHeader';
import { ChatMessage } from '../../../interfaces/chat-message.model';
import { CurrentUser } from '../../../interfaces/current-user';
import { ActionButtonMenu } from '../../action-buttons/ActionButtonMenu/ActionButtonMenu';
import { ExternalAction } from '../../../interfaces/external-action';

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

export const MobileContent: FC<MobileContentProps> = ({
  name,
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
  if (showChat) {
    items.push({ label: 'Chat', key: '0', children: chatContent });
  }
  items.push({ label: 'About', key: '2', children: aboutTabContent });
  if (supportFediverseFeatures) {
    items.push({ label: 'Followers', key: '3', children: followersTabContent });
  }

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
      {items.length > 1 ? (
        <Tabs
          className={styles.tabs}
          defaultActiveKey="0"
          items={items}
          renderTabBar={replacementTabBar}
        />
      ) : (
        aboutTabContent
      )}
    </div>
  );
};

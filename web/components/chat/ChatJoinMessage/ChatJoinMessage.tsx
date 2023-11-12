import { FC } from 'react';
import dynamic from 'next/dynamic';
import { ModerationBadge } from '../ChatUserBadge/ModerationBadge';

import styles from './ChatJoinMessage.module.scss';

// Lazy loaded components

const UsergroupAddOutlined = dynamic(() => import('@ant-design/icons/UsergroupAddOutlined'), {
  ssr: false,
});

export type ChatJoinMessageProps = {
  isAuthorModerator: boolean;
  userColor: number;
  displayName: string;
};

export const ChatJoinMessage: FC<ChatJoinMessageProps> = ({
  isAuthorModerator,
  userColor,
  displayName,
}) => {
  const color = `var(--theme-color-users-${userColor})`;

  return (
    <div className={styles.root}>
      <span style={{ color }}>
        <span className={styles.icon}>
          <UsergroupAddOutlined />
        </span>
        <span className={styles.user}>{displayName}</span>
        {isAuthorModerator && (
          <span className={styles.moderatorBadge}>
            <ModerationBadge userColor={userColor} />
          </span>
        )}
      </span>
      joined the chat.
    </div>
  );
};

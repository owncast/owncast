import { FC } from 'react';
import dynamic from 'next/dynamic';
import styles from './ChatJoinMessage.module.scss';
import { ModerationBadge } from '../ChatUserBadge/ModerationBadge';

// Lazy loaded components

const TeamOutlined = dynamic(() => import('@ant-design/icons/TeamOutlined'), {
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
        <span style={{ padding: '0 10px' }}>
          <TeamOutlined />
        </span>
        <span style={{ fontWeight: 'bold' }}>{displayName}</span>
        {isAuthorModerator && (
          <span>
            <ModerationBadge userColor={userColor} />
          </span>
        )}
      </span>{' '}
      joined the chat.
    </div>
  );
};

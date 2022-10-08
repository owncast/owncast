import { FC } from 'react';
import { TeamOutlined } from '@ant-design/icons';
import { ChatUserBadge } from '~/components/chat/ChatUserBadge/ChatUserBadge';
import styles from './ChatJoinMessage.module.scss';

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
            <ChatUserBadge badge="mod" userColor={userColor} />
          </span>
        )}
      </span>{' '}
      joined the chat.
    </div>
  );
};

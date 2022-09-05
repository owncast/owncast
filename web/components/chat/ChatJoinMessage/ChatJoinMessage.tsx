import { FC } from 'react';
import s from './ChatJoinMessage.module.scss';
import { ChatUserBadge } from '../ChatUserBadge/ChatUserBadge';

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
  const color = `var(--theme-user-colors-${userColor})`;

  return (
    <div className={s.join}>
      <span style={{ color }}>
        {displayName}
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

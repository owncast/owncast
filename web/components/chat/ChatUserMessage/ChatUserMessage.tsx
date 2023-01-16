/* eslint-disable react/no-danger */
import { FC, ReactNode, useEffect, useState } from 'react';
import cn from 'classnames';
import { Tooltip } from 'antd';
import { useRecoilValue } from 'recoil';
import dynamic from 'next/dynamic';
import { decodeHTML } from 'entities';
import styles from './ChatUserMessage.module.scss';
import { formatTimestamp } from './messageFmt';
import { ChatMessage } from '../../../interfaces/chat-message.model';
import { ChatUserBadge } from '../ChatUserBadge/ChatUserBadge';
import { accessTokenAtom } from '../../stores/ClientConfigStore';
import { User } from '../../../interfaces/user.model';

// Lazy loaded components

const LinkOutlined = dynamic(() => import('@ant-design/icons/LinkOutlined'), {
  ssr: false,
});

const ChatModerationActionMenu = dynamic(
  () =>
    import('../ChatModerationActionMenu/ChatModerationActionMenu').then(
      mod => mod.ChatModerationActionMenu,
    ),
  {
    ssr: false,
  },
);

const Highlight = dynamic(() => import('react-highlighter-ts').then(mod => mod.Highlight), {
  ssr: false,
});

export type ChatUserMessageProps = {
  message: ChatMessage;
  showModeratorMenu: boolean;
  highlightString: string;
  sentBySelf: boolean;
  sameUserAsLast: boolean;
  isAuthorModerator: boolean;
  isAuthorAuthenticated: boolean;
};

export type UserTooltipProps = {
  user: User;
  children: ReactNode;
};

const UserTooltip: FC<UserTooltipProps> = ({ children, user }) => {
  const { displayName, createdAt } = user;
  const content = `${displayName} first joined ${formatTimestamp(createdAt)}`;

  return (
    <Tooltip title={content} placement="topLeft" mouseEnterDelay={1}>
      {children}
    </Tooltip>
  );
};

export const ChatUserMessage: FC<ChatUserMessageProps> = ({
  message,
  highlightString,
  showModeratorMenu,
  sentBySelf, // Move the border to the right and render a background
  sameUserAsLast,
  isAuthorModerator,
  isAuthorAuthenticated,
}) => {
  const { id: messageId, body, user, timestamp } = message;
  const { id: userId, displayName, displayColor } = user;
  const accessToken = useRecoilValue<string>(accessTokenAtom);

  const color = `var(--theme-color-users-${displayColor})`;
  const formattedTimestamp = `Sent ${formatTimestamp(timestamp)}`;
  const [formattedMessage, setFormattedMessage] = useState<string>(body);

  const badgeNodes = [];
  if (isAuthorModerator) {
    badgeNodes.push(<ChatUserBadge key="mod" badge="mod" userColor={displayColor} />);
  }
  if (isAuthorAuthenticated) {
    badgeNodes.push(
      <ChatUserBadge
        key="auth"
        badge={<LinkOutlined title="authenticated" />}
        userColor={displayColor}
      />,
    );
  }

  useEffect(() => {
    setFormattedMessage(decodeHTML(body));
  }, [message]);

  return (
    <div
      className={cn(
        styles.messagePadding,
        sameUserAsLast && styles.messagePaddingCollapsed,
        'chat-message_user',
      )}
    >
      <div
        className={cn(styles.root, {
          [styles.ownMessage]: sentBySelf,
        })}
        style={{ borderColor: color }}
      >
        {!sameUserAsLast && (
          <UserTooltip user={user}>
            <div className={styles.user} style={{ color }}>
              <span className={styles.userName}>{displayName}</span>
              <span>{badgeNodes}</span>
            </div>
          </UserTooltip>
        )}
        <Tooltip title={formattedTimestamp} mouseEnterDelay={1}>
          <Highlight search={highlightString}>
            <div
              className={styles.message}
              dangerouslySetInnerHTML={{ __html: formattedMessage }}
            />
          </Highlight>
        </Tooltip>

        {showModeratorMenu && (
          <div className={styles.modMenuWrapper}>
            <ChatModerationActionMenu
              messageID={messageId}
              accessToken={accessToken}
              userID={userId}
              userDisplayName={displayName}
            />
          </div>
        )}
        <div className={styles.background} style={{ color }} />
      </div>
    </div>
  );
};

/* eslint-disable react/no-danger */
import { useEffect, useState } from 'react';
import { Highlight } from 'react-highlighter-ts';
import he from 'he';
import cn from 'classnames';
import s from './ChatUserMessage.module.scss';
import { formatTimestamp } from './messageFmt';
import { ChatMessage } from '../../../interfaces/chat-message.model';
import ChatModerationActionMenu from '../ChatModerationActionMenu/ChatModerationActionMenu';
import ChatUserBadge from './ChatUserBadge';

interface Props {
  message: ChatMessage;
  showModeratorMenu: boolean;
  highlightString: string;
  sentBySelf: boolean;
  sameUserAsLast: boolean;
  isAuthorModerator: boolean;
  isAuthorAuthenticated: boolean;
}

export default function ChatUserMessage({
  message,
  highlightString,
  showModeratorMenu,
  sentBySelf, // Move the border to the right and render a background
  sameUserAsLast,
  isAuthorModerator,
  isAuthorAuthenticated,
}: Props) {
  const { id: messageId, body, user, timestamp } = message;
  const { id: userId, displayName, displayColor } = user;

  const color = `var(--theme-user-colors-${displayColor})`;
  const formattedTimestamp = `Sent at ${formatTimestamp(timestamp)}`;
  const [formattedMessage, setFormattedMessage] = useState<string>(body);

  const badgeStrings = [isAuthorModerator && 'mod', isAuthorAuthenticated && 'auth'];
  const badges = badgeStrings
    .filter(badge => !!badge)
    .map(badge => <ChatUserBadge key={badge} badge={badge} userColor={displayColor} />);

  useEffect(() => {
    setFormattedMessage(he.decode(body));
  }, [message]);

  return (
    <div style={{ padding: 3.5 }}>
      <div
        className={cn(s.root, {
          [s.ownMessage]: sentBySelf,
        })}
        style={{ borderColor: color }}
        title={formattedTimestamp}
      >
        {!sameUserAsLast && (
          <div className={s.user} style={{ color }}>
            <span className={s.userName}>{displayName}</span>
            <span>{badges}</span>
          </div>
        )}
        <Highlight search={highlightString}>
          <div className={s.message}>{formattedMessage}</div>
        </Highlight>
        {showModeratorMenu && (
          <div className={s.modMenuWrapper}>
            <ChatModerationActionMenu
              messageID={messageId}
              accessToken=""
              userID={userId}
              userDisplayName={displayName}
            />
          </div>
        )}
        <div className={s.customBorder} style={{ color }} />
        <div className={s.background} style={{ color }} />
      </div>
    </div>
  );
}

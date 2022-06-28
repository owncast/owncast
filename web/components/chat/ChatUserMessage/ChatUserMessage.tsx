/* eslint-disable react/no-danger */
import { useEffect, useState } from 'react';
import { Highlight } from 'react-highlighter-ts';
import he from 'he';
import cn from 'classnames';
import { ChatMessage } from '../../../interfaces/chat-message.model';
import { formatTimestamp } from './messageFmt';
import s from './ChatUserMessage.module.scss';

interface Props {
  message: ChatMessage;
  showModeratorMenu: boolean;
  highlightString: string;
  renderAsPersonallySent: boolean;
}

export default function ChatUserMessage({
  message,
  highlightString,
  showModeratorMenu,
  renderAsPersonallySent, // Move the border to the right and render a background
}: Props) {
  const { body, user, timestamp } = message;
  const { displayName, displayColor } = user;
  const color = `hsl(${displayColor}, 100%, 65%)`;
  const bgColor = `hsl(${displayColor}, 20%, 25%)`;
  const formattedTimestamp = `Sent at ${formatTimestamp(timestamp)}`;
  const [formattedMessage, setFormattedMessage] = useState<string>(body);

  useEffect(() => {
    setFormattedMessage(he.decode(body));
  }, [message]);

  return (
    <div
      className={cn(s.root, {
        [s.ownMessage]: renderAsPersonallySent,
      })}
      style={{ borderColor: color, backgroundColor: bgColor }}
      title={formattedTimestamp}
    >
      <div className={s.user} style={{ color }}>
        {displayName}
      </div>
      <Highlight search={highlightString}>
        <div className={s.message}>{formattedMessage}</div>
      </Highlight>
      {showModeratorMenu && <div>Moderator menu</div>}
    </div>
  );
}

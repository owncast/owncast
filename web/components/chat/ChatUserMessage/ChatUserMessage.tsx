import { useEffect, useState } from 'react';
import { ChatMessage } from '../../../interfaces/chat-message.model';
import { formatTimestamp, formatMessageText } from './messageFmt';
import s from './ChatUserMessage.module.scss';

interface Props {
  message: ChatMessage;
  showModeratorMenu: boolean;
}

export default function ChatUserMessage({ message, showModeratorMenu }: Props) {
  const { body, user, timestamp } = message;
  const { displayName, displayColor } = user;
  const color = `hsl(${displayColor}, 100%, 70%)`;
  const formattedTimestamp = `Sent at ${formatTimestamp(timestamp)}`;
  const [formattedMessage, setFormattedMessage] = useState<string>(body);

  useEffect(() => {
    setFormattedMessage(formatMessageText(body));
  }, [message]);

  return (
    <div className={s.root} style={{ borderColor: color }} title={formattedTimestamp}>
      <div className={s.user} style={{ color }}>
        {displayName}
      </div>
      <div className={s.message} dangerouslySetInnerHTML={{ __html: formattedMessage }} />
      {showModeratorMenu && <div>Moderator menu</div>}
    </div>
  );
}

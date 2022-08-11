/* eslint-disable react/no-danger */
import { Highlight } from 'react-highlighter-ts';
import { ChatMessage } from '../../../interfaces/chat-message.model';
import s from './ChatSystemMessage.module.scss';

interface Props {
  message: ChatMessage;
  highlightString: string;
}
// eslint-disable-next-line @typescript-eslint/no-unused-vars
export default function ChatSystemMessage({ message, highlightString }: Props) {
  const { body, user } = message;
  const { displayName } = user;

  return (
    <div className={s.chatSystemMessage}>
      <div className={s.user}>
        <span className={s.userName}>{displayName}</span>
      </div>
      <Highlight search={highlightString}>
        <div className={s.message} dangerouslySetInnerHTML={{ __html: body }} />
      </Highlight>
    </div>
  );
}

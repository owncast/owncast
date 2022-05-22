import { ChatMessage } from '../../../interfaces/chat-message.model';
import s from './ChatUserMessage.module.scss';

interface Props {
  message: ChatMessage;
  showModeratorMenu: boolean;
}

export default function ChatUserMessage(props: Props) {
  const { message, showModeratorMenu } = props;
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  const { body, user, timestamp } = message;
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  const { displayName, displayColor } = user;

  const color = `hsl(${displayColor}, 100%, 70%)`;

  return (
    <div className={s.root} style={{ borderColor: color }}>
      <div className={s.user} style={{ color }}>
        {displayName}
      </div>
      <div className={s.message}>{body}</div>
      {showModeratorMenu && <div>Moderator menu</div>}
    </div>
  );
}

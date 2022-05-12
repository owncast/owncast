import { ChatMessage } from '../../interfaces/chat-message.model';

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

  // TODO: Convert displayColor (a hue) to a usable color.

  return (
    <div>
      <div>{displayName}</div>
      <div>{body}</div>
      {showModeratorMenu && <div>Moderator menu</div>}
    </div>
  );
}

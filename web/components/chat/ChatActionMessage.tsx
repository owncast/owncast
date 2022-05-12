import { ChatMessage } from '../../interfaces/chat-message.model';

interface Props {
  // eslint-disable-next-line react/no-unused-prop-types
  message: ChatMessage;
}

// eslint-disable-next-line @typescript-eslint/no-unused-vars
export default function ChatSystemMessage(props: Props) {
  return <div>Component goes here</div>;
}

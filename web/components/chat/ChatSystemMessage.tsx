import { ChatMessage } from '../../interfaces/chat-message.model';

interface Props {
  message: ChatMessage;
}

export default function ChatSystemMessage(props: Props) {
  return <div>Component goes here</div>;
}

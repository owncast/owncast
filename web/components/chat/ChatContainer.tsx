import { ChatMessage } from '../../interfaces/chat-message.model';

interface Props {
  messages: ChatMessage[];
}

export default function ChatContainer(props: Props) {
  return <div>Chat container goes here</div>;
}

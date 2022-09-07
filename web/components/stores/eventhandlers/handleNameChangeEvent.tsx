import { ChatMessage } from '../../../interfaces/chat-message.model';
import { ChatEvent } from '../../../interfaces/socket-events';

export function handleNameChangeEvent(
  message: ChatEvent,
  messages: ChatMessage[],
  setChatMessages,
) {
  const updatedMessages = [...messages, message];
  setChatMessages(updatedMessages);
}
export default handleNameChangeEvent;

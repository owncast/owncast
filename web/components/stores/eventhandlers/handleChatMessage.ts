import {
  ChatEvent,
  SocketEvent,
} from '../../../interfaces/socket-events';

export default function handleChatMessage(message: ChatEvent) {
  console.log('chat message', message);
}

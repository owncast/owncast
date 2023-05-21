import { ChatEvent } from '../../../interfaces/socket-events';

export function handleNameChangeEvent(message: ChatEvent, setChatMessages) {
  setChatMessages(currentState => [...currentState, message]);
}

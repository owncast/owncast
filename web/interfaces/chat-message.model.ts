import { SocketEvent } from './socket-events';
import { User } from './user.model';

export interface ChatMessage extends SocketEvent {
  user: User;
  body: string;
}

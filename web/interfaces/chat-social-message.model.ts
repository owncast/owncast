import { SocketEvent } from './socket-events';

export interface ChatSocialMessage extends SocketEvent {
  title: string;
  body: string;
  image: string;
  link: string;
}

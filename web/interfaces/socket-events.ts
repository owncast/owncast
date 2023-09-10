import { User } from './user.model';

export enum MessageType {
  CHAT = 'CHAT',
  PING = 'PING',
  NAME_CHANGE = 'NAME_CHANGE',
  COLOR_CHANGE = 'COLOR_CHANGE',
  PONG = 'PONG',
  SYSTEM = 'SYSTEM',
  USER_JOINED = 'USER_JOINED',
  USER_PARTED = 'USER_PARTED',
  CHAT_ACTION = 'CHAT_ACTION',
  FEDIVERSE_ENGAGEMENT_FOLLOW = 'FEDIVERSE_ENGAGEMENT_FOLLOW',
  FEDIVERSE_ENGAGEMENT_LIKE = 'FEDIVERSE_ENGAGEMENT_LIKE',
  FEDIVERSE_ENGAGEMENT_REPOST = 'FEDIVERSE_ENGAGEMENT_REPOST',
  CONNECTED_USER_INFO = 'CONNECTED_USER_INFO',
  ERROR_USER_DISABLED = 'ERROR_USER_DISABLED',
  ERROR_NEEDS_REGISTRATION = 'ERROR_NEEDS_REGISTRATION',
  ERROR_MAX_CONNECTIONS_EXCEEDED = 'ERROR_MAX_CONNECTIONS_EXCEEDED',
  VISIBILITY_UPDATE = 'VISIBILITY-UPDATE',
}

export interface SocketEvent {
  id: string;
  timestamp: Date;
  type: MessageType;
}

export interface ConnectedClientInfoEvent extends SocketEvent {
  user: User;
}
export class ChatEvent implements SocketEvent {
  constructor(message) {
    this.id = message.id;
    this.timestamp = message.timestamp;
    this.type = message.type;
    this.body = message.body;
    if (message.user) {
      this.user = new User(message.user);
    }
  }

  timestamp: Date;

  type: MessageType;

  id: string;

  user: User;

  body: string;
}

export interface NameChangeEvent extends SocketEvent {
  user: User;
  oldName: string;
}

export interface MessageVisibilityEvent extends SocketEvent {
  visible: boolean;
  ids: string[];
}

export interface FediverseEvent extends SocketEvent {
  title: string;
  image: string;
  link: string;
  body: string;
}

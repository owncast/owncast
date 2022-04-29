import { User } from './user';

export interface ChatMessage {
  id: string;
  type: string;
  timestamp: Date;
  user: User;
  body: string;
}

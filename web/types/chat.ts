export interface MessageType {
  user: User;
  body: string;
  id: string;
  key: string;
  name: string;
  timestamp: string;
  type: string;
  visible: boolean;
}

export interface User {
  id: string;
  displayName: string;
  createdAt: Date;
  disabledAt: Date;
  previousNames: [string];
  nameChangedAt: Date;
}

export interface UsernameHistory {
  displayName: string;
  changedAt: Date;
}

export interface UserConnectionInfo {
  connectedAt: Date;
  messageCount: number;
  userAgent: string;
}

export interface Client extends UserConnectionInfo {
  user: User;
}

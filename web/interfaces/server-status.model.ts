export interface ServerStatus {
  online: boolean;
  viewerCount: number;
  lastConnectTime?: Date;
  lastDisconnectTime?: Date;
  versionNumber?: string;
  streamTitle?: string;
  serverTime: Date;
}

export function makeEmptyServerStatus(): ServerStatus {
  return {
    online: false,
    viewerCount: 0,
    serverTime: new Date(),
  };
}

export interface ServerStatus {
  online: boolean;
  viewerCount: number;
  lastConnectTime?: Date;
  lastDisconnectTime?: Date;
  versionNumber?: string;
  streamTitle?: string;
  serverTime: Date;
  // streamId identifies the currently-live broadcast, needed to clip what is
  // happening right now. Absent while offline.
  streamId?: string;
}

export function makeEmptyServerStatus(): ServerStatus {
  return {
    online: false,
    viewerCount: 0,
    serverTime: new Date(),
  };
}

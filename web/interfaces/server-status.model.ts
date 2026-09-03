export interface ScheduledEventStatus {
  id: string;
  name: string;
  description: string;
  startTime: string;
  durationMinutes: number;
  chatOpen: boolean;
}

export interface ServerStatus {
  online: boolean;
  viewerCount: number;
  lastConnectTime?: Date;
  lastDisconnectTime?: Date;
  versionNumber?: string;
  streamTitle?: string;
  serverTime: Date;
  scheduledEvent?: ScheduledEventStatus;
}

export function makeEmptyServerStatus(): ServerStatus {
  return {
    online: false,
    viewerCount: 0,
    serverTime: new Date(),
  };
}

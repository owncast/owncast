export type ScheduledEventStatus = 'scheduled' | 'cancelled';

export interface ScheduledEvent {
  id: string;
  seriesId?: string;
  name: string;
  description: string;
  reminderMessage?: string;
  startTime: string;
  durationMinutes: number;
  timezone: string;
  status: ScheduledEventStatus;
}

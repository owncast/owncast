import { Alert, Empty, Modal, Spin, Tag, Typography } from 'antd';
import type { DatesSetArg } from '@fullcalendar/core';
import { FC, useCallback, useRef, useState } from 'react';
import { useTranslation } from 'next-export-i18n';
import type { ScheduledEvent } from '../../../interfaces/scheduled-event.model';
import { PUBLIC_SCHEDULE, PUBLIC_SCHEDULE_ICS, fetchData } from '../../../utils/apis';
import { Localization } from '../../../types/localization';
import { ScheduleCalendar } from './ScheduleCalendar';
import type { ScheduleCalendarView } from './ScheduleCalendar';
import styles from './ScheduleTab.module.scss';

const formatEventDateTime = (startTime: string) =>
  new Intl.DateTimeFormat(undefined, {
    weekday: 'long',
    month: 'long',
    day: 'numeric',
    year: 'numeric',
    hour: 'numeric',
    minute: '2-digit',
  }).format(new Date(startTime));

const getInitialView = (): ScheduleCalendarView =>
  typeof window !== 'undefined' && window.innerWidth <= 768 ? 'listMonth' : 'dayGridMonth';

export const ScheduleTab: FC = () => {
  const { t } = useTranslation();
  const [events, setEvents] = useState<ScheduledEvent[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [selectedEvent, setSelectedEvent] = useState<ScheduledEvent | null>(null);
  const requestId = useRef(0);
  const lastRange = useRef<DatesSetArg | null>(null);

  const loadEvents = useCallback(async (range: DatesSetArg) => {
    lastRange.current = range;
    const currentRequest = ++requestId.current;
    setLoading(true);
    setError(null);

    try {
      const url = new URL(PUBLIC_SCHEDULE, window.location.origin);
      url.searchParams.set('from', range.start.toISOString());
      url.searchParams.set('to', range.end.toISOString());
      const result = await fetchData<ScheduledEvent[]>(url.toString(), { auth: false });
      if (currentRequest === requestId.current) {
        setEvents(result);
      }
    } catch (loadError) {
      if (currentRequest === requestId.current) {
        setError(`${loadError}`);
      }
    } finally {
      if (currentRequest === requestId.current) {
        setLoading(false);
      }
    }
  }, []);

  const getText = (key: string, fallback: string) => {
    const translated = t(key);
    return translated === key ? fallback : translated;
  };
  const noEventsDescription = getText(
    Localization.Frontend.Schedule.noEvents,
    'No streams scheduled in this range',
  );
  const loadingText = getText(Localization.Frontend.Schedule.loading, 'Loading schedule...');
  const errorTitle = getText(Localization.Frontend.Schedule.error, 'Unable to load the schedule');
  const retryText = getText(Localization.Frontend.Schedule.retry, 'Try again');
  const cancelledText = getText(Localization.Frontend.Schedule.cancelled, 'Cancelled');
  const timezoneText = getText(Localization.Frontend.Schedule.timezone, 'Schedule timezone');
  const durationText = getText(Localization.Frontend.Schedule.duration, '{{minutes}} minutes');

  return (
    <div className={styles.container}>
      {error && (
        <Alert
          className={styles.alert}
          type="error"
          showIcon
          title={errorTitle}
          description={
            <button
              type="button"
              className={styles.retry}
              onClick={() => {
                if (lastRange.current) {
                  loadEvents(lastRange.current);
                }
              }}
            >
              {retryText}
            </button>
          }
        />
      )}
      {loading && (
        <div className={styles.loading} role="status" aria-live="polite">
          <Spin description={loadingText} />
        </div>
      )}
      {!loading && !error && events.length === 0 && (
        <Empty
          className={styles.empty}
          image={Empty.PRESENTED_IMAGE_SIMPLE}
          description={noEventsDescription}
        />
      )}
      <ScheduleCalendar
        events={events}
        initialView={getInitialView()}
        onDatesSet={loadEvents}
        onEventClick={setSelectedEvent}
        calendarAction={{
          text: getText(Localization.Frontend.Schedule.addToCalendar, 'Add to calendar'),
          onClick: () => window.open(PUBLIC_SCHEDULE_ICS, '_blank', 'noopener'),
        }}
      />
      <Modal
        open={!!selectedEvent}
        title={selectedEvent?.name}
        footer={null}
        onCancel={() => setSelectedEvent(null)}
      >
        {selectedEvent && (
          <div className={styles.eventDetails}>
            {selectedEvent.status === 'cancelled' && <Tag color="red">{cancelledText}</Tag>}
            <Typography.Paragraph>
              <strong>{formatEventDateTime(selectedEvent.startTime)}</strong>
              <br />
              {durationText.replace('{{minutes}}', `${selectedEvent.durationMinutes}`)}
            </Typography.Paragraph>
            {selectedEvent.description && (
              <Typography.Paragraph>{selectedEvent.description}</Typography.Paragraph>
            )}
            {selectedEvent.timezone !== Intl.DateTimeFormat().resolvedOptions().timeZone && (
              <Typography.Text type="secondary">
                {timezoneText}: {selectedEvent.timezone}
              </Typography.Text>
            )}
          </div>
        )}
      </Modal>
    </div>
  );
};

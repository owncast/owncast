import type { DatesSetArg, EventClickArg, EventInput } from '@fullcalendar/core';
import dayGridPlugin from '@fullcalendar/daygrid';
import FullCalendar from '@fullcalendar/react';
import listPlugin from '@fullcalendar/list';
import timeGridPlugin from '@fullcalendar/timegrid';
import type { FC, MouseEvent, ReactNode } from 'react';
import { useTranslation } from 'next-export-i18n';
import type { ScheduledEvent } from '../../../interfaces/scheduled-event.model';
import { Localization } from '../../../types/localization';
import styles from './ScheduleTab.module.scss';

export type ScheduleCalendarView = 'dayGridMonth' | 'timeGridWeek' | 'listMonth';

export interface ScheduleCalendarProps {
  events: ScheduledEvent[];
  initialView?: ScheduleCalendarView;
  initialDate?: string;
  headerAction?: {
    text: string;
    onClick: () => void;
  };
  calendarAction?: {
    text: string;
    onClick: () => void;
  };
  onDatesSet?: (range: DatesSetArg) => void;
  onEventClick?: (event: ScheduledEvent) => void;
  renderEventActions?: (event: ScheduledEvent) => ReactNode;
}

export const isEventLive = (event: ScheduledEvent) => {
  if (event.status === 'cancelled') {
    return false;
  }
  const start = new Date(event.startTime).getTime();
  const end = start + event.durationMinutes * 60 * 1000;
  const now = Date.now();
  return now >= start && now < end;
};

const toCalendarEvent = (event: ScheduledEvent): EventInput => ({
  id: event.id,
  title: event.name,
  start: event.startTime,
  end: new Date(new Date(event.startTime).getTime() + event.durationMinutes * 60 * 1000),
  classNames: [
    ...(event.status === 'cancelled' ? [styles.cancelledEvent] : []),
    ...(isEventLive(event) ? [styles.liveEvent] : []),
  ],
  extendedProps: { scheduleEvent: event },
});

export const ScheduleCalendar: FC<ScheduleCalendarProps> = ({
  events,
  initialView = 'dayGridMonth',
  initialDate,
  headerAction,
  calendarAction,
  onDatesSet,
  onEventClick,
  renderEventActions,
}) => {
  const { t } = useTranslation();
  const getText = (key: string, fallback: string) => {
    const translated = t(key);
    return translated === key ? fallback : translated;
  };
  const noEventsText = getText(
    Localization.Frontend.Schedule.noEvents,
    'No streams scheduled in this range',
  );
  const cancelledText = getText(Localization.Frontend.Schedule.cancelled, 'Cancelled');
  const liveText = getText(Localization.Frontend.Schedule.liveNow, 'Live now');

  const handleEventClick = (info: EventClickArg) => {
    onEventClick?.(info.event.extendedProps.scheduleEvent as ScheduledEvent);
  };

  const stopEventClick = (event: MouseEvent<HTMLDivElement>) => {
    event.stopPropagation();
  };

  return (
    <div className={styles.calendar} data-testid="schedule-calendar">
      <FullCalendar
        plugins={[dayGridPlugin, timeGridPlugin, listPlugin]}
        initialView={initialView}
        initialDate={initialDate}
        customButtons={
          headerAction || calendarAction
            ? {
                ...(calendarAction
                  ? {
                      calendarAction: {
                        text: calendarAction.text,
                        click: calendarAction.onClick,
                      },
                    }
                  : {}),
                ...(headerAction
                  ? {
                      scheduleAction: {
                        text: headerAction.text,
                        click: headerAction.onClick,
                      },
                    }
                  : {}),
              }
            : undefined
        }
        headerToolbar={{
          left: `prev,next today${calendarAction ? ' calendarAction' : ''}${
            headerAction ? ' scheduleAction' : ''
          }`,
          center: 'title',
          right: 'dayGridMonth,timeGridWeek,listMonth',
        }}
        buttonText={{
          today: getText(Localization.Frontend.Schedule.today, 'Today'),
          month: getText(Localization.Frontend.Schedule.month, 'Month'),
          week: getText(Localization.Frontend.Schedule.week, 'Week'),
          list: getText(Localization.Frontend.Schedule.list, 'List'),
        }}
        events={events.map(toCalendarEvent)}
        datesSet={onDatesSet}
        eventClick={handleEventClick}
        eventContent={info => {
          const scheduledEvent = info.event.extendedProps.scheduleEvent as ScheduledEvent;
          return (
            <div className={styles.eventContent}>
              <div className={styles.eventLabel}>
                {info.timeText && <span>{info.timeText} </span>}
                <span>{scheduledEvent.name}</span>
                {scheduledEvent.status === 'cancelled' && (
                  <span className={styles.eventStatus}>{cancelledText}</span>
                )}
                {isEventLive(scheduledEvent) && (
                  <span className={styles.eventStatus}>{liveText}</span>
                )}
              </div>
              {renderEventActions && info.view.type === 'listMonth' && (
                <div className={styles.eventActions} onClick={stopEventClick}>
                  {renderEventActions(scheduledEvent)}
                </div>
              )}
            </div>
          );
        }}
        eventDisplay="block"
        dayMaxEvents
        height="auto"
        slotDuration="01:00:00"
        noEventsContent={noEventsText}
      />
    </div>
  );
};

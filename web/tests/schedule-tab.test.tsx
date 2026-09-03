import type { DatesSetArg, EventClickArg, EventInput } from '@fullcalendar/core';
import type { ReactNode } from 'react';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom';
import { ScheduleCalendar } from '../components/ui/ScheduleTab/ScheduleCalendar';
import { ScheduleTab } from '../components/ui/ScheduleTab/ScheduleTab';
import type { ScheduledEvent } from '../interfaces/scheduled-event.model';
type MockCalendarProps = {
  events: EventInput[];
  datesSet?: (range: DatesSetArg) => void;
  eventClick: (info: EventClickArg) => void;
  initialView: string;
  customButtons?: Record<string, { text: string; click: () => void }>;
  eventContent: (info: {
    event: { extendedProps: Record<string, unknown> };
    timeText: string;
    view: { type: string };
  }) => ReactNode;
};

jest.mock('@fullcalendar/daygrid', () => ({ __esModule: true, default: {} }));
jest.mock('@fullcalendar/list', () => ({ __esModule: true, default: {} }));
jest.mock('@fullcalendar/timegrid', () => ({ __esModule: true, default: {} }));
jest.mock('@fullcalendar/react', () => {
  // eslint-disable-next-line @typescript-eslint/no-require-imports
  const React = require('react');
  const MockCalendar = ({
    events,
    datesSet,
    eventClick,
    eventContent,
    initialView,
    customButtons,
  }: MockCalendarProps) => {
    React.useEffect(() => {
      datesSet?.({
        start: new Date('2030-09-01T00:00:00Z'),
        end: new Date('2030-10-01T00:00:00Z'),
        startStr: '2030-09-01T00:00:00.000Z',
        endStr: '2030-10-01T00:00:00.000Z',
        view: {},
      } as DatesSetArg);
    }, [datesSet]);
    return React.createElement(
      'div',
      null,
      customButtons?.calendarAction &&
        React.createElement(
          'button',
          { type: 'button', onClick: customButtons.calendarAction.click },
          customButtons.calendarAction.text,
        ),
      customButtons?.scheduleAction &&
        React.createElement(
          'button',
          { type: 'button', onClick: customButtons.scheduleAction.click },
          customButtons.scheduleAction.text,
        ),
      events.map(event =>
        React.createElement(
          'div',
          {
            key: event.id,
            onClick: () =>
              eventClick({
                event: { extendedProps: event.extendedProps },
              } as unknown as EventClickArg),
          },
          eventContent({
            event: { extendedProps: event.extendedProps },
            timeText: '6:00 PM',
            view: { type: initialView },
          }),
        ),
      ),
    );
  };
  return { __esModule: true, default: MockCalendar };
});
jest.mock('next/router', () => ({
  useRouter: () => ({
    query: {},
    pathname: '/',
    asPath: '/',
    push: jest.fn(),
    replace: jest.fn(),
  }),
}));

const scheduledEvent = {
  id: 'scheduled-show',
  name: 'Scheduled show',
  description: 'A useful description.',
  startTime: '2030-09-08T18:00:00-04:00',
  durationMinutes: 90,
  timezone: 'America/New_York',
  status: 'scheduled',
} satisfies ScheduledEvent;

const cancelledEvent = {
  ...scheduledEvent,
  id: 'cancelled-show',
  name: 'Cancelled show',
  status: 'cancelled',
} satisfies ScheduledEvent;

const okResponse = (body: unknown) => ({
  ok: true,
  status: 200,
  text: async () => JSON.stringify(body),
});

const errorResponse = {
  ok: false,
  status: 503,
  text: async () => 'unavailable',
};

describe('ScheduleCalendar actions', () => {
  it('renders list actions without triggering the event click handler', () => {
    const onEdit = jest.fn();
    const onEventClick = jest.fn();
    const onAdd = jest.fn();
    const onCalendar = jest.fn();

    render(
      <ScheduleCalendar
        events={[scheduledEvent]}
        initialView="listMonth"
        headerAction={{ text: 'Add event', onClick: onAdd }}
        calendarAction={{ text: 'Add to calendar', onClick: onCalendar }}
        onEventClick={onEventClick}
        renderEventActions={event => (
          <button type="button" onClick={() => onEdit(event)}>
            Edit
          </button>
        )}
      />,
    );

    fireEvent.click(screen.getByRole('button', { name: 'Edit' }));
    expect(onEdit).toHaveBeenCalledWith(scheduledEvent);
    expect(onEventClick).not.toHaveBeenCalled();
    fireEvent.click(screen.getByRole('button', { name: 'Add event' }));
    expect(onAdd).toHaveBeenCalledTimes(1);
    fireEvent.click(screen.getByRole('button', { name: 'Add to calendar' }));
    expect(onCalendar).toHaveBeenCalledTimes(1);
  });
});

describe('ScheduleTab', () => {
  const originalFetch = global.fetch;

  afterEach(() => {
    global.fetch = originalFetch;
    jest.clearAllMocks();
  });

  it('opens the current iCalendar feed', async () => {
    global.fetch = jest
      .fn()
      .mockResolvedValue(okResponse([scheduledEvent])) as unknown as typeof fetch;
    const open = jest.spyOn(window, 'open').mockImplementation();

    render(<ScheduleTab />);

    fireEvent.click(await screen.findByRole('button', { name: 'Add to calendar' }));
    expect(open).toHaveBeenCalledWith('/api/schedule.ics', '_blank', 'noopener');
    open.mockRestore();
  });

  it('loads the visible range and renders scheduled and cancelled occurrences', async () => {
    global.fetch = jest
      .fn()
      .mockResolvedValue(okResponse([scheduledEvent, cancelledEvent])) as unknown as typeof fetch;

    render(<ScheduleTab />);

    expect(await screen.findByText('Scheduled show')).toBeInTheDocument();
    expect(screen.getByText('Cancelled show')).toBeInTheDocument();
    expect(screen.getAllByText('Cancelled').length).toBeGreaterThan(0);
    expect(screen.getByRole('button', { name: 'Add to calendar' })).toBeInTheDocument();
    expect(global.fetch).toHaveBeenCalledWith(
      expect.stringMatching(/\/api\/schedule\?from=.*&to=.*/),
      expect.objectContaining({ method: 'GET' }),
    );
  });

  it('offers a retry after a schedule request fails', async () => {
    const fetchMock = jest
      .fn()
      .mockResolvedValueOnce(errorResponse)
      .mockResolvedValueOnce(okResponse([scheduledEvent]));
    global.fetch = fetchMock as unknown as typeof fetch;

    render(<ScheduleTab />);

    expect(await screen.findByText('Unable to load the schedule')).toBeInTheDocument();
    fireEvent.click(screen.getByRole('button', { name: 'Try again' }));

    await waitFor(() => expect(screen.getByText('Scheduled show')).toBeInTheDocument());
    expect(fetchMock).toHaveBeenCalledTimes(2);
  });
});

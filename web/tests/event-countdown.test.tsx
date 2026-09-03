import { act, render, screen } from '@testing-library/react';
import { hydrateRoot } from 'react-dom/client';
import { renderToString } from 'react-dom/server.node';
import '@testing-library/jest-dom';
import { EventCountdown } from '../components/ui/EventCountdown/EventCountdown';
import type { ScheduledEventStatus } from '../interfaces/server-status.model';

jest.mock('next/router', () => ({
  useRouter: () => ({
    query: {},
    pathname: '/',
    asPath: '/',
    push: jest.fn(),
    replace: jest.fn(),
  }),
}));

const event = {
  id: 'event-1',
  name: 'Scheduled show',
  description: 'A useful description.',
  startTime: '2030-09-06T22:10:12Z',
  durationMinutes: 90,
  chatOpen: false,
} satisfies ScheduledEventStatus;

describe('EventCountdown', () => {
  beforeEach(() => {
    jest.useFakeTimers();
    jest.setSystemTime(new Date('2030-09-01T18:00:00Z'));
  });

  afterEach(() => {
    jest.useRealTimers();
  });

  it('shows the event details, countdown, and last-live metadata', async () => {
    await act(async () => {
      render(<EventCountdown event={event} lastLive={new Date('2030-08-31T18:00:00Z')} />);
    });
    expect(screen.getByText('Scheduled show')).toBeInTheDocument();
    expect(screen.getByRole('link', { name: 'Add to calendar' })).toHaveAttribute(
      'href',
      '/api/schedule.ics',
    );
  });

  it('hydrates without rendering a clock value on the server', async () => {
    const container = document.createElement('div');
    const errors: unknown[] = [];
    const originalError = console.error;
    console.error = (...args: unknown[]) => errors.push(args);

    try {
      container.innerHTML = renderToString(<EventCountdown event={event} />);
      jest.setSystemTime(new Date('2030-09-01T18:00:01Z'));
      await act(async () => {
        hydrateRoot(container, <EventCountdown event={event} />);
      });
    } finally {
      console.error = originalError;
    }

    expect(errors).toEqual([]);
  });

  it('updates to live now at the event start', () => {
    render(<EventCountdown event={event} />);

    act(() => {
      jest.setSystemTime(new Date(event.startTime));
      jest.advanceTimersByTime(1000);
    });

    expect(screen.getByText('Starting soon')).toBeInTheDocument();
  });
});

import type { Meta } from '@storybook/nextjs';
import { http, HttpResponse } from 'msw';
import { ScheduleEventPage } from './ScheduleEventPage';

const upcomingEvent = {
  id: 'weekly-show',
  name: 'Night Signal: live from the studio',
  description: 'A quiet hour of field recordings, new music, and conversation.',
  startTime: new Date(Date.now() + 2 * 60 * 60 * 1000).toISOString(),
  durationMinutes: 90,
  timezone: Intl.DateTimeFormat().resolvedOptions().timeZone || 'UTC',
  status: 'scheduled' as const,
};

const cancelledEvent = {
  ...upcomingEvent,
  id: 'cancelled-show',
  name: 'The cancelled broadcast',
  status: 'cancelled' as const,
};

const endedEvent = {
  ...upcomingEvent,
  id: 'ended-show',
  name: 'Yesterday on the air',
  startTime: new Date(Date.now() - 3 * 60 * 60 * 1000).toISOString(),
};

const config = {
  name: 'Night Signal',
  appearanceVariables: {},
  pluginStyles: '',
  customStyles: '',
};

const meta = {
  title: 'owncast/Components/Schedule event page',
  component: ScheduleEventPage,
  parameters: {
    msw: {
      handlers: [
        http.get('/api/schedule*', () =>
          HttpResponse.json([upcomingEvent, cancelledEvent, endedEvent]),
        ),
        http.get('/api/config', () => HttpResponse.json(config)),
      ],
    },
  },
} satisfies Meta<typeof ScheduleEventPage>;

export default meta;

export const Upcoming = {
  args: { eventID: upcomingEvent.id, logoSrc: '/project/logo-semisimple-white.svg' },
};

export const Cancelled = {
  args: { eventID: cancelledEvent.id, logoSrc: '/project/logo-semisimple-white.svg' },
};

export const Ended = {
  args: { eventID: endedEvent.id, logoSrc: '/project/logo-semisimple-white.svg' },
};

export const Loading = {
  args: { eventID: upcomingEvent.id, logoSrc: '/project/logo-semisimple-white.svg' },
  parameters: {
    msw: {
      handlers: [
        http.get('/api/schedule*', () => {
          const { promise, resolve } = Promise.withResolvers<Response>();
          window.setTimeout(() => resolve(HttpResponse.json([upcomingEvent])), 10000);
          return promise;
        }),
        http.get('/api/config', () => HttpResponse.json(config)),
      ],
    },
  },
};

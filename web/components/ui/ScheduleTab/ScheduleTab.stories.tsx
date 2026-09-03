import type { Meta, StoryFn } from '@storybook/nextjs';
import { http, HttpResponse } from 'msw';
import { ScheduleTab } from './ScheduleTab';

const events = [
  {
    id: 'weekly-show',
    name: 'Weekly show',
    description: 'A scheduled example stream.',
    startTime: '2026-09-08T18:00:00-04:00',
    durationMinutes: 90,
    timezone: 'America/New_York',
    status: 'scheduled',
  },
  {
    id: 'cancelled-show',
    name: 'Cancelled show',
    description: 'This occurrence will not go live.',
    startTime: '2026-09-10T18:00:00-04:00',
    durationMinutes: 60,
    timezone: 'America/New_York',
    status: 'cancelled',
  },
];

const meta = {
  title: 'owncast/Components/Schedule tab',
  component: ScheduleTab,
  parameters: {
    msw: {
      handlers: [http.get('/api/schedule*', () => HttpResponse.json(events))],
    },
  },
} satisfies Meta<typeof ScheduleTab>;

export default meta;

const Template: StoryFn<typeof ScheduleTab> = () => <ScheduleTab />;

export const Example = {
  render: Template,
};

export const Loading = {
  render: Template,
  parameters: {
    msw: {
      handlers: [
        http.get('/api/schedule*', async () => {
          await new Promise(resolve => setTimeout(resolve, 10000));
          return HttpResponse.json(events);
        }),
      ],
    },
  },
};

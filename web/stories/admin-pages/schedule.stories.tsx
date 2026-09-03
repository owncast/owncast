import { Meta, StoryObj } from '@storybook/nextjs';
import { http, HttpResponse } from 'msw';
import Schedule from '../../pages/admin/schedule';

// Render smoke for the admin Schedule page: the feature toggle, reminder
// message field, both tables (one-off + materialized occurrences, and the
// recurring series that produced them) and the add-event button. A throwing
// story here is the breakage net for the schedule admin surface.
const schedule = {
  series: [
    {
      id: 'series1',
      name: 'Weekly show',
      description: 'Same time every week.',
      recurrence: 'DTSTART;TZID=America/Los_Angeles:20260710T180000\nRRULE:FREQ=WEEKLY;BYDAY=FR',
      durationMinutes: 90,
      active: true,
    },
  ],
  events: [
    {
      id: 'event1',
      name: 'One-off special',
      description: 'A special stream.',
      startTime: '2026-07-20T01:00:00Z',
      durationMinutes: 120,
      timezone: 'America/Los_Angeles',
      status: 'scheduled',
    },
    {
      id: 'event2',
      seriesId: 'series1',
      name: 'Weekly show',
      description: 'Same time every week.',
      startTime: '2026-07-11T01:00:00Z',
      durationMinutes: 90,
      timezone: 'America/Los_Angeles',
      status: 'scheduled',
    },
    {
      id: 'event3',
      seriesId: 'series1',
      name: 'Weekly show',
      description: 'Same time every week.',
      startTime: '2026-07-18T01:00:00Z',
      durationMinutes: 90,
      timezone: 'America/Los_Angeles',
      status: 'cancelled',
    },
  ],
};

const meta = {
  title: 'owncast/Admin/Pages/Schedule',
  component: Schedule,
  parameters: {
    msw: {
      handlers: [http.get('*/api/admin/schedule', () => HttpResponse.json(schedule))],
    },
  },
} satisfies Meta<typeof Schedule>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

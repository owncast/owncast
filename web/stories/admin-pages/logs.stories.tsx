import { Meta, StoryObj } from '@storybook/nextjs';
import { http, HttpResponse } from 'msw';
import Logs from '../../pages/admin/logs';

// Render smoke for the admin Logs page. The page renders nothing at all
// without log data (LogTable returns null when empty), so the story feeds a
// few entries through msw to make the table, level tags and pagination
// visible.
const logs = [
  {
    level: 'info',
    time: '2026-07-01T18:00:00Z',
    message: 'Owncast v0.2.0 started. Web server is listening on port 8080.',
  },
  {
    level: 'warning',
    time: '2026-07-01T18:05:12Z',
    message: 'The stream is using a video codec that some devices may not support.',
  },
  {
    level: 'error',
    time: '2026-07-01T18:07:45Z',
    message: 'Unable to connect to the directory at https://directory.owncast.online.',
  },
];

const meta = {
  title: 'owncast/Admin/Pages/Logs',
  component: Logs,
  parameters: {
    msw: {
      handlers: [http.get('*/api/admin/logs', () => HttpResponse.json(logs))],
    },
  },
} satisfies Meta<typeof Logs>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

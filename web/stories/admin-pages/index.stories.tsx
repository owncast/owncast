import { Meta, StoryObj } from '@storybook/nextjs';
import { http, HttpResponse } from 'msw';
import AdminHome from '../../pages/admin/index';

// Render smoke for the admin Home page: a story that throws fails the
// Storybook build, catching page-level breakage on every PR. Renders with
// default (empty) server config context.
//
// The news feed fetches owncast.online for real; snapshots of live content
// change whenever a newsletter goes out, so pin the feed to empty.
const meta = {
  title: 'owncast/Admin/Pages/Home',
  component: AdminHome,
  parameters: {
    msw: {
      handlers: [
        http.get('https://owncast.online/news/index.json', () => HttpResponse.json({ items: [] })),
      ],
    },
  },
} satisfies Meta<typeof AdminHome>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

import { Meta, StoryObj } from '@storybook/nextjs';
import { http, HttpResponse } from 'msw';
import ConfigFeatured from '../../pages/admin/config-featured';

// Render smoke for the admin Config Featured Streams page: a story that throws fails the
// Storybook build, catching page-level breakage on every PR. Renders with
// default (empty) server config context.
//
// The page's hooks fetch three admin endpoints on mount; without handlers the
// failures pop error toasts that race the snapshot and make the story flaky.
const meta = {
  title: 'owncast/Admin/Pages/Config Featured Streams',
  component: ConfigFeatured,
  parameters: {
    msw: {
      handlers: [
        http.get('/api/admin/federation/servers', () => HttpResponse.json([])),
        http.get('/api/admin/federation/feature-requests', () => HttpResponse.json([])),
        http.get('/api/admin/followers/directory', () => HttpResponse.json([])),
      ],
    },
  },
} satisfies Meta<typeof ConfigFeatured>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

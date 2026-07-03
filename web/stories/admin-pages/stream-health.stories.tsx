import { Meta, StoryObj } from '@storybook/nextjs';
import StreamHealth from '../../pages/admin/stream-health';

// Render smoke for the admin Stream Health page: a story that throws fails the
// Storybook build, catching page-level breakage on every PR. Renders with
// default (empty) server config context, no backend.
const meta = {
  title: 'owncast/Admin/Pages/Stream Health',
  component: StreamHealth,
} satisfies Meta<typeof StreamHealth>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

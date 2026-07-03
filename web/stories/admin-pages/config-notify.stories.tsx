import { Meta, StoryObj } from '@storybook/nextjs';
import ConfigNotify from '../../pages/admin/config-notify';

// Render smoke for the admin Config Notifications page: a story that throws fails the
// Storybook build, catching page-level breakage on every PR. Renders with
// default (empty) server config context, no backend.
const meta = {
  title: 'owncast/Admin/Pages/Config Notifications',
  component: ConfigNotify,
} satisfies Meta<typeof ConfigNotify>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

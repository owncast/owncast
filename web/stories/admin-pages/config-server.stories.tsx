import { Meta, StoryObj } from '@storybook/nextjs';
import ConfigServer from '../../pages/admin/config/server';

// Render smoke for the admin Config Server page: a story that throws fails the
// Storybook build, catching page-level breakage on every PR. Renders with
// default (empty) server config context, no backend.
const meta = {
  title: 'owncast/Admin/Pages/Config Server',
  component: ConfigServer,
} satisfies Meta<typeof ConfigServer>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

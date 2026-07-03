import { Meta, StoryObj } from '@storybook/nextjs';
import ConfigChat from '../../pages/admin/config-chat';

// Render smoke for the admin Config Chat page: a story that throws fails the
// Storybook build, catching page-level breakage on every PR. Renders with
// default (empty) server config context, no backend.
const meta = {
  title: 'owncast/Admin/Pages/Config Chat',
  component: ConfigChat,
} satisfies Meta<typeof ConfigChat>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

import { Meta, StoryObj } from '@storybook/nextjs';
import ConfigSocialItems from '../../pages/admin/config-social-items';

// Render smoke for the admin Config Social Items page: a story that throws fails the
// Storybook build, catching page-level breakage on every PR. Renders with
// default (empty) server config context, no backend.
const meta = {
  title: 'owncast/Admin/Pages/Config Social Items',
  component: ConfigSocialItems,
} satisfies Meta<typeof ConfigSocialItems>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

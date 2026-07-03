import { Meta, StoryObj } from '@storybook/nextjs';
import ConfigFeatured from '../../pages/admin/config-featured';

// Render smoke for the admin Config Featured Streams page: a story that throws fails the
// Storybook build, catching page-level breakage on every PR. Renders with
// default (empty) server config context, no backend.
const meta = {
  title: 'owncast/Admin/Pages/Config Featured Streams',
  component: ConfigFeatured,
} satisfies Meta<typeof ConfigFeatured>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

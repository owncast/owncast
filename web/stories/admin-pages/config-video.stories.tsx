import { Meta, StoryObj } from '@storybook/nextjs';
import ConfigVideo from '../../pages/admin/config-video';

// Render smoke for the admin Config Video page: a story that throws fails the
// Storybook build, catching page-level breakage on every PR. Renders with
// default (empty) server config context, no backend.
const meta = {
  title: 'owncast/Admin/Pages/Config Video',
  component: ConfigVideo,
} satisfies Meta<typeof ConfigVideo>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

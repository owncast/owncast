import { Meta, StoryObj } from '@storybook/nextjs';
import Plugins from '../../pages/admin/plugins';

// Render smoke for the admin Plugins page: a story that throws fails the
// Storybook build, catching page-level breakage on every PR. Renders with
// default (empty) server config context, no backend.
const meta = {
  title: 'owncast/Admin/Pages/Plugins',
  component: Plugins,
} satisfies Meta<typeof Plugins>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

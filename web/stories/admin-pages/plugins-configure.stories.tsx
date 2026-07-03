import { Meta, StoryObj } from '@storybook/nextjs';
import PluginConfigure from '../../pages/admin/plugins/configure';

// Render smoke for the admin Plugin Configure page: a story that throws fails the
// Storybook build, catching page-level breakage on every PR. Renders with
// default (empty) server config context, no backend.
const meta = {
  title: 'owncast/Admin/Pages/Plugin Configure',
  component: PluginConfigure,
} satisfies Meta<typeof PluginConfigure>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

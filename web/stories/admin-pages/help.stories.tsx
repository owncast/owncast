import { Meta, StoryObj } from '@storybook/nextjs';
import Help from '../../pages/admin/help';

// Render smoke for the admin Help page: a story that throws fails the
// Storybook build, catching page-level breakage on every PR. Renders with
// default (empty) server config context, no backend.
const meta = {
  title: 'owncast/Admin/Pages/Help',
  component: Help,
} satisfies Meta<typeof Help>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

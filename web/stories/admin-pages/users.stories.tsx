import { Meta, StoryObj } from '@storybook/nextjs';
import Users from '../../pages/admin/users';

// Render smoke for the admin Users page: a story that throws fails the
// Storybook build, catching page-level breakage on every PR. Renders with
// default (empty) server config context, no backend.
const meta = {
  title: 'owncast/Admin/Pages/Users',
  component: Users,
} satisfies Meta<typeof Users>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

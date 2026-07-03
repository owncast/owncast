import { Meta, StoryObj } from '@storybook/nextjs';
import AdminHome from '../../pages/admin/index';

// Render smoke for the admin Home page: a story that throws fails the
// Storybook build, catching page-level breakage on every PR. Renders with
// default (empty) server config context, no backend.
const meta = {
  title: 'owncast/Admin/Pages/Home',
  component: AdminHome,
} satisfies Meta<typeof AdminHome>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

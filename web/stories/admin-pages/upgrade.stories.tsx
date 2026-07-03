import { Meta, StoryObj } from '@storybook/nextjs';
import Upgrade from '../../pages/admin/upgrade';

// Render smoke for the admin Upgrade page: a story that throws fails the
// Storybook build, catching page-level breakage on every PR. Renders with
// default (empty) server config context, no backend.
const meta = {
  title: 'owncast/Admin/Pages/Upgrade',
  component: Upgrade,
} satisfies Meta<typeof Upgrade>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

import { Meta, StoryObj } from '@storybook/nextjs';
import Actions from '../../pages/admin/actions';

// Render smoke for the admin Actions page: a story that throws fails the
// Storybook build, catching page-level breakage on every PR. Renders with
// default (empty) server config context, no backend.
const meta = {
  title: 'owncast/Admin/Pages/Actions',
  component: Actions,
} satisfies Meta<typeof Actions>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

import { Meta, StoryObj } from '@storybook/nextjs';
import Webhooks from '../../pages/admin/webhooks';

// Render smoke for the admin Webhooks page: a story that throws fails the
// Storybook build, catching page-level breakage on every PR. Renders with
// default (empty) server config context, no backend.
const meta = {
  title: 'owncast/Admin/Pages/Webhooks',
  component: Webhooks,
} satisfies Meta<typeof Webhooks>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

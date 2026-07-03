import { Meta, StoryObj } from '@storybook/nextjs';
import FederationFollowers from '../../pages/admin/federation/followers';

// Render smoke for the admin Federation Followers page: a story that throws fails the
// Storybook build, catching page-level breakage on every PR. Renders with
// default (empty) server config context, no backend.
const meta = {
  title: 'owncast/Admin/Pages/Federation Followers',
  component: FederationFollowers,
} satisfies Meta<typeof FederationFollowers>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

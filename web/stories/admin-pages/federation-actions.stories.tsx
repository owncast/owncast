import { Meta, StoryObj } from '@storybook/nextjs';
import FederationActions from '../../pages/admin/federation/actions';

// Render smoke for the admin Federation Actions page: a story that throws fails the
// Storybook build, catching page-level breakage on every PR. Renders with
// default (empty) server config context, no backend.
const meta = {
  title: 'owncast/Admin/Pages/Federation Actions',
  component: FederationActions,
} satisfies Meta<typeof FederationActions>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

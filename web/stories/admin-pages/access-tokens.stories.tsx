import { Meta, StoryObj } from '@storybook/nextjs';
import AccessTokens from '../../pages/admin/access-tokens';

// Render smoke for the admin Access Tokens page: a story that throws fails the
// Storybook build, catching page-level breakage on every PR. Renders with
// default (empty) server config context, no backend.
const meta = {
  title: 'owncast/Admin/Pages/Access Tokens',
  component: AccessTokens,
} satisfies Meta<typeof AccessTokens>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

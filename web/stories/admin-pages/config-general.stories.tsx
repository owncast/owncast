import { Meta, StoryObj } from '@storybook/nextjs';
import ConfigGeneral from '../../pages/admin/config/general';

// Render smoke for the admin Config General page: a story that throws fails the
// Storybook build, catching page-level breakage on every PR. Renders with
// default (empty) server config context, no backend.
const meta = {
  title: 'owncast/Admin/Pages/Config General',
  component: ConfigGeneral,
} satisfies Meta<typeof ConfigGeneral>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

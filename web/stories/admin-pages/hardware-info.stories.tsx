import { Meta, StoryObj } from '@storybook/nextjs';
import HardwareInfo from '../../pages/admin/hardware-info';

// Render smoke for the admin Hardware Info page: a story that throws fails the
// Storybook build, catching page-level breakage on every PR. Renders with
// default (empty) server config context, no backend.
const meta = {
  title: 'owncast/Admin/Pages/Hardware Info',
  component: HardwareInfo,
} satisfies Meta<typeof HardwareInfo>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

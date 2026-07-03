import { Meta, StoryObj } from '@storybook/nextjs';
import ViewerInfo from '../../pages/admin/viewer-info';

// Render smoke for the admin Viewer Info page: a story that throws fails the
// Storybook build, catching page-level breakage on every PR. Renders with
// default (empty) server config context, no backend.
const meta = {
  title: 'owncast/Admin/Pages/Viewer Info',
  component: ViewerInfo,
} satisfies Meta<typeof ViewerInfo>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

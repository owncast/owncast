import { Meta, StoryObj } from '@storybook/nextjs';
import { VideoVariantForm } from './VideoVariantForm';
import { DEFAULT_VARIANT_STATE } from '../../utils/config-constants';

// The densest form in the admin: sliders with marks and tooltips, collapse
// panels, toggle switches, numeric text fields and a Popconfirm. If any of
// these widgets render broken, stream output configuration is broken, so
// this form gets its own visual baseline.
const meta = {
  title: 'owncast/Admin/Video variant form',
  component: VideoVariantForm,
} satisfies Meta<typeof VideoVariantForm>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
  args: {
    dataState: DEFAULT_VARIANT_STATE,
    onUpdateField: () => null,
  },
};

// Passthrough disables the quality controls and shows the warning state.
export const VideoPassthroughEnabled: Story = {
  args: {
    dataState: { ...DEFAULT_VARIANT_STATE, videoPassthrough: true },
    onUpdateField: () => null,
  },
};

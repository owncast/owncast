import { Meta, StoryObj } from '@storybook/nextjs';
import { ClipPlayer } from './ClipPlayer';

const meta = {
  title: 'owncast/Player/Clip player',
  component: ClipPlayer,
} satisfies Meta<typeof ClipPlayer>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
  args: {
    source: 'https://watch.owncast.online/hls/stream.m3u8',
  },
};

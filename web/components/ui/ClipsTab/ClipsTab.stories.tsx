import { Meta, StoryObj } from '@storybook/nextjs';
import { ClipsTab } from './ClipsTab';

const clips = [
  {
    id: 'example-clip',
    streamId: 'example-stream',
    title: 'The best moment',
    streamTitle: 'Example stream',
    clippedBy: 'Viewer',
    relativeStartTime: 30,
    relativeEndTime: 75,
    durationSeconds: 45,
    manifest: 'https://watch.owncast.online/hls/stream.m3u8',
    timestamp: '2026-08-08T12:00:00Z',
  },
];

const meta = {
  title: 'owncast/Components/Clips tab',
  component: ClipsTab,
} satisfies Meta<typeof ClipsTab>;

export default meta;
type Story = StoryObj<typeof meta>;

export const WithClips: Story = {
  loaders: [
    async () => {
      window.fetch = async () => new Response(JSON.stringify(clips));
      return {};
    },
  ],
};

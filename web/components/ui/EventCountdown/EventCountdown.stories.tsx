import type { Meta, StoryObj } from '@storybook/nextjs';
import { EventCountdown } from './EventCountdown';

const meta = {
  title: 'owncast/Viewer/Event Countdown',
  component: EventCountdown,
} satisfies Meta<typeof EventCountdown>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
  args: {
    event: {
      id: 'event-1',
      name: 'Community stream',
      description: 'Join us for a live community stream.',
      startTime: new Date(Date.now() + 5 * 86400000 + 4 * 3600000).toISOString(),
      durationMinutes: 90,
      chatOpen: false,
    },
  },
};

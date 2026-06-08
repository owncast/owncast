import React from 'react';
import { Meta, StoryFn } from '@storybook/react';
import { StreamCard, StreamCardProps } from './StreamCard';

export default {
  title: 'owncast/Components/StreamCard',
  component: StreamCard,
  parameters: {
    docs: {
      description: {
        component:
          'A card component for displaying featured Owncast stream status and information.',
      },
    },
  },
} as Meta<typeof StreamCard>;

const Template: StoryFn<StreamCardProps> = args => (
  <div style={{ maxWidth: '320px' }}>
    <StreamCard {...args} />
  </div>
);

export const OnlineStream = Template.bind({});
OnlineStream.args = {
  serverName: 'Awesome Stream Server',
  serverUrl: 'https://example-stream.com',
  serverLogo: 'https://picsum.photos/seed/1/64/64',
  streamTitle: 'Live Coding Session: Building a Chat App',
  streamDescription:
    'Join us for an exciting live coding session where we build a real-time chat application from scratch using React and WebSockets.',
  tags: ['programming', 'react', 'live-coding'],
  thumbnail: 'https://picsum.photos/seed/2/640/360',
  isOnline: true,
};

export const OfflineStream = Template.bind({});
OfflineStream.args = {
  serverName: 'Gaming Central',
  serverUrl: 'https://gaming-stream.com',
  serverLogo: 'https://picsum.photos/seed/3/64/64',
  tags: ['gaming', 'speedrun'],
  isOnline: false,
};

export const OnlineWithoutThumbnail = Template.bind({});
OnlineWithoutThumbnail.args = {
  serverName: 'Music Vibes',
  serverUrl: 'https://music-stream.com',
  serverLogo: 'https://picsum.photos/seed/4/64/64',
  streamTitle: '24/7 Lo-fi Hip Hop Radio',
  streamDescription: 'Relax and study with our continuous lo-fi hip hop stream.',
  tags: ['music', 'lofi', '24/7'],
  isOnline: true,
};

export const LongContent = Template.bind({});
LongContent.args = {
  serverName: 'Very Long Server Name That Should Be Truncated',
  serverUrl: 'https://long-name-stream.com',
  serverLogo: 'https://picsum.photos/seed/5/64/64',
  streamTitle:
    'This is a very long stream title that should be truncated when it exceeds the available space',
  streamDescription:
    'This is an extremely long description that goes on and on and on. It contains lots of details about the stream, what viewers can expect, the schedule, and many other things that would normally not fit in the card but should be elegantly truncated with an ellipsis.',
  tags: ['tag1', 'tag2', 'tag3', 'tag4', 'tag5'],
  thumbnail: 'https://picsum.photos/seed/6/640/360',
  isOnline: true,
};

export const MinimalInfo = Template.bind({});
MinimalInfo.args = {
  serverName: 'Simple Stream',
  serverUrl: 'https://simple.com',
  isOnline: false,
};

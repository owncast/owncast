import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { VideoPoster } from './VideoPoster';

export default {
  title: 'owncast/Player/Video poster',
  component: VideoPoster,
  parameters: {
    docs: {
      description: {
        component: `
- Sits on top of the video player when playback is not taking place.
- Shows the instance logo when the video is offline.
- Initial image is the logo when online.
- When the stream is online, will transition, via cross-fades, through the thumbnail.
- Will be removed when playback starts.`,
      },
    },
  },
} as ComponentMeta<typeof VideoPoster>;

// eslint-disable-next-line @typescript-eslint/no-unused-vars
const Template: ComponentStory<typeof VideoPoster> = args => <VideoPoster {...args} />;

// eslint-disable-next-line @typescript-eslint/no-unused-vars
export const Example1 = Template.bind({});
Example1.args = {
  initialSrc: 'https://watch.owncast.online/logo',
  src: 'https://watch.owncast.online/thumbnail.jpg',
  online: true,
};

export const Example2 = Template.bind({});
Example2.args = {
  initialSrc: 'https://listen.batstationrad.io/logo',
  src: 'https://listen.batstationrad.io//thumbnail.jpg',
  online: true,
};

export const Offline = Template.bind({});
Offline.args = {
  initialSrc: 'https://watch.owncast.online/logo',
  src: 'https://watch.owncast.online/thumbnail.jpg',
  online: false,
};

import { Meta } from '@storybook/react';
import { VideoPoster } from './VideoPoster';

const meta = {
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
} satisfies Meta<typeof VideoPoster>;

export default meta;

export const Example1 = {
  args: {
    initialSrc: 'https://watch.owncast.online/logo',
    src: 'https://watch.owncast.online/thumbnail.jpg',
    online: true,
  },
};

export const Example2 = {
  args: {
    initialSrc: 'https://listen.batstationrad.io/logo',
    src: 'https://listen.batstationrad.io//thumbnail.jpg',
    online: true,
  },
};

export const Offline = {
  args: {
    initialSrc: 'https://watch.owncast.online/logo',
    src: 'https://watch.owncast.online/thumbnail.jpg',
    online: false,
  },
};

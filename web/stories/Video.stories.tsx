import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import OwncastPlayer from '../components/video/OwncastPlayer';

const streams = {
  DemoServer: `https://watch.owncast.online/hls/stream.m3u8`,
  RetroStrangeTV: `https://live.retrostrange.com/hls/stream.m3u8`,
  localhost: `http://localhost:8080/hls/stream.m3u8`,
};

export default {
  title: 'owncast/Player',
  component: OwncastPlayer,
  argTypes: {
    source: {
      options: Object.keys(streams),
      mapping: streams,
      control: {
        type: 'select',
      },
    },
  },
  parameters: {},
} as ComponentMeta<typeof OwncastPlayer>;

const Template: ComponentStory<typeof OwncastPlayer> = args => <OwncastPlayer {...args} />;

export const LiveDemo = Template.bind({});
LiveDemo.args = {
  online: true,
  source: 'https://watch.owncast.online/hls/stream.m3u8',
};

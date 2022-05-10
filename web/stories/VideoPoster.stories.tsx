import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import VideoPoster from '../components/video/VideoPoster';

export default {
  title: 'owncast/Video poster',
  component: VideoPoster,
  parameters: {},
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

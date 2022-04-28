import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import StatusBar from '../components/video/StatusBar';

export default {
  title: 'owncast/Status bar',
  component: StatusBar,
  parameters: {},
} as ComponentMeta<typeof StatusBar>;

// eslint-disable-next-line @typescript-eslint/no-unused-vars
const Template: ComponentStory<typeof StatusBar> = args => <StatusBar {...args} />;

// eslint-disable-next-line @typescript-eslint/no-unused-vars
export const Online = Template.bind({});
Online.args = {
  online: true,
  viewers: 42,
  timer: '10:42',
};

export const Offline = Template.bind({});
Offline.args = {
  online: false,
};

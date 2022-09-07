import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { subHours } from 'date-fns';
import { Statusbar } from './Statusbar';

export default {
  title: 'owncast/Player/Status bar',
  component: Statusbar,
  parameters: {},
} as ComponentMeta<typeof Statusbar>;

const Template: ComponentStory<typeof Statusbar> = args => <Statusbar {...args} />;

export const Online = Template.bind({});
Online.args = {
  online: true,
  viewerCount: 42,
  lastConnectTime: subHours(new Date(), 3),
};

export const Offline = Template.bind({});
Offline.args = {
  online: false,
  lastDisconnectTime: subHours(new Date(), 3),
};

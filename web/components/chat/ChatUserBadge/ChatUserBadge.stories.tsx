import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import ChatUserBadge from './ChatUserBadge';

export default {
  title: 'owncast/Chat/Messages/User Flag',
  component: ChatUserBadge,
  argTypes: {
    userColor: {
      options: ['0', '1', '2', '3', '4', '5', '6', '7'],
      control: { type: 'select' },
    },
  },
} as ComponentMeta<typeof ChatUserBadge>;

const Template: ComponentStory<typeof ChatUserBadge> = args => <ChatUserBadge {...args} />;

export const Moderator = Template.bind({});
Moderator.args = {
  badge: 'mod',
  userColor: '5',
};

export const Authenticated = Template.bind({});
Authenticated.args = {
  badge: 'auth',
  userColor: '6',
};

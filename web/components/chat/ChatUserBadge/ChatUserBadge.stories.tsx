import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { ChatUserBadge } from './ChatUserBadge';
import { ModerationBadge } from './ModerationBadge';
import { AuthedUserBadge } from './AuthedUserBadge';
import { BotUserBadge } from './BotUserBadge';

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
const ModerationTemplate: ComponentStory<typeof ModerationBadge> = args => (
  <ModerationBadge {...args} />
);

const AuthedTemplate: ComponentStory<typeof ModerationBadge> = args => (
  <AuthedUserBadge {...args} />
);

const BotTemplate: ComponentStory<typeof BotUserBadge> = args => <BotUserBadge {...args} />;

export const Authenticated = AuthedTemplate.bind({});
Authenticated.args = {
  userColor: '3',
};

export const Moderator = ModerationTemplate.bind({});
Moderator.args = {
  userColor: '5',
};

export const Bot = BotTemplate.bind({});
Bot.args = {
  userColor: '7',
};

export const Generic = Template.bind({});
Generic.args = {
  badge: '?',
  userColor: '6',
};

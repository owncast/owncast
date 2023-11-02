import { StoryFn, Meta } from '@storybook/react';
import { ChatUserBadge } from './ChatUserBadge';
import { ModerationBadge } from './ModerationBadge';
import { AuthedUserBadge } from './AuthedUserBadge';
import { BotUserBadge } from './BotUserBadge';

const meta = {
  title: 'owncast/Chat/Messages/User Flag',
  component: ChatUserBadge,
  argTypes: {
    userColor: {
      options: ['0', '1', '2', '3', '4', '5', '6', '7'],
      control: { type: 'select' },
    },
  },
} satisfies Meta<typeof ChatUserBadge>;

export default meta;

const ModerationTemplate: StoryFn<typeof ModerationBadge> = args => <ModerationBadge {...args} />;

const AuthedTemplate: StoryFn<typeof ModerationBadge> = args => <AuthedUserBadge {...args} />;

const BotTemplate: StoryFn<typeof BotUserBadge> = args => <BotUserBadge {...args} />;

export const Authenticated = {
  render: AuthedTemplate,

  args: {
    userColor: '3',
  },
};

export const Moderator = {
  render: ModerationTemplate,

  args: {
    userColor: '5',
  },
};

export const Bot = {
  render: BotTemplate,

  args: {
    userColor: '7',
  },
};

export const Generic = {
  args: {
    badge: '?',
    userColor: '6',
  },
};

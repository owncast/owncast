import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { ChatPartMessage } from './ChatPartMessage';
import Mock from '../../../stories/assets/mocks/chatmessage-action.png';

export default {
  title: 'owncast/Chat/Messages/Chat Part',
  component: ChatPartMessage,
  argTypes: {
    userColor: {
      options: ['0', '1', '2', '3', '4', '5', '6', '7'],
      control: { type: 'select' },
    },
  },
  parameters: {
    design: {
      type: 'image',
      url: Mock,
    },
    docs: {
      description: {
        component: `This is shown when a chat participant parts.`,
      },
    },
  },
} as ComponentMeta<typeof ChatPartMessage>;

const Template: ComponentStory<typeof ChatPartMessage> = args => <ChatPartMessage {...args} />;

export const Regular = Template.bind({});
Regular.args = {
  displayName: 'RandomChatter',
  isAuthorModerator: false,
  userColor: 3,
};

export const Moderator = Template.bind({});
Moderator.args = {
  displayName: 'RandomChatter',
  isAuthorModerator: true,
  userColor: 2,
};

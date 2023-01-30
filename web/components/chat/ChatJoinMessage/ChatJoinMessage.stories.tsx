import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { ChatJoinMessage } from './ChatJoinMessage';
import Mock from '../../../stories/assets/mocks/chatmessage-action.png';

export default {
  title: 'owncast/Chat/Messages/Chat Join',
  component: ChatJoinMessage,
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
        component: `This is the message design an action takes place, such as a join or a name change.`,
      },
    },
  },
} as ComponentMeta<typeof ChatJoinMessage>;

const Template: ComponentStory<typeof ChatJoinMessage> = args => <ChatJoinMessage {...args} />;

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

import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import ChatSystemMessage from '../components/chat/ChatSystemMessage/ChatSystemMessage';
import Mock from './assets/mocks/chatmessage-system.png';

export default {
  title: 'owncast/Chat/Messages/System',
  component: ChatSystemMessage,
  parameters: {
    design: {
      type: 'image',
      url: Mock,
    },
    docs: {
      description: {
        component: `This is the message design used when the server sends a message to chat.`,
      },
    },
  },
} as ComponentMeta<typeof ChatSystemMessage>;

const Template: ComponentStory<typeof ChatSystemMessage> = args => <ChatSystemMessage {...args} />;

// eslint-disable-next-line @typescript-eslint/no-unused-vars
export const Basic = Template.bind({});
Basic.args = {
  message: 'This is a system message.',
};

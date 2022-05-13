import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import ChatActionMessage from '../components/chat/ChatActionMessage';
import Mock from './assets/mocks/chatmessage-action.png';

export default {
  title: 'owncast/Chat/Messages/Chat action',
  component: ChatActionMessage,
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
} as ComponentMeta<typeof ChatActionMessage>;

// eslint-disable-next-line @typescript-eslint/no-unused-vars
const Template: ComponentStory<typeof ChatActionMessage> = args => <ChatActionMessage {...args} />;

// eslint-disable-next-line @typescript-eslint/no-unused-vars
export const Basic = Template.bind({});

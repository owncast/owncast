import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { ChatActionMessage } from './ChatActionMessage';
import Mock from '../../../stories/assets/mocks/chatmessage-action.png';

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

const Template: ComponentStory<typeof ChatActionMessage> = args => <ChatActionMessage {...args} />;

export const Basic = Template.bind({});
Basic.args = {
  body: 'This is a basic action message.',
};

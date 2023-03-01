import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { ChatNameChangeMessage } from './ChatNameChangeMessage';

export default {
  title: 'owncast/Chat/Messages/Chat name change',
  component: ChatNameChangeMessage,
} as ComponentMeta<typeof ChatNameChangeMessage>;

const Template: ComponentStory<typeof ChatNameChangeMessage> = args => (
  <ChatNameChangeMessage {...args} />
);

export const Basic = Template.bind({});
Basic.args = {
  message: {
    oldName: 'JohnnyOldName',
    user: {
      displayName: 'JohnnyNewName',
      displayColor: '3',
    },
  },
};

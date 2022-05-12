import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import ChatSystemMessage from '../components/chat/ChatSystemMessage';

export default {
  title: 'owncast/Chat/Messages/System',
  component: ChatSystemMessage,
  parameters: {},
} as ComponentMeta<typeof ChatSystemMessage>;

const Template: ComponentStory<typeof ChatSystemMessage> = args => <ChatSystemMessage {...args} />;

// eslint-disable-next-line @typescript-eslint/no-unused-vars
export const Basic = Template.bind({});

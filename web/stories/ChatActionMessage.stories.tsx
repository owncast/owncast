import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import ChatActionMessage from '../components/chat/ChatActionMessage';

export default {
  title: 'owncast/Chat/Messages/Chat action',
  component: ChatActionMessage,
  parameters: {},
} as ComponentMeta<typeof ChatActionMessage>;

// eslint-disable-next-line @typescript-eslint/no-unused-vars
const Template: ComponentStory<typeof ChatActionMessage> = args => <ChatActionMessage {...args} />;

// eslint-disable-next-line @typescript-eslint/no-unused-vars
export const Basic = Template.bind({});

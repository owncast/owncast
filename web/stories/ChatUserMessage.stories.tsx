import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import UserChatMessage from '../components/chat/ChatUserMessage';

export default {
  title: 'owncast/Chat/Messages/Standard user',
  component: UserChatMessage,
  parameters: {},
} as ComponentMeta<typeof UserChatMessage>;

// eslint-disable-next-line @typescript-eslint/no-unused-vars
const Template: ComponentStory<typeof UserChatMessage> = args => <UserChatMessage />;

// eslint-disable-next-line @typescript-eslint/no-unused-vars
export const Basic = Template.bind({});

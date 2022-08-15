import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import ChatSocialMessage from '../components/chat/ChatSocialMessage';

export default {
  title: 'owncast/Chat/Messages/Social-fediverse event',
  component: ChatSocialMessage,
  parameters: {},
} as ComponentMeta<typeof ChatSocialMessage>;

// eslint-disable-next-line @typescript-eslint/no-unused-vars
const Template: ComponentStory<typeof ChatSocialMessage> = args => <ChatSocialMessage {...args} />;

// eslint-disable-next-line @typescript-eslint/no-unused-vars
export const Basic = Template.bind({});

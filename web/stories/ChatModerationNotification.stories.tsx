import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import ChatModeratorNotification from '../components/chat/ChatModeratorNotification';

export default {
  title: 'owncast/Chat/Messages/Moderator notification',
  component: ChatModeratorNotification,
  parameters: {},
} as ComponentMeta<typeof ChatModeratorNotification>;

// eslint-disable-next-line @typescript-eslint/no-unused-vars
const Template: ComponentStory<typeof ChatModeratorNotification> = args => (
  <ChatModeratorNotification />
);

// eslint-disable-next-line @typescript-eslint/no-unused-vars
export const Basic = Template.bind({});

import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import ChatModerationActionMenu from '../components/chat/ChatModerationActionMenu/ChatModerationActionMenu';

export default {
  title: 'owncast/Chat/Moderation menu',
  component: ChatModerationActionMenu,
  parameters: {
    docs: {
      description: {
        component: `This should be a popup that is activated from a user's chat message. It should have actions to:
- Remove single message
- Ban user completely
- Open modal to see user details
        `,
      },
    },
  },
} as ComponentMeta<typeof ChatModerationActionMenu>;

// eslint-disable-next-line @typescript-eslint/no-unused-vars
const Template: ComponentStory<typeof ChatModerationActionMenu> = args => (
  <ChatModerationActionMenu
    accessToken="abc123"
    messageID="xxx"
    userDisplayName="Fake-User"
    userID="abc123"
  />
);

// eslint-disable-next-line @typescript-eslint/no-unused-vars
export const Basic = Template.bind({});

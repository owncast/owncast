import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import ChatModerationDetailsModal from '../components/chat/ChatModerationDetailsModal';

export default {
  title: 'owncast/Chat/Moderation modal',
  component: ChatModerationDetailsModal,
  parameters: {
    docs: {
      description: {
        component: `This should be a modal that gives the moderator more details about the user such as:
- When the user was created
- Other names they've used
- If they're authenticated, and using what method (IndieAuth, FediAuth)
`,
      },
    },
  },
} as ComponentMeta<typeof ChatModerationDetailsModal>;

// eslint-disable-next-line @typescript-eslint/no-unused-vars
const Template: ComponentStory<typeof ChatModerationDetailsModal> = args => (
  <ChatModerationDetailsModal />
);

// eslint-disable-next-line @typescript-eslint/no-unused-vars
export const Basic = Template.bind({});

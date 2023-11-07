import { Meta } from '@storybook/react';
import { ChatPartMessage } from './ChatPartMessage';
import Mock from '../../../stories/assets/mocks/chatmessage-action.png';

const meta = {
  title: 'owncast/Chat/Messages/Chat Part',
  component: ChatPartMessage,
  argTypes: {
    userColor: {
      options: ['0', '1', '2', '3', '4', '5', '6', '7'],
      control: { type: 'select' },
    },
  },
  parameters: {
    design: {
      type: 'image',
      url: Mock,
    },
    docs: {
      description: {
        component: `This is shown when a chat participant parts.`,
      },
    },
  },
} satisfies Meta<typeof ChatPartMessage>;

export default meta;

export const Regular = {
  args: {
    displayName: 'RandomChatter',
    isAuthorModerator: false,
    userColor: 3,
  },
};

export const Moderator = {
  args: {
    displayName: 'RandomChatter',
    isAuthorModerator: true,
    userColor: 2,
  },
};

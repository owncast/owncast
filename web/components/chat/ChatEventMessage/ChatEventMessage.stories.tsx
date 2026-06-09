import { Meta } from '@storybook/nextjs';
import { ChatEventMessage, ChatEventType } from './ChatEventMessage';
import Mock from '../../../stories/assets/mocks/chatmessage-action.png';

const meta = {
  title: 'owncast/Chat/Messages/Chat Event',
  component: ChatEventMessage,
  argTypes: {
    type: {
      options: [ChatEventType.Join, ChatEventType.Part],
      control: { type: 'radio' },
    },
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
        component: `This is shown when a chat participant joins or parts.`,
      },
    },
  },
} satisfies Meta<typeof ChatEventMessage>;

export default meta;

export const Join = {
  args: {
    type: ChatEventType.Join,
    displayName: 'RandomChatter',
    isAuthorModerator: false,
    userColor: 3,
  },
};

export const Part = {
  args: {
    type: ChatEventType.Part,
    displayName: 'RandomChatter',
    isAuthorModerator: false,
    userColor: 3,
  },
};

export const Moderator = {
  args: {
    type: ChatEventType.Join,
    displayName: 'RandomChatter',
    isAuthorModerator: true,
    userColor: 2,
  },
};

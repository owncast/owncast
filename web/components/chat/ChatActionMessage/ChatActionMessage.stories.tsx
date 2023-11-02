import { Meta } from '@storybook/react';
import { ChatActionMessage } from './ChatActionMessage';
import Mock from '../../../stories/assets/mocks/chatmessage-action.png';

const meta = {
  title: 'owncast/Chat/Messages/Chat action',
  component: ChatActionMessage,
  parameters: {
    design: {
      type: 'image',
      url: Mock,
    },
    docs: {
      description: {
        component: `This is the message design an action takes place, such as a join or a name change.`,
      },
    },
  },
} satisfies Meta<typeof ChatActionMessage>;

export default meta;

export const Basic = {
  args: {
    body: 'This is a basic action message.',
  },
};

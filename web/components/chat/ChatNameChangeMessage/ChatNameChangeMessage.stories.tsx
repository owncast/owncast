import { Meta } from '@storybook/react';
import { ChatNameChangeMessage } from './ChatNameChangeMessage';

const meta = {
  title: 'owncast/Chat/Messages/Chat name change',
  component: ChatNameChangeMessage,
} satisfies Meta<typeof ChatNameChangeMessage>;

export default meta;

export const Basic = {
  args: {
    message: {
      oldName: 'JohnnyOldName',
      user: {
        displayName: 'JohnnyNewName',
        displayColor: '3',
      },
    },
  },
};

import { Meta } from '@storybook/react';
import { FatalErrorStateModal } from './FatalErrorStateModal';

const meta = {
  title: 'owncast/Modals/Global error state',
  component: FatalErrorStateModal,
  parameters: {},
} satisfies Meta<typeof FatalErrorStateModal>;

export default meta;

export const Example = {
  args: {
    title: 'Example error title',
    message: 'Example error message',
  },
};

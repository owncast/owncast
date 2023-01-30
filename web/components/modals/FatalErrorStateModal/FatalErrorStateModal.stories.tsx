import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { FatalErrorStateModal } from './FatalErrorStateModal';

export default {
  title: 'owncast/Modals/Global error state',
  component: FatalErrorStateModal,
  parameters: {},
} as ComponentMeta<typeof FatalErrorStateModal>;

// eslint-disable-next-line @typescript-eslint/no-unused-vars
const Template: ComponentStory<typeof FatalErrorStateModal> = args => (
  <FatalErrorStateModal {...args} />
);

// eslint-disable-next-line @typescript-eslint/no-unused-vars
export const Example = Template.bind({});
Example.args = {
  title: 'Example error title',
  message: 'Example error message',
};

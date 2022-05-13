import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import NameChangeModal from '../components/modals/NameChangeModal';

export default {
  title: 'owncast/Modals/Name change',
  component: NameChangeModal,
  parameters: {},
} as ComponentMeta<typeof NameChangeModal>;

// eslint-disable-next-line @typescript-eslint/no-unused-vars
const Template: ComponentStory<typeof NameChangeModal> = args => <NameChangeModal />;

// eslint-disable-next-line @typescript-eslint/no-unused-vars
export const Basic = Template.bind({});

import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import IndieAuthModal from '../components/modals/IndieAuthModal';

const Example = () => (
  <div>
    <IndieAuthModal />
  </div>
);

export default {
  title: 'owncast/Modals/IndieAuth',
  component: IndieAuthModal,
  parameters: {},
} as ComponentMeta<typeof IndieAuthModal>;

// eslint-disable-next-line @typescript-eslint/no-unused-vars
const Template: ComponentStory<typeof IndieAuthModal> = args => <Example />;

// eslint-disable-next-line @typescript-eslint/no-unused-vars
export const Basic = Template.bind({});

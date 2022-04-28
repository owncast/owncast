import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import FediAuthModal from '../components/modals/FediAuthModal';

const Example = () => (
  <div>
    <FediAuthModal />
  </div>
);

export default {
  title: 'owncast/Modals/FediAuth',
  component: FediAuthModal,
  parameters: {},
} as ComponentMeta<typeof FediAuthModal>;

// eslint-disable-next-line @typescript-eslint/no-unused-vars
const Template: ComponentStory<typeof FediAuthModal> = args => <Example />;

// eslint-disable-next-line @typescript-eslint/no-unused-vars
export const Basic = Template.bind({});

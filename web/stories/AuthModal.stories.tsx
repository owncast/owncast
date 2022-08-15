import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import AuthModal from '../components/modals/AuthModal';

const Example = () => (
  <div>
    <AuthModal />
  </div>
);

export default {
  title: 'owncast/Modals/Auth',
  component: AuthModal,
  parameters: {},
} as ComponentMeta<typeof AuthModal>;

// eslint-disable-next-line @typescript-eslint/no-unused-vars
const Template: ComponentStory<typeof AuthModal> = args => <Example />;

// eslint-disable-next-line @typescript-eslint/no-unused-vars
export const Basic = Template.bind({});

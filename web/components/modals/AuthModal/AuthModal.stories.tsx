import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { RecoilRoot } from 'recoil';
import AuthModal from './AuthModal';

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
const Template: ComponentStory<typeof AuthModal> = args => (
  <RecoilRoot>
    <Example />
  </RecoilRoot>
);

// eslint-disable-next-line @typescript-eslint/no-unused-vars
export const Basic = Template.bind({});

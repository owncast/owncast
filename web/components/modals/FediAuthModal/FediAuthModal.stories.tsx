import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { FediAuthModal } from './FediAuthModal';
import FediAuthModalMock from '../../../stories/assets/mocks/fediauth-modal.png';

export default {
  title: 'owncast/Modals/FediAuth',
  component: FediAuthModal,
  parameters: {
    design: {
      type: 'image',
      url: FediAuthModalMock,
      scale: 0.5,
    },
  },
} as ComponentMeta<typeof FediAuthModal>;

const Template: ComponentStory<typeof FediAuthModal> = args => <FediAuthModal {...args} />;

export const NotYetAuthenticated = Template.bind({});
NotYetAuthenticated.args = {
  displayName: 'fake-user',
  authenticated: false,
  accessToken: '',
};

export const Authenticated = Template.bind({});
Authenticated.args = {
  displayName: 'fake-user',
  authenticated: true,
  accessToken: '',
};

import { Meta } from '@storybook/react';
import { FediAuthModal } from './FediAuthModal';
import FediAuthModalMock from '../../../stories/assets/mocks/fediauth-modal.png';

const meta = {
  title: 'owncast/Modals/FediAuth',
  component: FediAuthModal,
  parameters: {
    design: {
      type: 'image',
      url: FediAuthModalMock,
      scale: 0.5,
    },
  },
} satisfies Meta<typeof FediAuthModal>;

export default meta;

export const NotYetAuthenticated = {
  args: {
    displayName: 'fake-user',
    authenticated: false,
    accessToken: '',
  },
};

export const Authenticated = {
  args: {
    displayName: 'fake-user',
    authenticated: true,
    accessToken: '',
  },
};

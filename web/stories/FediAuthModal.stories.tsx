import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import FediAuthModal from '../components/modals/FediAuthModal';
import FediAuthModalMock from './assets/mocks/fediauth-modal.png';

const Example = () => (
  <div>
    <FediAuthModal />
  </div>
);

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

// eslint-disable-next-line @typescript-eslint/no-unused-vars
const Template: ComponentStory<typeof FediAuthModal> = args => <Example />;

// eslint-disable-next-line @typescript-eslint/no-unused-vars
export const Basic = Template.bind({});
